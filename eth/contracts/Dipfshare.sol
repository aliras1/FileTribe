pragma solidity ^0.4.23;


contract Dipfshare {
    struct User {
        address addr;
        bytes32 boxingKey;
        bytes32 verifyKey;
        string ipfsAddr;
        bool exists;
    }

    address public owner;
    mapping(string => User) private userBindings;

    event RegisteredUser(address addr);

    constructor () public {
        owner = msg.sender;
    }

    function registerUser(
        string name, 
        bytes32 boxingKey, 
        bytes32 verifyKey, 
        string ipfsAddr
    ) 
        public
    {
        require(!userBindings[name].exists, "Username already exists");
        userBindings[name] = User(msg.sender, boxingKey, verifyKey, ipfsAddr, true);
        emit RegisteredUser(msg.sender);
    }

    function isUserRegistered(string name) public view returns(bool) {
        if (userBindings[name].exists)
            return true;
        return false;
    }

    function getUser(string name) public view returns(address, bytes32, bytes32, string) {
        require(userBindings[name].exists);
        User memory user = userBindings[name];
        return (user.addr, user.boxingKey, user.verifyKey, user.ipfsAddr);
    }
}