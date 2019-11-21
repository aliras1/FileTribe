pragma solidity ^0.5.0;

import "../IAccount.sol";
import "../IConsensus.sol";
import "../IGroup.sol";

interface IConsensusFactory {

    function create(IAccount proposer, IGroup group) external returns(IConsensus);
}
