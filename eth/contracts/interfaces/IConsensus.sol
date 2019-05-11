pragma solidity ^0.5.0;

interface IConsensus {
    function propose(bytes calldata payload, uint256 counter) external;

    function getProposer() external returns(address);

    function getId() external returns(uint256);
}
