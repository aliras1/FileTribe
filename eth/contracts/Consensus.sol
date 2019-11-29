pragma solidity ^0.5.0;

import "./common/Ownable.sol";
import "./interfaces/IAccount.sol";
import "./interfaces/IConsensusCallback.sol";
import "./interfaces/IConsensus.sol";

contract Consensus is IConsensus {
    IAccount private _proposer;
    IConsensusCallback private _callback;
    bytes _payload;
    mapping(bytes32 => bool) _hasVoted;
    uint256 _numApprovals;
    uint256 _numDeclinations;
    uint256 _id;

    event Debug(uint256 state);
    event DebugCons(address msg);

    constructor (IAccount proposer, IConsensusCallback callback) public {
        _proposer = proposer;
        _callback = callback;
    }

    modifier isAuthorized() {
        require(_callback.isAuthorized(msg.sender), "user is not member of group");
        _;
    }

    function propose(bytes calldata payload, uint256 id) external {
        require(msg.sender == address(_callback), "msg.sender is not group owner");

        _payload = payload;
        _id = id;
        _numApprovals = 1;

        bytes32 key = keccak256(abi.encodePacked(_id, Ownable(address(_proposer)).owner()));
        _hasVoted[key] = true;
    }

    function approve() public isAuthorized {
        bytes32 key = keccak256(abi.encodePacked(_id, msg.sender));
        if (_hasVoted[key]) {
            return;
        }

        _hasVoted[key] = true;

        if (++_numApprovals > _callback.threshold()) {
            _callback.onConsensusSuccess(_payload);
        }
    }

    function decline() public isAuthorized {
        bytes32 key = keccak256(abi.encodePacked(_id, msg.sender));
        if (_hasVoted[key]) {
            return;
        }

        _hasVoted[key] = true;

        if (++_numDeclinations > _callback.threshold()) {
            _callback.onConsensusFailure(_payload);
        }
    }

    function getId() external view returns(uint256) {
        return _id;
    }

    function getProposer() external view returns(IAccount) {
        return _proposer;
    }    

    function payload() public view returns(bytes memory) {
        return _payload;
    }

    function proposer() public view returns(IAccount) {
        return _proposer;
    }
}
