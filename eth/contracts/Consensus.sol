pragma solidity ^0.5.0;

import "./interfaces/IGroup.sol";
import "./interfaces/IConsensus.sol";

contract Consensus {
    bool private _isActive;
    address private _proposer;
    address private _group;
    bytes32 _digest;
    bytes _payload;
    address[] _membersThatApproved;

    event Debug(uint256 state);
    event DebugCons(bytes msg);

    constructor (address proposer, address group) public {
        _proposer = proposer;
        _group = group;
    }

    function propose(bytes32 digest, bytes calldata payload) external {
        require(msg.sender == _group, "msg.sender is not group owner");

        _digest = digest;
        _payload = payload;

        _membersThatApproved.length = 0;
        _membersThatApproved.push(_proposer);
        _isActive = true;
    }

    function invalidate() external {
        require(msg.sender == _group, "msg.sender is not group owner");

        _isActive = false;
    }

    function approve(bytes32 r, bytes32 s, uint8 v) public onlyMembers onlyWhenActive {
        require(memberNotVotedYet(msg.sender), "member already voted");
//        require(verify(msg.sender, _digest, v, r, s), "invalid approval: invalid signature");

        _membersThatApproved.push(msg.sender);

        if (_membersThatApproved.length > IGroup(_group).threshold()) {
            IGroup(_group).onChangeIpfsHashConsensus(_payload);
        }
    }

    function verify(address addr, bytes32 hash, uint8 v, bytes32 r, bytes32 s) private pure returns(bool) {
        return ecrecover(hash, v, r, s) == addr;
    }

    function memberNotVotedYet(address member) private view returns(bool) {
        for (uint256 i = 0; i < _membersThatApproved.length; i++) {
            if (_membersThatApproved[i] == member) {
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
