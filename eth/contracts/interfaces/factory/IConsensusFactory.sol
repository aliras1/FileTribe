pragma solidity ^0.5.0;

import "../IConsensus.sol";

interface IConsensusFactory {

    function create(IConsensus.Type cType, address proposer, address group) external returns(address);
}
