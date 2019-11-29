pragma solidity ^0.5.11;

import "./../Consensus.sol";
import "./../interfaces/IConsensus.sol";
import "./../interfaces/IGroup.sol";
import "./../common/Ownable.sol";

contract FakeGroup is Ownable, IGroup {

    uint256 _threshold;
    IConsensus _consensus;

    constructor() public Ownable(msg.sender) {}
    
    function isMember(address owner) external view returns(bool) {
        return true;
    }

    function onChangeIpfsHashConsensus(bytes calldata payload) external {}

    function threshold() external view returns(uint256) {
        return _threshold;
    }
    
    function setThreshold(uint256 t) public {
        _threshold = t;
    }
    
    function createConsensus(IAccount acc) public returns(address) {
        _consensus = new Consensus(acc, this);
        return address(_consensus);
    }
    
    function propose() public {
        _consensus.propose("0x0012", 345);
    }
}