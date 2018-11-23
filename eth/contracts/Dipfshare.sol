pragma solidity ^0.4.23;

//import "../HeapSortLib.sol";

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
        bytes ipfsHash; // encrypted with group key
        mapping(address => bool) canInvite;
        bool isKeyDirty;
        bool exists;

        uint256 leaderIdx;
        uint leaderStart;
    }

    address public owner;
    mapping(address => User) private users;
    mapping(bytes32 => Group) private groups;

    event KeyDirty(bytes32 groupId);
    event GroupKeyChanged(bytes32 groupId, bytes ipfsHash);
    event UserRegistered(address addr);
    event GroupRegistered(bytes32 id);
    event GroupLeft(bytes32 groupId, address user);
    event GroupInvitation(address from, address to, bytes32 groupId);
    event GroupUpdateIpfsHash(bytes32 groupId, bytes ipfsHash);
    event Debug(int msg);

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
        bytes ipfsHash
    ) 
        public
    {
        require(!groups[id].exists, "A group with the given id already exists");

        groups[id].owner = msg.sender;
        groups[id].name = name;
        groups[id].members.push(msg.sender);
        groups[id].ipfsHash = ipfsHash;
        groups[id].exists = true;
        groups[id].canInvite[msg.sender] = true;
        groups[id].leaderIdx = 0;

        emit GroupRegistered(id);
    }

    function getGroup(bytes32 groupId) public view returns(string name, address[] members, bytes ipfsHash, address leader) {
        require(groups[groupId].exists, "Group does not exists");

        return (
            groups[groupId].name,
            groups[groupId].members,
            groups[groupId].ipfsHash,
            getLeader(groupId));
    }

    function leaveGroup(bytes32 groupId) public {
        require(isUserInGroup(groupId, msg.sender), "User is not a member of the given group");

        groups[groupId].isKeyDirty = true;

        uint256 len = groups[groupId].members.length;
        for (uint256 i = 0; i < len; i++) {
            if (groups[groupId].members[i] == msg.sender) {
                groups[groupId].members[i] = groups[groupId].members[len - 1];
                groups[groupId].members.length--;
                break;
            }
        }

        groups[groupId].leaderIdx++;
        groups[groupId].leaderStart = now;

        emit GroupLeft(groupId, msg.sender);
        emit KeyDirty(groupId);
    }

    function getLeader(bytes32 groupId) public view returns(address) {
        uint256 idx = groups[groupId].leaderIdx % groups[groupId].members.length;

        return groups[groupId].members[idx];
    }

    function changeLeader(bytes32 groupId) public {
        require(groups[groupId].exists, "group does not exist");
        require(isUserInGroup(groupId, msg.sender), "user is not member of group");
        require(groups[groupId].isKeyDirty, "can not change group leader: key is not dirty");

        // do not punish users who were requesting change leader
        // on timeout at the same time
        if (now < groups[groupId].leaderStart + 5 minutes) {
            return;
        }

        groups[groupId].leaderIdx++;
        groups[groupId].leaderStart = now;
    }

    function checkGroupConsensus(
        bytes32 groupId,
        bytes32 digest,
        address[] memory members,
        bytes32[] memory rs,
        bytes32[] memory ss,
        uint8[] memory vs)
    internal returns(bool) {
        require(rs.length == members.length, "invalid r length");
        require(ss.length == members.length, "invalid s length");
        require(vs.length == members.length, "invalid v length");
        require(members.length > groups[groupId].members.length / 2, "not enough approvals");

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

        return true;
    }

    function changeGroupKey(
        bytes32 groupId,
        bytes newIpfsHash,
        address[] members,
        bytes32[] rs,
        bytes32[] ss,
        uint8[] vs)
    public {
        require(groups[groupId].exists, "group does not exist");

        bytes32 digest = keccak256(groups[groupId].ipfsHash, newIpfsHash, getLeader(groupId));
        require(checkGroupConsensus(groupId, digest, members, rs, ss, vs));

        groups[groupId].ipfsHash = newIpfsHash;
        groups[groupId].isKeyDirty = false;

        emit GroupKeyChanged(groupId, newIpfsHash);
    }

    function inviteUser(bytes32 groupId, address newMember, bool hasInviteRight) public {
        require(groups[groupId].canInvite[msg.sender] == true, "User can not invite");
        require(users[newMember].exists, "Can not invite non existent user");

        groups[groupId].members.push(newMember);
        users[newMember].groups.push(groupId);
        groups[groupId].canInvite[newMember] = hasInviteRight;

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

    function grantInviteRight(bytes32 groupId, address member) public {
        require(groups[groupId].canInvite[msg.sender] == true, "User can not grant invite right");
        require(users[member].exists, "Can not grant invite right to non existent user");
        require(isUserInGroup(groupId, member), "Can not grant invite right to a non member user");

        groups[groupId].canInvite[member] = true;
    }

    function revokeInviteRight(bytes32 groupId, address member) public {
        require(groups[groupId].canInvite[msg.sender] == true, "User can not revoke invite right");
        require(users[member].exists, "Can not revoke invite right from non existent user");
        require(isUserInGroup(groupId, member), "Can not revoke invite right from a non member user");
        require(member != groups[groupId].owner, "Can not revoke invite right from the owner");

        groups[groupId].canInvite[member] = false;
    }

    function verify(address user, bytes32 hash, uint8 v, bytes32 r, bytes32 s) internal constant returns(bool) {
        return ecrecover(hash, v, r, s) == user;
    }

    function updateGroupIpfsHash(
        bytes32 groupId,
        bytes newIpfsHash,
        address[] members,
        bytes32[] rs,
        bytes32[] ss,
        uint8[] vs
    )
        public
    {
        require(groups[groupId].exists, "group does not exist");

        bytes32 digest = keccak256(groups[groupId].ipfsHash, newIpfsHash);
        require(checkGroupConsensus(groupId, digest, members, rs, ss, vs));

        groups[groupId].ipfsHash = newIpfsHash;

        emit GroupUpdateIpfsHash(groupId, newIpfsHash);
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