pragma solidity ^0.5.0;

interface IAccount {
    function invite() external;

    function onInvitationAccepted() external;

    function onInvitationDeclined() external;

    function onGroupLeft(address group) external;

    function owner() external returns(address);
}
