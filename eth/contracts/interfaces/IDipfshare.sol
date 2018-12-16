pragma solidity ^0.5.0;

import "./IConsensus.sol";

interface IDipfshare {
    function createGroup(string calldata name) external returns(address group);

    function createConsensus(
        IConsensus.Type cType,
        address proposer,
        bytes32 digest,
        bytes calldata payload)
    external returns(address consensus);

    function onInvitationAccepted(address group) external;

    function onInvitationDeclined() external;

    function owner() external returns(address);
}
