pragma solidity ^0.5.0;

import "./interfaces/IAccount.sol";
import "./interfaces/IFileTribeDApp.sol";
import "./interfaces/IConsensus.sol";
import "./common/Ownable.sol";

contract Group is Ownable {
    address _fileTribe;
    string private _name;
    address[] private _members;
    bytes private _ipfsHash; // encrypted with group key
    mapping(address => bool) private _canInvite;
    uint256 _leaderIdx;
    address[] private _consensuses;
    mapping(address => uint256) private _memberToIdx;
    mapping(address => address) private _invitations; // owner -> account address

    event GroupRegistered(bytes32 id);
    event GroupLeft(bytes32 groupId, address user);
    event GroupInvitation(address from, address to, bytes32 groupId);
    event GroupUpdateIpfsHash(bytes32 groupId, bytes ipfsHash);
    event InvitationAccepted(address groupAddress, address account);
    event InvitationSent(address group, address account);
    event InvitationDeclined(address group, address account);
    event NewConsensus(address group, address consensus);
    event IpfsHashChanged(address group, bytes ipfsHash, address proposer);
    event MemberLeft(address group, address account);
    event Debug(int msg);

    constructor (
        address fileTribe,
        address account,
        string memory name,
        bytes memory ipfsHash)
    public Ownable(account) {
        _fileTribe = fileTribe;
        _name = name;
        _ipfsHash = ipfsHash;
        _members.push(account);
        _consensuses.push(IFileTribeDApp(_fileTribe).createConsensus(account));
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

    function changeIpfsHash(bytes memory newIpfsHash) public onlyMembers {
        bytes32 digest = keccak256(abi.encodePacked(_ipfsHash, newIpfsHash));

        if (_members.length == 1) {
            emit IpfsHashChanged(address(this), newIpfsHash, IFileTribeDApp(_fileTribe).getAccount(msg.sender));

            _ipfsHash = newIpfsHash;

            return;
        }

        uint256 idx = _memberToIdx[msg.sender];
        IConsensus(_consensuses[idx]).propose(digest, newIpfsHash);

        emit NewConsensus(address(this), _consensuses[idx]);
    }

    function onChangeIpfsHashConsensus(bytes calldata payload) external {
        address proposer = IConsensus(msg.sender).getProposer();
        address proposerOwner = IAccount(proposer).owner();

        require(_consensuses[_memberToIdx[proposerOwner]] != address(0), "msg.sender is no group consensus");

        for (uint256 i = 0; i < _consensuses.length; i++) {
            IConsensus(_consensuses[i]).invalidate();
        }

        emit IpfsHashChanged(address(this), payload, proposer);

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

        address memberLeft = _members[i];

        _members[i] = _members[_members.length - 1];
        _consensuses[i] = _consensuses[_consensuses.length - 1];
        _members.length--;

        emit MemberLeft(address(this), memberLeft);
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
        _consensuses[i] = _consensuses[_consensuses.length - 1];
        _members.length--;

        emit MemberLeft(address(this), member);
    }

    function name() public view returns(string memory) {
        return _name;
    }

    function members() public view returns(address[] memory) {
        return _members;
    }

    function invite(address account) public onlyMembers {
        uint256 i;
        for (i = 0; i < _members.length; i++) {
            if (_members[i] == account) {
                revert("account is already member of the group");
            }
        }

        address accountOwner = IAccount(account).owner();
        require(_invitations[accountOwner] == address(0), "account has already been invited");

        IAccount(account).invite();

        emit InvitationSent(address(this), account);

        _invitations[accountOwner] = account;
    }

    function join() public {
        address account = _invitations[msg.sender];
        require(account != address(0), "account was not invited");

        _memberToIdx[IAccount(account).owner()] = _members.length;
        _members.push(account);

        address c = IFileTribeDApp(_fileTribe).createConsensus(account);
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

    function getConsensus(address owner) public view returns(address) {
        return _consensuses[_memberToIdx[owner]];
    }

    function threshold() public view returns(uint256) {
        return _members.length / 2;
    }

    function ipfsHash() public view returns(bytes memory) {
        return _ipfsHash;
    }
}