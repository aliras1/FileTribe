pragma solidity ^0.5.0;

import "./IAccount.sol";

interface IConsensus {
    function propose(bytes calldata payload, uint256 counter) external;

    function getProposer() external view returns(IAccount);

    function getId() external view returns(uint256);
}
