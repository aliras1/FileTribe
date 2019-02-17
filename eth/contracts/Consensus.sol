pragma solidity ^0.5.0;

import "./interfaces/IGroup.sol";
import "./interfaces/IConsensus.sol";

contract Consensus {
    IConsensus.Type private _type;
    address private _proposer;
    address private _group;
    bytes32 _digest;
    bytes _payload;
    address[] _membersThatApproved;

    event Debug(uint256 state);
    event DebugCons(bytes msg);

    constructor (
        IConsensus.Type cType,
        address proposer,
        address group,
        bytes32 digest,
        bytes memory payload)
    public {
        _type = cType;
        _proposer = proposer;
        _group = group;
        _digest = digest;
        _payload = payload;

        _membersThatApproved.push(proposer);
    }

    function approve(bytes32 r, bytes32 s, uint8 v) public onlyMembers {
        require(memberNotApprovedYet(msg.sender), "member already approved");
        // require(verify(msg.sender, _digest, v, r, s), "invalid approval: invalid signature");

        approveInternal(msg.sender);
    }

    function approveExternal(address sender, bytes32 r, bytes32 s, uint8 v) external {
        require(IGroup(_group).isMember(msg.sender), "user is not member of group");
        require(memberNotApprovedYet(sender), "member already approved");
        // require(verify(msg.sender, _digest, v, r, s), "invalid approval: invalid signature");

        approveInternal(sender);
    }

    function approveInternal(address sender) private {
        if (_membersThatApproved.length + 1 > IGroup(_group).threshold()) {
            emit Debug(101);

            if (_type == IConsensus.Type.IPFS_HASH) {
                IGroup(_group).onChangeIpfsHashConsensus(_payload);
            } else if (_type == IConsensus.Type.KEY) {
                IGroup(_group).onChangeKeyConsensus(_payload);
            } else if (_type == IConsensus.Type.LEADER) {
                IGroup(_group).onChangeLeaderConsensus();
            }
        } else {
            emit Debug(102);
            _membersThatApproved.push(sender);
        }
    }

    function verify(address addr, bytes32 hash, uint8 v, bytes32 r, bytes32 s) internal pure returns(bool) {
        return ecrecover(hash, v, r, s) == addr;
    }

    function memberNotApprovedYet(address member) private view returns(bool) {
        for (uint256 i = 0; i < _membersThatApproved.length; i++) {
            if (_membersThatApproved[i] == member) {
                return false;
            }
        }

        return true;
    }

    modifier onlyMembers() {
        require(IGroup(_group).isMember(msg.sender), "user is not member of group");
        _;
    }

    function ctype() public view returns(IConsensus.Type) {
        return _type;
    }

    function digest() public view returns(bytes32) {
        return _digest;
    }

    function payload() public view returns(bytes memory) {
        return _payload;
    }

    function proposer() public view returns(address) {
        return _proposer;
    }
}
