pragma solidity ^0.5.0;

interface IConsensus {
    enum Type {
        IPFS_HASH,
        KEY,
        LEADER
    }

    function approve(address sender, bytes32 r, bytes32 s, uint8 v) external;
}
