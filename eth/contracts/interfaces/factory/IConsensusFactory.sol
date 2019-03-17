pragma solidity ^0.5.0;

interface IConsensusFactory {

    function create(address proposer, address group) external returns(address);
}
