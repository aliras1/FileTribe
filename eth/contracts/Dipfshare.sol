pragma solidity ^0.4.23;


contract Dipfshare {
    struct User {
        string name;
        bytes32 boxingKey;
        byte[] verifyKey;
        string ipfsAddr;
        bool exists;
    }

    address public owner;
    mapping(address => User) private userBindings; // user.verify_key = key

    event UserRegistered(address addr);
    event MessageSent(byte[] message);

    constructor () public {
        owner = msg.sender;
    }

    function registerUser(
        string name, 
        bytes32 boxingKey, 
        byte[] verifyKey, 
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

    function getUser(address id) public view returns(string, bytes32, byte[], string) {
        require(userBindings[id].exists);
        User memory user = userBindings[id];
        return (user.name, user.boxingKey, user.verifyKey, user.ipfsAddr);
    }

    function sendMessage(byte[] message) public {
        emit MessageSent(message);
    }
}