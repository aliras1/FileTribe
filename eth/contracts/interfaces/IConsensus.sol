pragma solidity ^0.5.0;

interface IConsensus {
    enum Type {
        IPFS_HASH,
        KEY,
        LEADER
    }

    function approveExternal(address sender, bytes32 r, bytes32 s, uint8 v) external;

    function propose(bytes32 digest, bytes calldata payload) external;

    function invalidate() external;

    function setProposer(address proposer) external;
}
