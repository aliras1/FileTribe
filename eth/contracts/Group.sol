pragma solidity ^0.5.0;

import "./interfaces/IAccount.sol";
import "./interfaces/IDipfshare.sol";
import "./interfaces/IConsensus.sol";
import "./common/Ownable.sol";
import "./common/Invitable.sol";
import "./common/Governable.sol";

contract Group is Ownable, Invitable, Governable {
    enum State {
        NORMAL,
        KEY_DIRTY
    }

    address _parent; // Dipfshare
    string private _name;
    address[] private _members;
    bytes private _ipfsHash; // encrypted with group key
    mapping(address => bool) private _canInvite;
    State private _state;
    uint256 _leaderIdx;
    address _keyConsensus;
    address _leaderConsensus;

    event KeyDirty(address group);
    event GroupKeyChanged(bytes32 groupId, bytes ipfsHash);
    event GroupRegistered(bytes32 id);
    event GroupLeft(bytes32 groupId, address user);
    event GroupInvitation(address from, address to, bytes32 groupId);
    event GroupUpdateIpfsHash(bytes32 groupId, bytes ipfsHash);
    event InvitationAccepted(address groupAddress, address account);
    event InvitationSent(address group, address account);
    event InvitationDeclined(address group, address account);
    event NewConsensus(address group, address consensus);
    event GroupConsensusFailed(address group, address consensus);
    event ConsensusReached(address group, bytes ipfsHash);
    event IpfsHashChanged(address group, bytes ipfsHash);
    event MemberLeft(address group, address account);
    event GroupLeaderChanged(address group, address leader);
    event ApprovedChangeLeader(address account);
    event Debug(int msg);

    constructor (
        address parent,
        address account,
        string memory name,
        bytes memory ipfsHash)
    public Ownable(account) Invitable() Governable() {
        _parent = parent;
        _name = name;
        _ipfsHash = ipfsHash;
        _state = State.NORMAL;
        _members.push(account);
    }

    modifier onlyMembers() {
        require(isMember(msg.sender));
        _;
    }

    function isMember(address addr) public returns(bool) {
        for (uint256 i = 0; i < _members.length; i++) {
            if (addr == IAccount(_members[i]).owner()) {
                return true;
            }
        }

        return false;
    }

    // onlyInDirtyState
    function changeLeader(bytes32 r, bytes32 s, uint8 v) public onlyMembers  {
        if (_leaderConsensus == address(0)) {
            emit Debug(0);
            _leaderConsensus = IDipfshare(_parent).createConsensus(
                IConsensus.Type.LEADER,
                msg.sender,
                keccak256(abi.encodePacked(leader())),
                new bytes(0x00)  // not intrested in additional data
            );
        }

        emit ApprovedChangeLeader(msg.sender);

        IConsensus(_leaderConsensus).approveExternal(msg.sender, r, s, v);
    }

    function onChangeLeaderConsensus() external {
        require(msg.sender == _leaderConsensus, "msg.sender is not leaderConsensus");

        _leaderIdx++;
        _leaderConsensus = address(0);

        emit GroupLeaderChanged(address(this), leader());
    }

    function changeKey(bytes memory newIpfsHash) public onlyMembers onlyInDirtyState {
        bytes32 digest = keccak256(abi.encodePacked(_ipfsHash, newIpfsHash));

        if (_keyConsensus == address(0)) {
            return;
        }

        _keyConsensus = IDipfshare(_parent).createConsensus(
            IConsensus.Type.KEY,
            msg.sender,
            digest,
            newIpfsHash);

        emit NewConsensus(address(this), address(_keyConsensus));
    }

    function onChangeKeyConsensus(bytes calldata payload) external {
        require(msg.sender == address(_keyConsensus), "msg.sender is not keyConsensus");

        _ipfsHash = payload;
        _state = State.NORMAL;
        _leaderIdx++;
        _leaderConsensus = address(0);

        emit ConsensusReached(address(this), payload);
    }

    function changeIpfsHash(bytes memory newIpfsHash) public onlyMembers onlyInNormalState {
        bytes32 digest = keccak256(abi.encodePacked(_ipfsHash, newIpfsHash));

        address c = IDipfshare(_parent).createConsensus(
            IConsensus.Type.IPFS_HASH,
            msg.sender,
            digest,
            newIpfsHash);

        emit NewConsensus(address(this), address(c));

        addConsensus(c);
    }

    function onChangeIpfsHashConsensus(bytes calldata payload) external {
        require(_state == State.NORMAL, "operation is not possible in the current state");
        require(isConsensus(msg.sender), "msg.sender is no consensus");

        emit IpfsHashChanged(address(this), payload);

        clearConsensuses();
        _ipfsHash = payload;
    }

    // no need for onlyMembers modifier because we have to
    // iterate over the array to get the user's index anyway
    function leave() public {
        uint256 i;

        for (i = 0; i < _members.length; i++){
            if (IAccount(_members[i]).owner() == msg.sender) {
                break;
            }
        }

        require(i < _members.length, "msg.sender is not group member");

        IAccount(_members[i]).onGroupLeft(address(this));

        _members[i] = _members[_members.length - 1];
        _members.length--;

        _state = State.KEY_DIRTY;
        clearConsensuses();

        emit KeyDirty(address(this));
    }

    function kick(address member) public onlyMembers {
        uint256 i;

        for (i = 0; i < _members.length; i++){
            if (_members[i] == member) {
                break;
            }
        }
        require(i < _members.length, "msg.sender is not group member");

        IAccount(member).onGroupLeft(address(this));

        _members[i] = _members[_members.length - 1];
        _members.length--;

        _state = State.KEY_DIRTY;
        clearConsensuses();

        emit KeyDirty(address(this));
    }

    function leader() public view returns(address account) {
        return address(_members[_leaderIdx % _members.length]);
    }

    function name() public view returns(string memory) {
        return _name;
    }

    function members() public view returns(address[] memory) {
        return _members;
    }

    function invite(address account) public onlyMembers {
        address inv = IAccount(account).invite();

        emit InvitationSent(address(this), account);

        addInvitation(inv);
    }

    function onInvitationAccepted(address account) external onlyInvitations {
        removeInvitation(msg.sender);
        _members.push(account);

        emit InvitationAccepted(address(this), account);
    }

    function onInvitationDeclined(address account) external onlyInvitations {
        removeInvitation(msg.sender);
        _members.push(account);

        emit InvitationDeclined(address(this), account);
    }

    function threshold() public view returns(uint) {
        return _members.length / 2;
    }

    function ipfsHash() public view returns(bytes memory) {
        return _ipfsHash;
    }

    modifier onlyInNormalState() {
        require(_state == State.NORMAL, "operation is not possible in this state");
        _;
    }

    modifier onlyInDirtyState() {
        require(_state == State.KEY_DIRTY, "operation is not possible in this state");
        _;
    }

    modifier onlyLeader() {
        require(msg.sender == leader(), "msg.sender is not the leader");
        _;
    }

    modifier onlyKeyConsensus() {
        require(msg.sender == address(_keyConsensus), "msg.sender is not the keyConsensus");
        _;
    }

    modifier onlyLeaderConsensus() {
        require(msg.sender == address(_leaderConsensus), "msg.sender is not the leaderConsensus");
        _;
    }
}