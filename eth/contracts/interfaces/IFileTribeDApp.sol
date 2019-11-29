pragma solidity ^0.5.0;

import "./IGroup.sol";
import "./IAccount.sol";
import "./IConsensus.sol";

interface IFileTribeDApp {
    function createGroup(string calldata name) external returns(IGroup group);

    function createConsensus(IAccount proposer) external returns(IConsensus consensus);

    // function onInvitationAccepted(address group) external;

    // function onInvitationDeclined() external;

    function getAccountOf(address owner) external view returns (IAccount);
    
    function createDkg() external returns(address);
}
