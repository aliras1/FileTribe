pragma solidity ^0.5.0;

import "../Account.sol";
import "../common/Ownable.sol";

contract AccountFactory is Ownable {
    address _parent;

    constructor() public Ownable(msg.sender) {

    }

    function create(address owner, string calldata name, string calldata ipfsId, bytes32 key) external returns(address) {
        return address(new Account(_parent, owner, name, ipfsId, key));
    }

    function setParent(address parent) public onlyOwner {
        _parent = parent;
    }
}
