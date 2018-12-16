pragma solidity ^0.5.0;

import "../IConsensus.sol";

interface IConsensusFactory {

    function create(
        IConsensus.Type cType,
        address proposer,
        address group,
        bytes32 digest,
        bytes calldata payload)
    external returns(address);
}
