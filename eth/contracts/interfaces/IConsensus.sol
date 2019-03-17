pragma solidity ^0.5.0;

interface IConsensus {
    function propose(bytes32 digest, bytes calldata payload) external;

    function invalidate() external;

    function getProposer() external returns(address);
}
