pragma solidity ^0.5.0;

import "../Consensus.sol";
import "../interfaces/IConsensus.sol";
import "../common/Ownable.sol";

contract ConsensusFactory is Ownable {
    address _parent;

    constructor() Ownable(msg.sender) public {

    }

    function create(
        IConsensus.Type cType,
        address proposer,
        address group,
        bytes32 digest,
        bytes calldata payload)
    external returns(address) {
        return address(new Consensus(cType, proposer, group, digest, payload));
    }

    function setParent(address parent) public onlyOwner {
        _parent = parent;
    }
}
