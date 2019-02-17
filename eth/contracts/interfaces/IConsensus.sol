pragma solidity ^0.5.0;

interface IConsensus {
    enum Type {
        IPFS_HASH,
        KEY,
        LEADER
    }

    function approveExternal(address sender, bytes32 r, bytes32 s, uint8 v) external;
}
