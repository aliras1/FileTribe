pragma solidity ^0.5.0;


interface IDipfshare {
    function createGroup(string calldata name) external returns(address group);

    function createConsensus(address proposer) external returns(address consensus);

    function onInvitationAccepted(address group) external;

    function onInvitationDeclined() external;

    function owner() external returns(address);
}
