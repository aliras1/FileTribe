pragma solidity ^0.4.23;

import "./openzeppelin-solidity/contracts/ECRecovery.sol";
import "./HeapSortLib.sol";

contract Dipfshare {
    struct User {
        string name;
        string ipfsPeerId;
        bytes32 boxingKey;
        bytes verifyKey;
        bytes32[] groups;
        bool exists;
    }

    struct Friendship {
        bytes from;        // encrypted with pk_to
        bytes to;        // encrypted with pk_from
        
        bytes fromSigningKey;  // encrypted with pk_to
        bytes toSigningKey;  // encrypted with pk_from
        address fromVerifyAddress;
        address toVerifyAddress;
        
        bytes dirOfToByFrom; // encrypted with common key
        bytes dirOfFromByTo; // encrypted with common key
        
        bool confirmed;
    }

    struct Group {        
        address owner;
        string name;
        address[] members;
        string ipfsPath; // TODO: encrypted with group key
        mapping(address => bool) canInvite;
        bool exists;
    }

    struct Signature {
        bytes32 r;
        bytes32 s;
        uint8 v;
    }

    struct Approval {
        address user;
        Signature sig;
    }

    address public owner;
    mapping(address => User) private users; // user.verify_key = key
    mapping(bytes32 => Friendship) private friendships;
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
        bytes32 boxingKey, 
        bytes verifyKey
    ) 
        public
    {
        require(!users[msg.sender].exists, "Username already exists");

        users[msg.sender].name = name;
        users[msg.sender].ipfsPeerId = ipfsPeerId;
        users[msg.sender].boxingKey = boxingKey;
        users[msg.sender].verifyKey = verifyKey;
        users[msg.sender].exists = true;

        emit UserRegistered(msg.sender);
    }

    function isUserRegistered(address id) public view returns(bool) {
        return users[id].exists;
    }

    function getUser(address id) public view returns(string, string, bytes32, bytes) {
        require(users[id].exists, "User does not exist");

        User memory user = users[id];

        return (user.name, user.ipfsPeerId, user.boxingKey, user.verifyKey);
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

        uint length = members.length;
        require(rs.length == length);
        require(ss.length == length);
        require(vs.length == length);
        require(length > groups[groupId].members.length / 2);

        bytes32 digest = keccak256(groups[groupId].ipfsPath, newIpfsPath);

        HeapSortLib.heapSort(members);

        for (uint i = 0; i < members.length; i++) {
            require(isUserInGroup(groupId, members[i]), "invalid approval: user is not a group member");
            require(verify(members[i], digest, vs[i], rs[i], ss[i]), "invalid approval: invalid signature");
            if (i == 0) {
                continue;
            }

            // in a sorted array we can be sure, that
            // if there is no matching addresses next
            // to each other than there is not any in
            // the whole array
            require(members[i] != members[i - 1]);
        }

        // TODO: re-entrance danger
        groups[groupId].ipfsPath = newIpfsPath;

        emit GroupUpdateIpfsPath(groupId, newIpfsPath);
    }

    function sendMessage(bytes message) public {
        emit MessageSent(message);
    }
}