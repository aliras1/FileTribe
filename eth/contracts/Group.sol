pragma solidity ^0.5.0;

import "./interfaces/IAccount.sol";
import "./interfaces/IFileTribeDApp.sol";
import "./interfaces/IConsensus.sol";
import "./common/Ownable.sol";

contract Group is Ownable {
    struct Member {
        address account;
        address consensus;
        bool canInvite;
        bool exists;
    }

    address _fileTribe;
    string public _name;
    address[] private _memberOwners;
    mapping(address => Member) _members;
    bytes public _ipfsHash; // encrypted with group key
    uint256 _currConsId;
    mapping(address => address) private _invitations; // owner -> account address

    event GroupRegistered(bytes32 id);
    event GroupLeft(bytes32 groupId, address user);
    event GroupInvitation(address from, address to, bytes32 groupId);
    event GroupUpdateIpfsHash(bytes32 groupId, bytes ipfsHash);
    event InvitationAccepted(address groupAddress, address account);
    event InvitationSent(address group, address account);
    event InvitationDeclined(address group, address account);
    event NewConsensus(address group, address consensus, uint256 id);
    event IpfsHashChanged(address group, bytes ipfsHash, address proposer, uint256 id);
    event MemberLeft(address group, address account);
    event Debug(uint256 msg);

    constructor (
        address fileTribe,
        address account,
        string memory name,
        bytes memory ipfsHash)
    public Ownable(account) {
        address owner = IAccount(account).owner();

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

    function commit(bytes memory newIpfsHash) public onlyMembers {
        if (_memberOwners.length == 1) {
            emit IpfsHashChanged(address(this), newIpfsHash, IFileTribeDApp(_fileTribe).getAccount(msg.sender), _currConsId);

            _ipfsHash = newIpfsHash;

            return;
        }

        address cons = _members[msg.sender].consensus;
        IConsensus(cons).propose(newIpfsHash, _currConsId + 1);

        emit NewConsensus(address(this), cons, _currConsId + 1);
    }

    function onChangeIpfsHashConsensus(bytes calldata payload) external {
        uint256 id = IConsensus(msg.sender).getId();
        require(id > _currConsId, "Consensus expired!");

        address proposer = IConsensus(msg.sender).getProposer();
        address proposerOwner = IAccount(proposer).owner();

        require(_members[proposerOwner].consensus == msg.sender, "Consensus does not belong to the group!");

        _ipfsHash = payload;
        _currConsId = id;

        emit IpfsHashChanged(address(this), payload, proposer, _currConsId);
    }

    function leave() onlyMembers public {
        address account = _members[msg.sender].account;
        IAccount(account).onGroupLeft(address(this));
        _members[msg.sender].exists = false;

        uint256 i;
        for (i = 0; i < _memberOwners.length; i++) {
            if (_memberOwners[i] == msg.sender) {
                _memberOwners[i] = _memberOwners[_memberOwners.length - 1];
                _memberOwners.length--;

                break;
            }
        }

        emit MemberLeft(address(this), account);
    }

    function kick(address memberOwner) public onlyMembers {
        address account = _members[memberOwner].account;
        IAccount(account).onGroupLeft(address(this));
        _members[memberOwner].exists = false;

        uint256 i;
        for (i = 0; i < _memberOwners.length; i++){
            if (_memberOwners[i] == memberOwner) {
                _memberOwners[i] = _memberOwners[_memberOwners.length - 1];
                _memberOwners.length--;
                break;
            }
        }

        emit MemberLeft(address(this), account);
    }

    function invite(address account) public onlyMembers {
        address accountOwner = IAccount(account).owner();

        require(!_members[accountOwner].exists, "The user to be invited is already a member");
        require(_invitations[accountOwner] == address(0), "account has already been invited");

        IAccount(account).invite();

        emit InvitationSent(address(this), account);

        _invitations[accountOwner] = account;
    }

    function join() public {
        address account = _invitations[msg.sender];
        require(account != address(0), "account was not invited");

        _memberOwners.push(msg.sender);

        _members[msg.sender].account = account;
        _members[msg.sender].consensus = IFileTribeDApp(_fileTribe).createConsensus(account);
        _members[msg.sender].exists = true;

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

    function getConsensus(address owner) public view returns(address) {
        return _members[owner].consensus;
    }

    function memberOwners() public view returns(address[] memory) {
        return _memberOwners;
    }

    function threshold() public view returns(uint256) {
        return _memberOwners.length / 2;
    }
}