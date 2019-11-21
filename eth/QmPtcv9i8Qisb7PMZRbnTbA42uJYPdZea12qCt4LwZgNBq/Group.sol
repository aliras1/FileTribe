pragma solidity ^0.5.11;

import "./interfaces/IAccount.sol";
import "./interfaces/IFileTribeDApp.sol";
import "./interfaces/IConsensus.sol";
import "./common/Ownable.sol";
import "./ecOps.sol";
import "./dkg.sol";

contract Group is Ownable, IGroup {
    struct Member {
        IAccount account;
        IConsensus consensus;
        bool canInvite;
        bool exists;
    }

    IFileTribeDApp _fileTribe;
    string public _name;
    uint256[4] public _verifyKey;
    address[] private _memberOwners;
    mapping(address => Member) _members;
    bytes public _ipfsHash; // encrypted with group key
    uint256 _currConsId;
    mapping(address => IAccount) private _invitations; // owner -> account address

    // G1 generator (on the curve)
    uint256[2] public _g1 = [
        0x0000000000000000000000000000000000000000000000000000000000000001,
        0x0000000000000000000000000000000000000000000000000000000000000002
    ];
    // G2 generator (on the curve)
    uint256[4] public _g2 = [
        0x198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c2,
        0x1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed,
        0x90689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b,
        0x12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa
    ];

    event GroupUpdateIpfsHash(bytes32 groupId, bytes ipfsHash);
    event InvitationAccepted(IGroup groupAddress, IAccount account);
    event InvitationSent(IGroup group, IAccount account);
    event InvitationDeclined(IGroup group, IAccount account);
    event NewConsensus(IGroup group, IConsensus consensus, uint256 id);
    event IpfsHashChanged(IGroup group, bytes ipfsHash, IAccount proposer, uint256 id);
    event MemberLeft(IGroup group, IAccount account);
    event Debug(uint256 msg);

    constructor (
        IFileTribeDApp fileTribe,
        IAccount account,
        string memory name,
        bytes memory ipfsHash,
        uint256[4] memory verifyKey)
    public Ownable(address(account)) {
        address owner = Ownable(address(account)).owner(); // in version 0.5.11 interfaces cannot inherite, hence the back and forth casting

        _verifyKey = verifyKey;
        _fileTribe = fileTribe;
        _name = name;
        _ipfsHash = ipfsHash;
        _currConsId = 0;
        _memberOwners.push(owner);

        _members[owner].account = account;
        _members[owner].consensus = IFileTribeDApp(_fileTribe).createConsensus(account);
        _members[owner].exists = true;
    }

    modifier onlyMembers() {
        require(isMember(msg.sender));
        _;
    }

    function isMember(address owner) public view returns(bool) {
        return _members[owner].exists;
    }

    function commitWithGroupSig(bytes memory newIpfsHash, uint256[2] memory sig) public onlyMembers {
        uint256[2] memory hash = ecOps.hashToG1(newIpfsHash);

        require(ecOps.pairingCheck(hash, _verifyKey, sig, _g2), "invalid signature");

        emit IpfsHashChanged(this, newIpfsHash, _fileTribe.getAccountOf(msg.sender), _currConsId);
        _ipfsHash = newIpfsHash;
    }

    function commit(bytes memory newIpfsHash) public onlyMembers {
        if (_memberOwners.length == 1) {
            emit IpfsHashChanged(this, newIpfsHash, _fileTribe.getAccountOf(msg.sender), _currConsId);

            _ipfsHash = newIpfsHash;

            return;
        }

        IConsensus cons = _members[msg.sender].consensus;
        IConsensus(cons).propose(newIpfsHash, _currConsId + 1);

        emit NewConsensus(this, cons, _currConsId + 1);
    }

    function onChangeIpfsHashConsensus(bytes calldata payload) external {
        uint256 id = IConsensus(msg.sender).getId();
        require(id > _currConsId, "Consensus expired!");

        IAccount proposer = IConsensus(msg.sender).getProposer();
        address proposerOwner = Ownable(address(proposer)).owner();

        require(address(_members[proposerOwner].consensus) == msg.sender, "Consensus does not belong to the group!");

        _ipfsHash = payload;
        _currConsId = id;

        emit IpfsHashChanged(this, payload, proposer, _currConsId);
    }

    function leave() public onlyMembers {
        IAccount account = _members[msg.sender].account;
        IAccount(account).onGroupLeft(this);
        _members[msg.sender].exists = false;

        uint256 i;
        for (i = 0; i < _memberOwners.length; i++) {
            if (_memberOwners[i] == msg.sender) {
                _memberOwners[i] = _memberOwners[_memberOwners.length - 1];
                _memberOwners.length--;

                break;
            }
        }

        emit MemberLeft(this, account);
    }

    function kick(address memberOwner) public onlyMembers {
        IAccount account = _members[memberOwner].account;
        IAccount(account).onGroupLeft(this);
        _members[memberOwner].exists = false;

        uint256 i;
        for (i = 0; i < _memberOwners.length; i++){
            if (_memberOwners[i] == memberOwner) {
                _memberOwners[i] = _memberOwners[_memberOwners.length - 1];
                _memberOwners.length--;
                break;
            }
        }

        emit MemberLeft(this, account);
    }

    function invite(IAccount account) public onlyMembers {
        address accountOwner = Ownable(address(account)).owner();

        require(!_members[accountOwner].exists, "The user to be invited is already a member");
        require(address(_invitations[accountOwner]) == address(0), "account has already been invited");

        account.invite();

        emit InvitationSent(this, account);

        _invitations[accountOwner] = account;
    }

    function join() public {
        IAccount account = _invitations[msg.sender];
        require(address(account) != address(0), "account was not invited");

        _memberOwners.push(msg.sender);

        _members[msg.sender].account = account;
        _members[msg.sender].consensus = IFileTribeDApp(_fileTribe).createConsensus(account);
        _members[msg.sender].exists = true;

        IAccount(account).onInvitationAccepted();
        _invitations[msg.sender] = IAccount(address(0));

        emit InvitationAccepted(this, account);
    }

    function decline() public {
        IAccount account = _invitations[msg.sender];
        require(address(account) != address(0), "account was not invited");

        IAccount(account).onInvitationDeclined();
        _invitations[msg.sender] = IAccount(address(0));

        emit InvitationDeclined(this, account);
    }

    function getConsensus(address owner) public view returns(IConsensus) {
        return _members[owner].consensus;
    }

    function memberOwners() public view returns(address[] memory) {
        return _memberOwners;
    }

    function threshold() public view returns(uint256) {
        return _memberOwners.length / 2;
    }
}