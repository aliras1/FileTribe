pragma solidity ^0.5.0;

import "../interfaces/IFileTribeDApp.sol";
import "../interfaces/factory/IGroupFactory.sol";
import "../Group.sol";
import "../common/Ownable.sol";

contract GroupFactory is Ownable, IGroupFactory {
    IFileTribeDApp _parent;

    constructor() Ownable(msg.sender) public {

    }

    function create(IAccount account, string calldata name) external returns(IGroup) {
        return new Group(_parent, account, name, new bytes(0x00), [uint256(0), uint256(0), uint256(0), uint256(0)]);
    }

    function setParent(IFileTribeDApp parent) public onlyOwner {
        _parent = parent;
    }
}
