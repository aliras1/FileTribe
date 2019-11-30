pragma solidity ^0.5.10;

import "./openzeppelin/contracts/cryptography/MerkleProof.sol";
import "./openzeppelin/contracts/cryptography/MerkleTree.sol";

import "./interfaces/IConsensusCallback.sol";
import "./interfaces/IConsensus.sol";
import "./interfaces/IFileTribeDApp.sol";
import "./interfaces/IGroup.sol";
import "./common/Ownable.sol";
import "./ecOps.sol";


contract GroupDkg is Ownable, IConsensusCallback {

    /** 
     * DKG phases (optimistic case):
     * 
     * 0) Deploy the contract with a threshold (t), number of participants (n) 
     * and deposit in Wei.
     * 
     * 1) Enrollment - Each of the participants sends a deposit and a public key 
     * address (for encryption purpose) that he owns (i.e., has the corresponding 
     * secret key). At the end of the enrollment each participant receives an
     * unique index (1,...,n) that would be used with his identity for the rest
     * of the DKG protocol and the threshold signature scheme.
     *
     * (not on the contract) - each participant generates t+1 random sampled 
     * coefficients from the (mod q) field.
     *
     * 2) Each of the participants sends its public commitments (the generator 
     * exponentiated with t+1 coefficients) and encrypted private commitments 
     * for all of the other particpants.
     *
     * After all have validly committed each of the participant can take the
     * committed data and locally compute the followings:
     *  a) its own secret key and public key;
     *  b) the public keys for the rest of the participants;
     *  c) the group's public key (no one would know the corresponding secret 
     *     key);
     * With the above data each participant is able to sign, verify signature- 
     * shares and reconstruct from t+1 signatures the group signature.
     * 
     *
     *
     * In case one (or more) of the participants deviates from the instructions
     * of the protocol, we enable the other participants to file a complaint to
     * the contract. Such a complaint would terminate the DKG protocol with a
     * failure. Once a complaint is filed, a code is run to check the validity of
     * the complaint. When done the code decides which participant's deposit to 
     * slash and divide it among the other participants ("code is law").
     *
     * Each participant can send a complaint tx about one of the followings:
     *  a) 2 distinct participants offered the same public commitment
           (one is enough). (TODO)
     *  b) Some participant offered invalid commitment (invalid is: 
     *     duplicated, insufficient, unmatching commitments G1 to G2)
     *  c) Umatched private and public commitments.
     *  d) Time out.
     *
     */


    /**
     * Important note: at this point this contract purpose is as a
     * POC only, therefore its security is unreliable.
     */



    struct Participant {
        address payable ethPk; // Ethereum pk
        uint256[2] encPk; // pk for encryption
        bytes32 rootPubCommitG1;
        bytes32 rootPubCommitG2;
        bytes32 rootPrvCommit;
        uint256[2] yG1;
        string commitIpfsHash;
        bool isCommitted;
    }

    enum Phase { Enrollment, Commit, PostCommit, EndSuccess, EndFail } // start from 0

    event NewComplaint(
        address consensus
    );
    event PhaseChange(
        Phase phase
    );
    event NewCommit(
        uint16 committerIndex,
        bytes32 rootPubCommitG1,
        bytes32 rootPubCommitG2,
        bytes32 rootPrvCommit,
        uint256[2] yG1,
        string commitIpfsHash
    );
    event ParticipantJoined(
        uint16 index
    );
    event Debug(
        string message
    );
    event DebugProofG1(
        bytes32[] proof,
        bytes32 root,
        bytes32 leaf
    );
    event DebugProofG2(
        bytes32[] proof,
        bytes32 root,
        bytes32 leaf
    );


    Phase public curPhase;
    
     
    //uint256 public constant a = 0;
    //uint256 public constant b = 3;

    // G1 generator (on the curve)
    uint256[2] public g1 = [
        0x0000000000000000000000000000000000000000000000000000000000000001, 
        0x0000000000000000000000000000000000000000000000000000000000000002
    ];
    // G2 generator (on the curve)
    uint256[4] public g2 = [
        0x198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c2, 
        0x1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed,
        0x90689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b,
        0x12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa
    ];

    uint256 public depositWei;
    

    uint16 public t; // threshold
    uint16 public n; // numer of participants;
    uint16 public curN; // current num of participants
    uint16 public curNumCommittedLeft; // current num of participants that haven't committed 
    
    uint256 public phaseStart;
    uint256 public constant joinTimeout = 12;
    uint256 public constant commitTimeout = 11;
    uint256 public constant postCommitTimeout = 5;


    // mapping from node's index to a participant
    mapping (uint16 => Participant) public participants;
    uint256[2] public VK;
    uint16 finalizerIdx;
    IGroup _group;
    IFileTribeDApp _fileTribe;
    mapping (address => uint16) _cons2Idx;

    constructor() public Ownable(msg.sender) { }

    function init(uint16 threshold, uint16 numParticipants) public {
        t = threshold;
        n = numParticipants;
        curNumCommittedLeft = numParticipants;
        //depositWei = deposit;
        curPhase = Phase.Enrollment;

        require(n > t && t > 0, "wrong input");

        curN = 0;

        phaseStart = block.number;
    }

    modifier checkDeposit() {
        require(msg.value == depositWei, "wrong deposit");
        _;
    }
    modifier checkAuthorizedSender(uint16 index) {
        require(participants[index].ethPk == msg.sender, "not authorized sender");
        _; 
    }
    modifier beFalse(bool term) {
        require(!term);
        _;
    }
    modifier inPhase(Phase phase) {
        require(curPhase == phase, "wrong phase");
        _;
    }
    modifier notInPhase(Phase phase) {
        require(curPhase != phase, "wrong phase");
        _;
    }


    // Join the DKG (enrollment - phase 1).
    // A point on G1 that represents this participant's pk for encryption have
    // to be published. The publisher have to know the secret that generates
    // this point.
    function join(uint256[2] calldata encPk)
        external payable
        //checkDeposit()
        inPhase(Phase.Enrollment)
        returns(uint16 index)
    {
        // TODO: check pk

        uint16 cn = curN;
        address payable sender = msg.sender;

        // Check the pk isn't registered already
        for(uint16 i = 1; i <= cn; i++) {
            require(participants[i].ethPk != sender, "already joined");
        }

        cn++;
        participants[cn] = Participant({
            ethPk: sender,
            encPk: encPk,
            isCommitted: false,
            rootPubCommitG1: "",
            rootPubCommitG2: "",
            rootPrvCommit: "",
            commitIpfsHash: "",
            yG1: [uint256(0), uint256(0)]
        });

        curN = cn;
        if(cn == 1) {
            phaseStart = block.number;
        }

        // Abort if capacity on participants was reached
        if(cn == n) {
            curPhase = Phase.Commit;
            emit PhaseChange(Phase.Commit);
        }

        emit ParticipantJoined(cn);
        return cn;
    }    
    
    
    // Send commitments (phase 2). 
    //
    // pubCommitG1 is composed of t+1 commitments to local randomly sampled
    // coefficients. Each commitment should be on the G1 curve (affine 
    // coordinates) and therefore it has 2 coordinates. Thus, the input array
    // is of size (2t+2) and the i'th commitment will be in indices (2i) and
    // (2i+1).
    //
    // pubCommitG2 is composed of t+1 commitments to same sampled coefficients
    // from pubCommitG1. Each commitment should be on the G2 curve (affine 
    // coordinates) and therefore it has 4 coordinates. Thus, the input array
    // is of size (4t+4) and the i'th commitment will be in indices (4i),(4i+1),
    // (4i+2),(4i+3).   
    //
    // prCommit is an array of size n, where the first index matches the
    // first participant (participant index 1) and so forth. The commitment
    // is a calculation on the localy generated polynomial in the particpant's
    // index. This calculation is encrypted by the recevier pk for encryption.
    // The senderIndex private commitment is ignored and can be anything 
    // (but can't be skipped).
    // 
    // Note that this function does not verifies the committed data, it
    // should be done outside of this contract scope. In case of an
    // invalid committed data use complaints.
    function commit(
        uint16 senderIndex,
        bytes32 rootPubCommitG1,
        bytes32 rootPubCommitG2,
        bytes32 rootEncPrvCommit,
        uint256[2] calldata yG1,
        string calldata commitIpfsHash
    )
        external
        inPhase(Phase.Commit)
        checkAuthorizedSender(senderIndex)
        beFalse(participants[senderIndex].isCommitted)
    {
        // TODO: phase timeout, make prCommit encrypted, verify sender
        // index matches the sender's address.

        participants[senderIndex].rootPubCommitG1 = rootPubCommitG1;
        participants[senderIndex].rootPubCommitG2 = rootPubCommitG2;
        participants[senderIndex].rootPrvCommit = rootEncPrvCommit;
        participants[senderIndex].yG1 = yG1;
        participants[senderIndex].commitIpfsHash = commitIpfsHash;
        participants[senderIndex].isCommitted = true;
        emit NewCommit(
            senderIndex,
            rootPubCommitG1,
            rootPubCommitG2,
            rootEncPrvCommit,
            yG1,
            commitIpfsHash);

        uint16 committedNumLeft = curNumCommittedLeft - 1;
        curNumCommittedLeft = committedNumLeft;

        if(committedNumLeft == 0) {
            curPhase = Phase.PostCommit;
            phaseStart = block.number;
            emit PhaseChange(Phase.PostCommit);
        }
    }


    // Call this when in Phase.PostCommit for more than postCommitTimeout
    // blocks and no comlaint has to be made.
    function postCommitTimedOut(uint16 senderIndex, uint256[2] calldata vk)
        external
        inPhase(Phase.PostCommit)
        checkAuthorizedSender(senderIndex)
    {

        uint curBlockNum = block.number;

        require(curBlockNum > (phaseStart+postCommitTimeout), "hasn't reached timeout yet");

        curPhase = Phase.EndSuccess;
        finalizerIdx = senderIndex;
        VK = vk;

        emit PhaseChange(Phase.EndSuccess);
        //slash(0);
    }

    // Call this when in Phase.Enrollment for more than joinTimeout
    // blocks and not enough members have joined.
    function joinTimedOut()
        external
        inPhase(Phase.Enrollment)
    {
        uint curBlockNum = block.number;

        require(curBlockNum > (phaseStart+joinTimeout), "hasn't reached timeout yet");
        curPhase = Phase.EndFail;
        emit PhaseChange(Phase.EndFail);
        slash(0);
    }

    // Call this when in Phase.Commit for more than commitTimeout
    // blocks and not enough members have committed.
    function commitTimedOut()
        external
        inPhase(Phase.Commit)
    {
        uint curBlockNum = block.number;

        require(curBlockNum > (phaseStart+commitTimeout), "hasn't reached timeout yet");
        curPhase = Phase.EndFail;
        emit PhaseChange(Phase.EndFail);
        slashUncommitted();
    }

    // Returns the group PK.
    // This can only be performed after the DKG has ended. This
    // means only when the current phase is Phase.End .
    function getGroupPK()
        //inPhase(Phase.EndSuccess)
        public view returns(uint256[2] memory)
    {
        return VK;
    }

    function calculateGroupPK() private view returns(uint256[2] memory vk) {
        vk = participants[1].yG1;

        for(uint16 i = 2; i <= n; i++) {
            vk = ecOps.ecadd(vk, participants[i].yG1);
        }
    }


    ////////////////
    // Complaints //
    ////////////////

    function complaintInvalidIpfsHash(uint16 complainerIdx, uint16 accusedIdx)
        public
        checkAuthorizedSender(complainerIdx)
        inPhase(Phase.Commit)
    {
        IAccount acc = _fileTribe.getAccountOf(msg.sender);
        IConsensus cons = _fileTribe.createConsensus(acc);
        _cons2Idx[address(cons)] = complainerIdx;

        cons.propose(uint16ToBytes(accusedIdx), 0);

        emit NewComplaint(address(cons));
    }

    // A complaint on some public commit. If for some reason this
    // function fails it will slash the complainer deposit! (unless some
    // unauthorized address made the transaction or the wrong phase).
    //
    // The complaint should be called when the public commitments coming
    // from the G1 group does not match to the ones from G2 group (using pairing).
    function complaintPublicCommit(
        uint16 complainerIndex,
        uint16 accusedIndex,
        uint256[2] memory pubCommitG1,
        uint256[4] memory pubCommitG2,
        bytes32[] memory proofG1, // TODO: restrict the length of proof based on t
        bytes32[] memory proofG2 // TODO: restrict the length of proof based on t
    )
        public
        checkAuthorizedSender(complainerIndex)
        notInPhase(Phase.EndFail)
        notInPhase(Phase.EndSuccess)
    {
        curPhase = Phase.EndFail;
        emit PhaseChange(Phase.EndFail);

        Participant storage accused = participants[accusedIndex];
        if(!accused.isCommitted) {
            slash(complainerIndex);
            emit Debug("accused has not commtied: slashed complainer");
            return;
        }

        bytes32 leafG1 = keccak256(abi.encodePacked(pubCommitG1));
        emit DebugProofG1(proofG1, accused.rootPubCommitG1, leafG1);
        if (!MerkleProof.verify(proofG1, accused.rootPubCommitG1, leafG1)) {
            slash(complainerIndex);
            emit Debug("invalid merkle proof G1: slashed complainer");
            return;
        }
        bytes32 leafG2 = keccak256(abi.encodePacked(pubCommitG2));
        emit DebugProofG2(proofG2, accused.rootPubCommitG2, leafG2);
        if (!MerkleProof.verify(proofG2, accused.rootPubCommitG2, leafG2)) {
            slash(complainerIndex);
            emit Debug("invalid merkle proof G2: slashed complainer");
            return;
        }

        if (ecOps.pairingCheck(pubCommitG1, g2, g1, pubCommitG2)) {
            slash(complainerIndex);
            emit Debug("pairing was correct: slashed complainer");
        } else {
            slash(accusedIndex);
            emit Debug("pairing was incorrect: slashed accused");
        }
    }

    // A complaint on some private commitment. If for some reason this
    // function fails it will slash the complainer deposit! (unless some
    // unauthorized address made the transaction or the wrong phase).
    //
    // The complaint should be called when some private commitment does
    // not match to the public commitment.
    // The complainer has to publish the secret key from which its pk
    // for encryption is derived.
    function complaintPrivateCommit(
        uint16 complainerIndex,
        uint16 accusedIndex,
        uint256 complainerSk,
        uint256 encPrvCommit, // from the accused to the complainer
        bytes32[] memory proofPrvCommit,
        uint256[2][] memory pubCommitsG1 // accused's public G1 commitments
    )
        public
        checkAuthorizedSender(complainerIndex)
        notInPhase(Phase.EndFail)
        notInPhase(Phase.EndSuccess)
    {
        // TODO: a check for edge cases has to be
        // done (e.g., when no one has yet committed)

        curPhase = Phase.EndFail;
        emit PhaseChange(Phase.EndFail);


        Participant storage accused = participants[accusedIndex];
        if(!accused.isCommitted) {
            slash(complainerIndex);
            emit Debug("accused has not commited yet: slashed complainer");
            return;
        }

        bytes32 leaf = keccak256(abi.encodePacked(encPrvCommit));
        emit DebugProofG2(proofPrvCommit, accused.rootPrvCommit, leaf);
        if (!MerkleProof.verify(proofPrvCommit, accused.rootPrvCommit, leaf)) {
            slash(complainerIndex);
            emit Debug("invalid merkle proof prv: slashed complainer");
            return;
        }

        if(!ecOps.isEqualPoints(participants[complainerIndex].encPk, ecOps.ecmul(g1, complainerSk))) {
            slash(complainerIndex);
            emit Debug("invalid sk: slashed complainer");
            return;
        }

        if (getMerkleRoot(pubCommitsG1) != accused.rootPubCommitG1) {
            slash(complainerIndex);
            emit Debug("invalid G1 commits: merkle roots do not match: slashed complainer");
            return;
        }


        uint256 prvCommit = uint256(decrypt(accused.encPk, complainerSk, bytes32(encPrvCommit)));
        if (isPrvMatchPubCommit(complainerIndex, prvCommit, pubCommitsG1)) {
            emit Debug("correct pairing: slashed complainer");
            slash(complainerIndex);
        }
        else {
            emit Debug("incorrect pairing: slashed accused");
            slash(accusedIndex);
        }
    }

    function isPrvMatchPubCommit(
        uint16 complainerIndex,
        uint256 prvCommit,
        uint256[2][] memory pubCommitsG1
    )
        internal
        view
        returns (bool isMatch)
    {
        uint256[2] memory temp;
        uint256[2] memory RHS;
        uint256[2] memory LHS = ecOps.ecmul(g1, prvCommit);

        for(uint16 i = 0; i < t+1; i++) {
            temp = ecOps.ecmul(pubCommitsG1[i], complainerIndex**i);
            if(i == 0) {
                RHS = temp;
            }
            else {
                RHS = ecOps.ecadd(RHS, temp);
            }
        }

        return ecOps.isEqualPoints(LHS, RHS);
    }

    function getMerkleRoot(uint256[2][] memory pubCommitsG1) internal pure returns (bytes32) {
        bytes32[] memory leaves = new bytes32[](pubCommitsG1.length);
        for (uint256 i = 0; i < pubCommitsG1.length; i++) {
            leaves[i] = keccak256(abi.encodePacked(pubCommitsG1[i]));
        }
        return MerkleTree.root(leaves);
    }

    function checkKeccak(bytes32[2] memory b) public pure returns (bytes32) {        
            return keccak256(abi.encodePacked(b));
    }

    function leaves(bytes memory byts) public pure returns (bytes32[] memory) {
        bytes32[] memory leaves = new bytes32[](byts.length);
        for (uint256 i = 0; i < byts.length; i++) {
            leaves[i] = keccak256(abi.encodePacked(byts[i]));
        }
        return leaves;
    }

    function getMerkleRootTest(bytes memory byts) public pure returns (bytes32) {
        bytes32[] memory leaves = new bytes32[](byts.length);
        for (uint256 i = 0; i < byts.length; i++) {
            leaves[i] = keccak256(abi.encodePacked(byts[i]));
        }
        return MerkleTree.root(leaves);
    }

    function checkMerkleProof(bytes32[] memory proof, bytes32 root, bytes32 leaf) public returns(bool) {
        return MerkleProof.verify(proof, root, leaf);
    }

    function complaintVK(uint16 complainerIdx)
        public
        checkAuthorizedSender(complainerIdx)
        inPhase(Phase.EndSuccess)
    {
        uint256[2] memory expectedVk = calculateGroupPK();
        if (expectedVk[0] != VK[0] && expectedVk[1] != VK[1]) {
            slash(finalizerIdx);
        } else {
            slash(complainerIdx);
        }
    }

    // Divides the deposited balance in the contract between
    // the enrolled participants except for the participant
    // with the slashedIndex. Send slashedIndex = 0 in order
    // to divide it between all the participants (no slashing).
    function slash(uint16 slashedIndex) private {

        uint16 nParticipants = curN;
        uint256 amount;
        if (slashedIndex == 0) {
            amount = address(this).balance/nParticipants;
        }
        else {
            amount = address(this).balance/(nParticipants-1);
        }

        for (uint16 i = 1; i < (nParticipants+1); i++) {
            if (i != slashedIndex) {
                participants[i].ethPk.transfer(amount);
            }
        }
    }


    // Divides the deposited balance in the contract between
    // all the committed paricipants.
    function slashUncommitted() private {

        uint16 nParticipants = curN;
        uint16 committedNum = nParticipants - curNumCommittedLeft;
        uint256 amount = address(this).balance/committedNum;

        for (uint16 i = 1; i < (nParticipants+1); i++) {
            Participant memory part = participants[i];

            if (part.isCommitted) {
                part.ethPk.transfer(amount);
            }
        }
    }
    
    
    function decrypt(uint256[2] memory encrypterPk, uint256 decrypterSk, bytes32 encData)
        internal view
        returns(bytes32 decryptedData)
    {
        bytes32 secret = keccak256(abi.encodePacked(ecOps.ecmul(encrypterPk, decrypterSk)));
        return encData^secret;
    }


////////////////////////////////////////////////////////////////////////////

    function getParticipantPkEnc(uint16 participantIndex)
        external
        view    
        returns(uint256[2] memory encPk)
    {
        return participants[participantIndex].encPk;
    }

    // function getParticipantPubCommitG1(uint16 participantIndex, uint16 coefIndex) 
    //     view
    //     external 
    //     returns(uint256[2] memory publicCommitmentsG1)
    // {
    //     return participants[participantIndex].publicCommitmentsG1[coefIndex];
    // }

    // function getParticipantPubCommitG2(uint16 participantIndex, uint16 coefIndex) 
    //     view
    //     external 
    //     returns(uint256[4] memory publicCommitmentsG2)
    // {
    //     return participants[participantIndex].publicCommitmentsG2[coefIndex];
    // }

    // function getParticipantPrvCommit(uint16 participantIndex, uint16 committedToIndex) 
    //     view
    //     external 
    //     returns(uint256 encPrivateCommitments)
    // {
    //     return participants[participantIndex].encPrivateCommitments[committedToIndex];
    // }
    function getParticipantYG1(uint16 participantIndex)
        external
        view
        returns(uint256[2] memory yG1)
    {
        return participants[participantIndex].yG1;
    }

    function getParticipantIsCommitted(uint16 participantIndex)
        external
        view
        returns(bool isCommitted)
    {
        return participants[participantIndex].isCommitted;
    }

    // IConsensusCallback implementation

    // Consensus is used to determine if a complaint about a member's
    // invalid IPFS hash was invalid indeed. Success means that others
    // agree that the accused was faulty.
    function onConsensusSuccess(bytes calldata payload) external {
        require(_cons2Idx[msg.sender] != uint16(0), "invalid caller");

        slash(bytesToUint16(payload));
        curPhase = Phase.EndFail;

        emit PhaseChange(Phase.EndFail);
    }

    function onConsensusFailure(bytes calldata payload) external {
        uint16 complainerIdx = _cons2Idx[msg.sender];
        require(complainerIdx != uint16(0), "invalid caller");

        slash(complainerIdx);
        curPhase = Phase.EndFail;

        emit PhaseChange(Phase.EndFail);
    }

    function threshold() external view returns(uint256) {
        return t;
    }

    function isAuthorized(address sender) external view returns(bool) {
        return _group.isMember(sender);
    }

    // Utility

    function bytesToUint16(bytes memory bys) private pure returns (uint16 integer) {
        assembly {
            integer := mload(add(bys,2))
        }
    }

    function uint16ToBytes(uint16 x) public returns (bytes memory b) {
        b = new bytes(2);
        assembly {
            mstore(add(b, 2), x)
        }
    }
}


/**
 Test parameters:

    n=2
    t=1

 coefficients:
    a0) 54379457673493
    a1) 23950433293405

    b0) 453845345602931234235
    b1) 976507650679506234134

    
 public commitments:
    a0)
    1368041971066725411361239018179403078339688804905262551154469895335979601856
    1618821492510491564023544834479645350362276877645830361512548330678288690656
    a1)
    2631817276443378672842587294280308402376367224288772184477397977654447923479
    10839063031804877909681801875549944362615153185887194276974645822919437293936
    
    
    

    b0)
    13557179362105442634413454166511695479402464592547795407903173557675549038583
    14036788543633373773860064791695546493243519155298095713201690292908488603901
    b1)
    1410561832783565967033505993299263011778160659525151177822679323808323926727
    13048336431799703732972170157397576895462246813007390422018149570212241640252

    

 

 sks for decryption:
    a)9163450 (0x8bd2ba)
    b)197435619 (0xbc4a0e3)

 corresponding pks:
    a)
    8010568142108063162131985722241290538226183630284938005062991216917635380726
    19057704389966790842586953303181803223520621589676400762135428513026389437262
    b)
    20560911457302142627844632840089789900964525435158247870992611991075785849245
    6050521612570208504883690818928114121824160212045800551636642965355310285042

 private commitments:
    fa(2)
    102280324260303
    
    fb(1)
    1430352996282437468369

private commitmente encrypted:
    fa(2)
    0x492cb4e02f3d22db552acd7d0d37ac3813a17bb0f62bbf314443cb5d4dece465
    
    fb(1)
    0x492cb4e02f3d22db552acd7d0d37ac3813a17bb0f62bbf7cce6131901dd9c57b

Group PK:
    5837875810486042432133468170632608166601746087751836698563291230806091472694,
    20890074343725101418521449931290202203738610819454311822513791170653396069990

    ## Join
    ["0x11b5d2263b698dd637fb356ea748350b072265cf1acfaf374201f8e99c5bb5f6","0x2a2247476997f4e72285cc8adc57bb0350a105d8f109e523836fc7611d8deb4e"]
    ["0x2d75104069619e845ea0f055105e3adb22f07fe1206c093880b9fee9942cb99d","0xd60794fcd581fed59e19e802dcc263a5d53f6a04ddebba96a77745474f700f2"]

    ## Commit

    1,
    ["0x30648c8ef4e8e38d2db668db8a4cab5513343aad935530559090e8a51354fc0",
    "0x39438725e6ce47a9b49d4a0b2d90e1cee07d3d7e9a44adb9c0a3cf84078ade0",
    "0x5d18e484aeddc886ba162e2fa4bf8bcc125d32230a3fbea6e39ef74de3d6117",
    "0x17f6b138a7105622c493ac45d228e9c858544c47227f27a548942c2f01d59970"],
    [1,2,3,4,5,6,7,8],
    ["0x0000000000000000000000000000000000000000000000000000000000000000","0x492cb4e02f3d22db552acd7d0d37ac3813a17bb0f62bbf314443cb5d4dece465"]

    2,
    ["0x1df91772c249f1b2a7e539242ed9eb60e1475f159a376614e91de79c644097f7",
    "0x1f088a7004f9c9035af5f4686a5494f576415da8de528c40a67702c5399338fd",
    "0x31e598642c78a683eedf66cf7cd4a35a3dd5b5fd8ea947a1c53ab867154fac7",
    "0x1cd918c17d9ea92a1a3efb8a999d577d06058a1b205e99769bdc06b6686c8b3c"],
    [1,2,3,4,5,6,7,8],
    ["0x492cb4e02f3d22db552acd7d0d37ac3813a17bb0f62bbf7cce6131901dd9c57b","0x0000000000000000000000000000000000000000000000000000000000000000"]


 */