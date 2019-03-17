pragma solidity ^0.5.0;

import "./interfaces/factory/IAccountFactory.sol";
import "./interfaces/factory/IGroupFactory.sol";
import "./interfaces/factory/IConsensusFactory.sol";
import "./common/Ownable.sol";

contract FileTribeDApp is Ownable {
    address private _accountFactory;
    address private _groupFactory;
    address private _consensusFactory;
    mapping(address => address) private _accounts;

    event KeyDirty(bytes32 groupId);
    event GroupKeyChanged(bytes32 groupId, bytes ipfsHash);
    event AccountCreated(address owner, address account);
    event GroupRegistered(bytes32 id);
    event GroupLeft(bytes32 groupId, address user);
    event GroupInvitation(address from, address to, bytes32 groupId);
    event GroupUpdateIpfsHash(bytes32 groupId, bytes ipfsHash);
    event Debug(int msg);
    event DebugBytes(bytes msg);

    constructor () public Ownable(msg.sender) {
    }

    function setConsensusFactory(address factory) public onlyOwner {
        _consensusFactory = factory;
    }

    function setAccountFactory(address factory) public onlyOwner {
        _accountFactory = factory;
    }

    function setGroupFactory(address factory) public onlyOwner {
        _groupFactory = factory;
    }

    function createAccount(string memory name, string memory ipfsPeerId, bytes32 boxingKey) public {
        require(_accounts[msg.sender] == address(0), "Account already exists");

        address acc = IAccountFactory(_accountFactory).create(msg.sender, name, ipfsPeerId, boxingKey);
        _accounts[msg.sender] = acc;

        emit AccountCreated(msg.sender, acc);
    }

    function createGroup(string calldata name) external returns(address group) {
        return IGroupFactory(_groupFactory).create(msg.sender, name);
    }

    function createConsensus(address proposer) external returns(address consensus) {
        return IConsensusFactory(_consensusFactory).create(proposer, msg.sender);
    }


    //    function isUserRegistered(address id) public view returns(bool) {
    //        return _users[id].exists;
    //    }

    function getAccount(address addr) public view returns(address) {
        return _accounts[addr];
    }
}
