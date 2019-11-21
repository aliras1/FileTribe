pragma solidity ^0.5.0;

import "../interfaces/IAccount.sol";
import "../interfaces/IFileTribeDApp.sol";
import "../interfaces/factory/IAccountFactory.sol";
import "../Account.sol";
import "../common/Ownable.sol";

contract AccountFactory is Ownable, IAccountFactory {
    IFileTribeDApp _parent;

    constructor() public Ownable(msg.sender) {

    }

    function create(address owner, string calldata name, string calldata ipfsId, bytes32 key) external returns(IAccount) {
        return new Account(_parent, owner, name, ipfsId, key);
    }

    function setParent(IFileTribeDApp parent) public onlyOwner {
        _parent = parent;
    }
}
