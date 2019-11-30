pragma solidity ^0.5.0;

import "./interfaces/IConsensusCallback.sol";
import "./interfaces/IFileTribeDApp.sol";
import "./interfaces/factory/IAccountFactory.sol";
import "./interfaces/factory/IGroupFactory.sol";
import "./interfaces/factory/IConsensusFactory.sol";
import "./common/Ownable.sol";
import "./GroupDkg.sol";

contract FileTribeDApp is Ownable, IFileTribeDApp {
    IAccountFactory private _accountFactory;
    IGroupFactory private _groupFactory;
    IConsensusFactory private _consensusFactory;
    mapping(address => IAccount) private _accounts;

    event AccountCreated(address owner, IAccount account);
    event GroupRegistered(bytes32 id);
    event Debug(int msg);
    event DebugBytes(bytes msg);

    constructor () public Ownable(msg.sender) {
    }

    function createDkg() external returns(address) {
        return address(new GroupDkg());
    }

    function setConsensusFactory(IConsensusFactory factory) public onlyOwner {
        _consensusFactory = factory;
    }

    function setAccountFactory(IAccountFactory factory) public onlyOwner {
        _accountFactory = factory;
    }

    function setGroupFactory(IGroupFactory factory) public onlyOwner {
        _groupFactory = factory;
    }

    function createAccount(string memory name, string memory ipfsPeerId, bytes32 boxingKey) public {
        require(address(_accounts[msg.sender]) == address(0), "Account already exists");

        IAccount account = _accountFactory.create(msg.sender, name, ipfsPeerId, boxingKey);
        _accounts[msg.sender] = account;

        emit AccountCreated(msg.sender, account);
    }

    function removeAccount() public {
        _accounts[msg.sender] = IAccount(address(0));
    }

    function createGroup(string calldata name) external returns(IGroup group) {
        return _groupFactory.create(IAccount(msg.sender), name);
    }

    function createConsensus(IAccount proposer) external returns(IConsensus consensus) {
        return _consensusFactory.create(proposer, IConsensusCallback(msg.sender));
    }


    //    function isUserRegistered(address id) public view returns(bool) {
    //        return _users[id].exists;
    //    }

    function getAccountOf(address owner) public view returns(IAccount) {
        return _accounts[owner];
    }
}
