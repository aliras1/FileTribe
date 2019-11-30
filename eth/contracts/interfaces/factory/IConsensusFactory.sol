pragma solidity ^0.5.0;

import "../IAccount.sol";
import "../IConsensus.sol";
import "../IConsensusCallback.sol";

interface IConsensusFactory {

    function create(IAccount proposer, IConsensusCallback callback) external returns(IConsensus);
}
