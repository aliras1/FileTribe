pragma solidity ^0.5.0;

/**
 * @dev These functions deal with verification of Merkle trees (hash trees),
 */
library MerkleTree {

    // function root(bytes32[] memory leaves) internal pure returns (bytes32) {
    //     uint256 curL = leaves.length;
    //     uint256 prevL = leaves.length;
    //     uint256 i = 0;
    //     uint256 j = 0;
    //     uint256 k = 1;
    //     bool rest = false;
    //     while (curL > 1 || (curL == 1 && rest)) {
    //         if (leaves[j] < leaves[k]) {
    //             leaves[i] = keccak256(abi.encodePacked(leaves[j], leaves[k]));
    //         } else {
    //             leaves[i] = keccak256(abi.encodePacked(leaves[k], leaves[j]));
    //         }

    //         i++;
    //         j += 2;
    //         k += 2;

    //         if (k >= curL) {
    //             i = 0;
    //             k = k - curL;
    //             prevL = curL;
    //             curL = curL / 2;
    //             if (2*curL != prevL) {
    //                 rest = true;
    //             }
    //         }
    //         if (j >= prevL) {
    //             j = j - prevL;
    //         }
    //     }
    //     return leaves[0];
    // }

    function root(bytes32[] memory leaves) internal pure returns (bytes32) {
        for (uint256 curLength = leaves.length; curLength > 1; curLength /= 2) {
            for (uint256 i = 0; i < curLength / 2; i++) {
                if (leaves[2*i] < leaves[2*i + 1]) {
                    leaves[i] = keccak256(abi.encodePacked(leaves[2*i], leaves[2*i + 1]));
                } else {
                    leaves[i] = keccak256(abi.encodePacked(leaves[2*i + 1], leaves[2*i]));
                }
            }
            if (curLength % 2 == 1) {
                leaves[curLength / 2] = leaves[curLength - 1];
                curLength++;
            }
        }
        return leaves[0];
    }

    // function getRoot(bytes32[] memory leaves, uint256 length) internal pure returns (bytes32) {
    //     if (length == 1) {
    //         return leaves[0];
    //     }

    //     for (uint256 i = 0; i < leaves.length / 2; i++) {
    //         if (leaves[2*i] < leaves[2*i + 1]) {
    //             leaves[i] = keccak256(abi.encodePacked(leaves[2*i], leaves[2*i + 1]));
    //         } else {
    //             leaves[i] = keccak256(abi.encodePacked(leaves[2*i + 1], leaves[2*i]));
    //         }
    //     }

    //     return getRoot(leaves, length / 2);
    // }
}
