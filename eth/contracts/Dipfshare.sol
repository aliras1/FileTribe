pragma solidity ^0.4.23;

import "./openzeppelin-solidity/contracts/ECRecovery.sol";

contract Dipfshare {
    struct User {
        string name;
        bytes32 boxingKey;
        bytes verifyKey;
        string ipfsAddr;
        bool exists;
    }

    struct Friendship {
        bytes from;        // encrypted with pk_to
        bytes to;        // encrypted with pk_from
        bytes signingKey;  // encrypted with pk_to
        bytes dirOfToByFrom; // encrypted with common key
        bytes dirOfFromByTo; // encrypted with common key
        bytes32 digest;     // (32 0's | rand), encrypted with common_key = Gen(pk_from, sk_to) = Gen(pk_to, sk_from)
        address verifyAddress;
        bool confirmed;
    }

    address public owner;
    mapping(address => User) private userBindings; // user.verify_key = key
    mapping(bytes32 => Friendship) private friendships;

    event UserRegistered(address addr);
    event MessageSent(bytes message);
    event NewFriendRequest(bytes32 id, bytes from, bytes to, bytes signingKey, bytes32 digest);
    event FriendshipConfirmed(bytes32 id, bytes dirOfToByFrom, bytes dirOfFromByTo);
    event Debug(address addr);

    constructor () public {
        owner = msg.sender;
    }

    function registerUser(
        string name, 
        bytes32 boxingKey, 
        bytes verifyKey, 
        string ipfsAddr
    ) 
        public
    {
        require(!userBindings[msg.sender].exists, "Username already exists");

        userBindings[msg.sender] = User(name, boxingKey, verifyKey, ipfsAddr, true);
        
        emit UserRegistered(msg.sender);
    }

    function isUserRegistered(address id) public view returns(bool) {
        if (userBindings[id].exists)
            return true;
        return false;
    }

    function getUser(address id) public view returns(string, bytes32, bytes, string) {
        require(userBindings[id].exists);
        User memory user = userBindings[id];
        return (user.name, user.boxingKey, user.verifyKey, user.ipfsAddr);
    }

    function sendMessage(bytes message) public {
        emit MessageSent(message);
    }

    function addFriend(
        bytes32 id,
        bytes from, 
        bytes to, 
        bytes signingKey,
        bytes dirOfToByFrom,
        bytes32 digest,
        address verifyAddress
    )
    public {
        emit Debug(msg.sender);

        friendships[id] = Friendship(from, to, signingKey, dirOfToByFrom, new bytes(0), digest, verifyAddress, false);

        emit NewFriendRequest(id, from, to, signingKey, digest);
    }

    function confirmFriendship(bytes32 id, bytes signature, bytes dirOfFromByTo) public {
        address addr = ECRecovery.recover(friendships[id].digest, signature);
        emit Debug(addr);
        
        require(!friendships[id].confirmed);
        require(addr == friendships[id].verifyAddress);

        friendships[id].confirmed = true;
        friendships[id].dirOfFromByTo = dirOfFromByTo;

        emit FriendshipConfirmed(id, friendships[id].dirOfToByFrom, dirOfFromByTo);
    }
}