pragma solidity ^0.5.0;

import "./interfaces/IAccount.sol";
import "./interfaces/IDipfshare.sol";
import "./interfaces/IConsensus.sol";
import "./common/Ownable.sol";

contract Group is Ownable {
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
    address[] private _consensuses;
    mapping(address => uint256) private _memberToIdx;
    address _keyConsensus;
    address _leaderConsensus;
    mapping(address => address) private _invitations; // owner -> account address

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
    public Ownable(account) {
        _parent = parent;
        _name = name;
        _ipfsHash = ipfsHash;
        _state = State.NORMAL;
        _members.push(account);
        _keyConsensus = IDipfshare(_parent).createConsensus(IConsensus.Type.KEY, address(0));
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

    function changeKey(bytes memory newIpfsHash) public onlyMembers onlyInDirtyState {
        bytes32 digest = keccak256(abi.encodePacked(_ipfsHash, newIpfsHash));

        IConsensus(_keyConsensus).propose(digest, newIpfsHash);
        IConsensus(_keyConsensus).setProposer(msg.sender);

        emit NewConsensus(address(this), address(_keyConsensus));
    }

    function onChangeKeyConsensus(bytes calldata payload) external {
        require(msg.sender == address(_keyConsensus), "msg.sender is not keyConsensus");

        _ipfsHash = payload;
        _state = State.NORMAL;

        emit ConsensusReached(address(this), payload);
    }

    function changeIpfsHash(bytes memory newIpfsHash) public onlyMembers onlyInNormalState {
        bytes32 digest = keccak256(abi.encodePacked(_ipfsHash, newIpfsHash));

        uint256 idx = _memberToIdx[msg.sender];
        IConsensus(_consensuses[idx]).propose(digest, newIpfsHash);

        emit NewConsensus(address(this), _consensuses[idx]);
    }

    function onChangeIpfsHashConsensus(bytes calldata payload) external {
        require(_state == State.NORMAL, "operation is not possible in the current state");

        uint256 i;
        for (i = 0; i < _consensuses.length; i++) {
            IConsensus(_consensuses[i]).invalidate();
            if (msg.sender == _consensuses[i]) {
                break;
            }
        }
        require(i < _consensuses.length, "msg.sender is no group consensus");
        for (; i < _consensuses.length; i++) {
            IConsensus(_consensuses[i]).invalidate();
        }

        emit IpfsHashChanged(address(this), payload);

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

        emit KeyDirty(address(this));
    }

    function onKeyDeclined(address proposer) external {
        require(msg.sender == _keyConsensus, "Only the group's Key-Consensus contract can call this function");
        require(_state == State.KEY_DIRTY, "This function can be called only in DirtyKey state");

        // TODO: ban proposer for some time

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
        IAccount(account).invite();
        address accountOwner = IAccount(account).owner();

        emit InvitationSent(address(this), account);

        _invitations[accountOwner] = account;
    }

    function join() public {
        address account = _invitations[msg.sender];
        require(account != address(0), "account was not invited");

        _memberToIdx[IAccount(account).owner()] = _members.length;
        _members.push(account);

        address c = IDipfshare(_parent).createConsensus(IConsensus.Type.IPFS_HASH, msg.sender);
        _consensuses.push(c);

        IAccount(account).onInvitationAccepted();
        _invitations[msg.sender] = address(0);

        emit InvitationAccepted(address(this), account);
    }

    function decline() public {
        address account = _invitations[msg.sender];
        require(account != address(0), "account was not invited");

        IAccount(account).onInvitationDeclined();
        _invitations[msg.sender] = address(0);

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