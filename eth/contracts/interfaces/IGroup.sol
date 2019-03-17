pragma solidity ^0.5.0;

interface IGroup {
    function isMember(address addr) external returns(bool);

    function onChangeIpfsHashConsensus(bytes calldata payload) external;

    function onInvitationAccepted(address account) external;

    function onInvitationDeclined(address account) external;

    function threshold() external view returns(uint256);

}
