pragma solidity ^0.5.0;

import "../interfaces/factory/IConsensusFactory.sol";
import "../Consensus.sol";
import "../interfaces/IGroup.sol";
import "../interfaces/IAccount.sol";
import "../interfaces/IConsensus.sol";
import "../common/Ownable.sol";

contract ConsensusFactory is Ownable, IConsensusFactory {
    address _parent;

    constructor() Ownable(msg.sender) public {

    }

    function create(IAccount proposer, IGroup group) external returns(IConsensus) {
        return new Consensus(proposer, group);
    }

    function setParent(address parent) public onlyOwner {
        _parent = parent;
    }
}
