pragma solidity ^0.4.23;


contract Dipfshare {
    struct User {
        address addr;
        string name;
        bytes32 boxingKey;
        bytes32 verifyKey;
        string ipfsAddr;
        bool exists;
    }

    address public owner;
    mapping(bytes32 => User) private userBindings; // user.verify_key = key

    event UserRegistered(address addr);
    event MessageSent(byte[] message);

    constructor () public {
        owner = msg.sender;
    }

    function registerUser(
        bytes32 id,
        string name, 
        bytes32 boxingKey, 
        bytes32 verifyKey, 
        string ipfsAddr
    ) 
        public
    {
        require(!userBindings[verifyKey].exists, "Username already exists");

        userBindings[id] = User(msg.sender, name, boxingKey, verifyKey, ipfsAddr, true);
        
        emit UserRegistered(msg.sender);
    }

    function isUserRegistered(bytes32 id) public view returns(bool) {
        if (userBindings[id].exists)
            return true;
        return false;
    }

    function getUser(bytes32 id) public view returns(address, string, bytes32, bytes32, string) {
        require(userBindings[id].exists);
        User memory user = userBindings[id];
        return (user.addr, user.name, user.boxingKey, user.verifyKey, user.ipfsAddr);
    }

    function sendMessage(byte[] message) public {
        emit MessageSent(message);
    }
}