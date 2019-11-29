pragma solidity ^0.5.0;

/**
 * @dev These functions deal with verification of Merkle trees (hash trees),
 */
library MerkleTree {

    function getRoot(bytes32[] memory leaves, uint256 length) internal pure returns (bytes32) {
        if (length == 1) {
            return leaves[0];
        }

        for (uint256 i = 0; i < leaves.length / 2; i++) {
            if (leaves[2*i] < leaves[2*i + 1]) {
                leaves[i] = keccak256(abi.encodePacked(leaves[2*i], leaves[2*i + 1]));
            } else {
                leaves[i] = keccak256(abi.encodePacked(leaves[2*i + 1], leaves[2*i]));
            }
        }

        return getRoot(leaves, length / 2);
    }
}
