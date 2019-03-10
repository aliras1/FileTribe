pragma solidity ^0.5.0;

import "./interfaces/IGroup.sol";
import "./interfaces/IConsensus.sol";

contract Consensus {
    IConsensus.Type private _type;
    bool private _isActive;
    address private _proposer;
    address private _group;
    bytes32 _digest;
    bytes _payload;
    address[] _membersThatVoted;
    uint256 _numAccept;
    uint256 _numDecline;

    event Debug(uint256 state);
    event DebugCons(bytes msg);

    constructor (IConsensus.Type cType, address proposer, address group) public {
        _type = cType;
        _proposer = proposer;
        _group = group;
    }

    function propose(bytes32 digest, bytes calldata payload) external {
        require(msg.sender == _group, "msg.sender is not group owner");

        _digest = digest;
        _payload = payload;

        delete _membersThatVoted;
        _membersThatVoted.push(_proposer);
        _isActive = true;
        _numAccept = 1;
        _numDecline = 0;
    }

    function setProposer(address proposer) external {
        require(msg.sender == _group, "msg.sender is not group owner");

        _proposer = proposer;
    }

    function invalidate() external {
        require(msg.sender == _group, "msg.sender is not group owner");

        _isActive = false;
    }

    function approve(bytes32 r, bytes32 s, uint8 v) public onlyMembers onlyWhenActive {
        require(memberNotVotedYet(msg.sender), "member already voted");
        // require(verify(msg.sender, _digest, v, r, s), "invalid approval: invalid signature");
        _membersThatVoted.push(msg.sender);

        if (++_numAccept > IGroup(_group).threshold()) {
            if (_type == IConsensus.Type.IPFS_HASH) {
                IGroup(_group).onChangeIpfsHashConsensus(_payload);
            } else if (_type == IConsensus.Type.KEY) {
                IGroup(_group).onChangeKeyConsensus(_payload);
            }
        }
    }

    function decline(bytes32 r, bytes32 s, uint8 v) public onlyMembers onlyWhenActive {
        require(_type == IConsensus.Type.KEY, "Non Key-Consensus contracts can not be declined");
        require(memberNotVotedYet(msg.sender), "member already voted");
        // require(verify(msg.sender, _digest, v, r, s), "invalid approval: invalid signature");

        _membersThatVoted.push(msg.sender);

        if (++_numDecline > IGroup(_group).threshold()) {
            _isActive = false;
            IGroup(_group).onKeyDeclined(_proposer);
        }
    }

    function verify(address addr, bytes32 hash, uint8 v, bytes32 r, bytes32 s) internal pure returns(bool) {
        return ecrecover(hash, v, r, s) == addr;
    }

    function memberNotVotedYet(address member) private view returns(bool) {
        for (uint256 i = 0; i < _membersThatVoted.length; i++) {
            if (_membersThatVoted[i] == member) {
                return false;
            }
        }

        return true;
    }

    function getProposer() external returns(address) {
        return _proposer;
    }

    modifier onlyMembers() {
        require(IGroup(_group).isMember(msg.sender), "user is not member of group");
        _;
    }

    modifier onlyWhenActive() {
        require(_isActive, "consensus is not active");
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
