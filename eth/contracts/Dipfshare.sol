pragma solidity ^0.4.23;

//import "./openzeppelin-solidity/contracts/ECRecovery.sol";
//import "./HeapSortLib.sol";

contract Dipfshare {
    struct User {
        string name;
        string ipfsPeerId;
        bytes32 boxingKey;
        bytes32[] groups;
        bool exists;
    }

    struct Group {
        address owner;
        string name;
        address[] members;
        string ipfsPath; // TODO: encrypted with group key
        mapping(address => bool) canInvite;
        bool exists;
    }

    address public owner;
    mapping(address => User) private users; // user.verify_key = key
    mapping(bytes32 => Group) private groups;

    event UserRegistered(address addr);
    event GroupRegistered(bytes32 id);
    event GroupInvitation(address from, address to, bytes32 groupId);
    event GroupUpdateIpfsPath(bytes32 groupId, string ipfsPath);
    event MessageSent(bytes message);
    event Debug(bytes msg);

    constructor () public {
        owner = msg.sender;
    }

    function registerUser(
        string name,
        string ipfsPeerId,
        bytes32 boxingKey
    ) 
        public
    {
        require(!users[msg.sender].exists, "Username already exists");

        users[msg.sender].name = name;
        users[msg.sender].ipfsPeerId = ipfsPeerId;
        users[msg.sender].boxingKey = boxingKey;
        users[msg.sender].exists = true;

        emit UserRegistered(msg.sender);
    }

    function isUserRegistered(address id) public view returns(bool) {
        return users[id].exists;
    }

    function getUser(address id) public view returns(string, string, bytes32) {
        require(users[id].exists, "User does not exist");

        User memory user = users[id];

        return (user.name, user.ipfsPeerId, user.boxingKey);
    }

    function createGroup(
        bytes32 id,
        string name,
        string ipfsPath
    ) 
        public
    {
        require(!groups[id].exists, "A group with the given id already exists");

        groups[id].owner = msg.sender;
        groups[id].name = name;
        groups[id].members.push(msg.sender);
        groups[id].ipfsPath = ipfsPath;
        groups[id].exists = true;
        groups[id].canInvite[msg.sender] = true;
        
        emit GroupRegistered(id);
    }

    function getGroup(bytes32 groupId) public view returns(string, address[], string) {
        require(groups[groupId].exists, "Group does not exists");

        return (groups[groupId].name, groups[groupId].members, groups[groupId].ipfsPath);
    }

    function inviteUser(bytes32 groupId, address newMember) public {
        require(groups[groupId].canInvite[msg.sender] == true, "User can not invite");
        require(users[newMember].exists, "Can not invite non existent user");

        groups[groupId].members.push(newMember);
        users[newMember].groups.push(groupId);

        emit GroupInvitation(msg.sender, newMember, groupId);
    }

    function isUserInGroup(bytes32 groupId, address user) internal returns(bool) {
        for (uint i = 0; i < groups[groupId].members.length; i++) {
            if (groups[groupId].members[i] == user) {
                return true;
            }
        }
        return false;
    }

    function verify(address user, bytes32 hash, uint8 v, bytes32 r, bytes32 s) internal constant returns(bool) {
        return ecrecover(hash, v, r, s) == user;
    }

    function updateGroupIpfsPath(
        bytes32 groupId,
        string newIpfsPath,
        address[] members,
        bytes32[] rs,
        bytes32[] ss,
        uint8[] vs
    )
        public
    {
        require(groups[groupId].exists, "group does not exist");

        require(rs.length == members.length);
        require(ss.length == members.length);
        require(vs.length == members.length);
        require(members.length > groups[groupId].members.length / 2);

        bytes32 digest = keccak256(groups[groupId].ipfsPath, newIpfsPath);

        for (uint256 i = 0; i < members.length; i++) {
            require(isUserInGroup(groupId, members[i]), "invalid approval: user is not a group member");
            require(verify(members[i], digest, vs[i], rs[i], ss[i]), "invalid approval: invalid signature");
        }

        heapSort(members);
        for (i = 0; i < members.length; i++) {
            // in a sorted array we can be sure, that
            // if there is no matching addresses next
            // to each other than there is not any in
            // the whole array
            if (i == 0) {
                continue;
            }
            require(members[i] != members[i - 1], "duplicate approvals detected");
        }

        // TODO: re-entrance danger
        groups[groupId].ipfsPath = newIpfsPath;

        emit GroupUpdateIpfsPath(groupId, newIpfsPath);
    }

    function sendMessage(bytes message) public {
        emit MessageSent(message);
    }

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