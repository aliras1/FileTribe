pragma solidity ^0.4.23;

library HeapSortLib {
    /// @dev Utility function for heapSort
    /// @param index The index of child node
    /// @return pI The parent node index
    function getParentI(uint256 index) private pure returns (uint256 pI) {
        uint256 i = index - 1;
        pI = i/2;
    }

    /// @dev Utility function for heapSort
    /// @param index The index of parent node
    /// @return lcI The index of left child
    function getLeftChildI(uint256 index) private pure returns (uint256 lcI) {
        uint256 i = index * 2;
        lcI = i + 1;
    }

    /// @dev Sorts given array in place
    /// @param self Storage array containing uint256 type variables
    function heapSort(address[] memory self) public {
        uint256 end = self.length - 1;
        uint256 start = getParentI(end);
        uint256 root = start;
        uint256 lChild;
        uint256 rChild;
        uint256 swap;
        address temp;
        while(start >= 0){
            root = start;
            lChild = getLeftChildI(start);
            while(lChild <= end){
                rChild = lChild + 1;
                swap = root;
                if(self[swap] < self[lChild])
                    swap = lChild;
                if((rChild <= end) && (self[swap]<self[rChild]))
                    swap = rChild;
                if(swap == root)
                    lChild = end+1;
                else {
                    temp = self[swap];
                    self[swap] = self[root];
                    self[root] = temp;
                    root = swap;
                    lChild = getLeftChildI(root);
                }
            }
            if(start == 0)
                break;
            else
                start = start - 1;
        }
        while(end > 0){
            temp = self[end];
            self[end] = self[0];
            self[0] = temp;
            end = end - 1;
            root = 0;
            lChild = getLeftChildI(0);
            while(lChild <= end){
                rChild = lChild + 1;
                swap = root;
                if(self[swap] < self[lChild])
                    swap = lChild;
                if((rChild <= end) && (self[swap]<self[rChild]))
                    swap = rChild;
                if(swap == root)
                    lChild = end + 1;
                else {
                    temp = self[swap];
                    self[swap] = self[root];
                    self[root] = temp;
                    root = swap;
                    lChild = getLeftChildI(root);
                }
            }
        }
    }

}