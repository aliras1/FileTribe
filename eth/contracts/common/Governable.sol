pragma solidity ^0.5.0;

contract Governable {
    address[] private _pendingConsensuses;

    constructor () internal {

    }

    modifier onlyConsensuses() {
        require(isConsensus(msg.sender));
        _;
    }

    function isConsensus(address consensus) internal view returns(bool) {
        for (uint256 i = 0; i < _pendingConsensuses.length; i++) {
            if (_pendingConsensuses[i] == consensus) {
                return true;
            }
        }

        return false;
    }

    function removeConsensus(address consensus) internal {
        for (uint256 i = 0; i < _pendingConsensuses.length; i++) {
            if (consensus == _pendingConsensuses[i]) {
                _pendingConsensuses[i] = _pendingConsensuses[_pendingConsensuses.length - 1];
                _pendingConsensuses.length--;
            }
        }
    }

    function addConsensus(address consensus) internal {
        _pendingConsensuses.push(consensus);
    }

    function consensuses() public view returns(address[] memory) {
        return _pendingConsensuses;
    }

    function clearConsensuses() internal {
        _pendingConsensuses.length = 0;
    }
}
