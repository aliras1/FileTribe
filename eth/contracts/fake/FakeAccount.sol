pragma solidity ^0.5.0;

import "./../interfaces/IGroup.sol";
import "./../interfaces/IAccount.sol";
import "./../common/Ownable.sol";

contract FakeAccount is Ownable, IAccount {

    constructor() public Ownable(msg.sender) {}

    function invite() external {}

    function onInvitationAccepted() external {}

    function onInvitationDeclined() external {}

    function onGroupLeft(IGroup group) external {}
}