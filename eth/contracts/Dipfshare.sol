pragma solidity ^0.5.0;

import "./interfaces/factory/IAccountFactory.sol";
import "./interfaces/factory/IGroupFactory.sol";
import "./interfaces/factory/IConsensusFactory.sol";
import "./common/Ownable.sol";

contract Dipfshare is Ownable {
    address private _accountFactory;
    address private _groupFactory;
    address private _consensusFactory;
    mapping(address => address) private _accounts;

    event KeyDirty(bytes32 groupId);
    event GroupKeyChanged(bytes32 groupId, bytes ipfsHash);
    event AccountCreated(address owner, address account);
    event GroupRegistered(bytes32 id);
    event GroupLeft(bytes32 groupId, address user);
    event GroupInvitation(address from, address to, bytes32 groupId);
    event GroupUpdateIpfsHash(bytes32 groupId, bytes ipfsHash);
    event Debug(int msg);
    event DebugBytes(bytes msg);

    constructor () public Ownable(msg.sender) {

    }

    function setConsensusFactory(address factory) public onlyOwner {
        _consensusFactory = factory;
    }

    function setAccountFactory(address factory) public onlyOwner {
        _accountFactory = factory;
    }

    function setGroupFactory(address factory) public onlyOwner {
        _groupFactory = factory;
    }

    function createAccount(string memory name, string memory ipfsPeerId, bytes32 boxingKey) public {
        address acc = IAccountFactory(_accountFactory).create(msg.sender, name, ipfsPeerId, boxingKey);

        _accounts[msg.sender] = acc;

        emit AccountCreated(msg.sender, acc);
    }

    function createGroup(string calldata name) external returns(address group) {
        return IGroupFactory(_groupFactory).create(msg.sender, name);
    }

    function createConsensus(
        IConsensus.Type cType,
        address proposer,
        bytes32 digest,
        bytes calldata payload)
    external returns(address consensus) {
        return IConsensusFactory(_consensusFactory).create(
            cType,
            proposer,
            msg.sender,
            digest,
            payload);
    }


    //    function isUserRegistered(address id) public view returns(bool) {
    //        return _users[id].exists;
    //    }

    function getAccount(address addr) public view returns(address) {
        return _accounts[addr];
    }
}

//    function getGroup(bytes32 groupId) public view returns(string name, address[] members, bytes ipfsHash, address leader) {
//        require(groups[groupId].exists, "Group does not exists");
//
//        return (
//            groups[groupId].name,
//            groups[groupId].members,
//            groups[groupId].ipfsHash,
//            getLeader(groupId));
//    }
//
//    function leaveGroup(bytes32 groupId) public {
//        require(isUserInGroup(groupId, msg.sender), "User is not a member of the given group");
//
//        groups[groupId].isKeyDirty = true;
//
//        uint256 len = groups[groupId].members.length;
//        for (uint256 i = 0; i < len; i++) {
//            if (groups[groupId].members[i] == msg.sender) {
//                groups[groupId].members[i] = groups[groupId].members[len - 1];
//                groups[groupId].members.length--;
//                break;
//            }
//        }
//
//        groups[groupId].leaderIdx++;
//        groups[groupId].leaderStart = now;
//
//        emit GroupLeft(groupId, msg.sender);
//        emit KeyDirty(groupId);
//    }
//
//    function getLeader(bytes32 groupId) public view returns(address) {
//        uint256 idx = groups[groupId].leaderIdx % groups[groupId].members.length;
//
//        return groups[groupId].members[idx];
//    }
//
//    function changeLeader(bytes32 groupId) public {
//        require(groups[groupId].exists, "group does not exist");
//        require(isUserInGroup(groupId, msg.sender), "user is not member of group");
//        require(groups[groupId].isKeyDirty, "can not change group leader: key is not dirty");
//
//        // do not punish users who were requesting change leader
//        // on timeout at the same time
//        if (now < groups[groupId].leaderStart + 5 minutes) {
//            return;
//        }
//
//        groups[groupId].leaderIdx++;
//        groups[groupId].leaderStart = now;
//    }
//
//    function checkGroupConsensus(
//        bytes32 groupId,
//        bytes32 digest,
//        address[] memory members,
//        bytes32[] memory rs,
//        bytes32[] memory ss,
//        uint8[] memory vs)
//    internal returns(bool) {
//        require(rs.length == members.length, "invalid r length");
//        require(ss.length == members.length, "invalid s length");
//        require(vs.length == members.length, "invalid v length");
//        require(members.length > groups[groupId].members.length / 2, "not enough approvals");
//
//        for (uint256 i = 0; i < members.length; i++) {
//            require(isUserInGroup(groupId, members[i]), "invalid approval: user is not a group member");
//            require(verify(members[i], digest, vs[i], rs[i], ss[i]), "invalid approval: invalid signature");
//        }
//
//        heapSort(members);
//        for (i = 0; i < members.length; i++) {
//            // in a sorted array we can be sure, that
//            // if there is no matching addresses next
//            // to each other than there is not any in
//            // the whole array
//            if (i == 0) {
//                continue;
//            }
//            require(members[i] != members[i - 1], "duplicate approvals detected");
//        }
//
//        return true;
//    }
//
//    function changeGroupKey(
//        bytes32 groupId,
//        bytes newIpfsHash,
//        address[] members,
//        bytes32[] rs,
//        bytes32[] ss,
//        uint8[] vs)
//    public {
//        require(groups[groupId].exists, "group does not exist");
//
//        bytes32 digest = keccak256(groups[groupId].ipfsHash, newIpfsHash, getLeader(groupId));
//        require(checkGroupConsensus(groupId, digest, members, rs, ss, vs));
//
//        groups[groupId].ipfsHash = newIpfsHash;
//        groups[groupId].isKeyDirty = false;
//
//        emit GroupKeyChanged(groupId, newIpfsHash);
//    }
//
//    function inviteUser(bytes32 groupId, address newMember, bool hasInviteRight) public {
//        require(groups[groupId].canInvite[msg.sender] == true, "User can not invite");
//        require(users[newMember].exists, "Can not invite non existent user");
//
//        groups[groupId].members.push(newMember);
//        users[newMember].groups.push(groupId);
//        groups[groupId].canInvite[newMember] = hasInviteRight;
//
//        emit GroupInvitation(msg.sender, newMember, groupId);
//    }
//
//    function isUserInGroup(bytes32 groupId, address user) internal returns(bool) {
//        for (uint i = 0; i < groups[groupId].members.length; i++) {
//            if (groups[groupId].members[i] == user) {
//                return true;
//            }
//        }
//        return false;
//    }
//
//    function grantInviteRight(bytes32 groupId, address member) public {
//        require(groups[groupId].canInvite[msg.sender] == true, "User can not grant invite right");
//        require(users[member].exists, "Can not grant invite right to non existent user");
//        require(isUserInGroup(groupId, member), "Can not grant invite right to a non member user");
//
//        groups[groupId].canInvite[member] = true;
//    }
//
//    function revokeInviteRight(bytes32 groupId, address member) public {
//        require(groups[groupId].canInvite[msg.sender] == true, "User can not revoke invite right");
//        require(users[member].exists, "Can not revoke invite right from non existent user");
//        require(isUserInGroup(groupId, member), "Can not revoke invite right from a non member user");
//        require(member != groups[groupId].owner, "Can not revoke invite right from the owner");
//
//        groups[groupId].canInvite[member] = false;
//    }
//
//    function verify(address user, bytes32 hash, uint8 v, bytes32 r, bytes32 s) internal constant returns(bool) {
//        return ecrecover(hash, v, r, s) == user;
//    }
//
//    function updateGroupIpfsHash(
//        bytes32 groupId,
//        bytes newIpfsHash,
//        address[] members,
//        bytes32[] rs,
//        bytes32[] ss,
//        uint8[] vs
//    )
//        public
//    {
//        require(groups[groupId].exists, "group does not exist");
//
//        bytes32 digest = keccak256(groups[groupId].ipfsHash, newIpfsHash);
//        require(checkGroupConsensus(groupId, digest, members, rs, ss, vs));
//
//        groups[groupId].ipfsHash = newIpfsHash;
//
//        emit GroupUpdateIpfsHash(groupId, newIpfsHash);
//    }



