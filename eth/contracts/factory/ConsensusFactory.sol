pragma solidity ^0.5.0;

import "../Consensus.sol";
import "../common/Ownable.sol";

contract ConsensusFactory is Ownable {
    address _parent;

    constructor() Ownable(msg.sender) public {

    }

    function create(address proposer, address group) external returns(address) {
        return address(new Consensus(proposer, group));
    }

    function setParent(address parent) public onlyOwner {
        _parent = parent;
    }
}
