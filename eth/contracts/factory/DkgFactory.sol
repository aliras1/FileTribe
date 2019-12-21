pragma solidity ^0.5.0;

import "../interfaces/IFileTribeDApp.sol";
import "../interfaces/IGroup.sol";
import "../GroupDkg.sol";
import "../common/Ownable.sol";

contract DkgFactory is Ownable {
    IFileTribeDApp _parent;

    constructor() Ownable(msg.sender) public {

    }

    function create(IGroup group) external returns(address) {
        return address(new GroupDkg(group));
    }

    function setParent(IFileTribeDApp parent) public onlyOwner {
        _parent = parent;
    }
}
