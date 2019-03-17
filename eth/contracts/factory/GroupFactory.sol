pragma solidity ^0.5.0;

import "../Group.sol";
import "../common/Ownable.sol";

contract GroupFactory is Ownable {
    address _parent;

    constructor() Ownable(msg.sender) public {

    }

    function create(address owner, string calldata name) external returns(address) {
        return address(new Group(_parent, owner, name, new bytes(0x00)));
    }

    function setParent(address parent) public onlyOwner {
        _parent = parent;
    }
}
