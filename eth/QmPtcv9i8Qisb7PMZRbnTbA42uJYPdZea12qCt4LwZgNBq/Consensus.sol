pragma solidity ^0.5.0;

import "./common/Ownable.sol";
import "./interfaces/IAccount.sol";
import "./interfaces/IGroup.sol";
import "./interfaces/IConsensus.sol";

contract Consensus is IConsensus {
    IAccount private _proposer;
    IGroup private _group;
    bytes _payload;
    mapping(bytes32 => bool) _memberApproved;
    uint256 _numApprovals;
    uint256 _id;

    event Debug(uint256 state);
    event DebugCons(address msg);

    constructor (IAccount proposer, IGroup group) public {
        _proposer = proposer;
        _group = group;
    }

    function propose(bytes calldata payload, uint256 id) external {
        require(msg.sender == address(_group), "msg.sender is not group owner");

        _payload = payload;
        _id = id;
        _numApprovals = 1;

        bytes32 key = keccak256(abi.encodePacked(_id, Ownable(address(_proposer)).owner()));
        _memberApproved[key] = true;
    }

    function approve() public onlyMembers {
        bytes32 key = keccak256(abi.encodePacked(_id, msg.sender));
        if (_memberApproved[key]) {
            return;
        }

        _memberApproved[key] = true;

        if (++_numApprovals > IGroup(_group).threshold()) {
            IGroup(_group).onChangeIpfsHashConsensus(_payload);
        }
    }

    function getId() external view returns(uint256) {
        return _id;
    }

    function getProposer() external view returns(IAccount) {
        return _proposer;
    }

    modifier onlyMembers() {
        require(IGroup(_group).isMember(msg.sender), "user is not member of group");
        _;
    }

    function payload() public view returns(bytes memory) {
        return _payload;
    }

    function proposer() public view returns(IAccount) {
        return _proposer;
    }
}
