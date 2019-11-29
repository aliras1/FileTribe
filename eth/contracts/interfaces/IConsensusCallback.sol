pragma solidity ^0.5.0;

interface IConsensusCallback {
    function onConsensusSuccess(bytes calldata payload) external;

    function onConsensusFailure(bytes calldata payload) external;

    function threshold() external view returns(uint256);

    function isAuthorized(address sender) external view returns(bool);
}
