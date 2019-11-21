pragma solidity ^0.5.0;

import "./IGroup.sol";

interface IAccount {
    function invite() external;

    function onInvitationAccepted() external;

    function onInvitationDeclined() external;

    function onGroupLeft(IGroup group) external;
}
