pragma solidity ^0.5.0;

interface IAccountFactory {

    function create(address owner, string calldata name, string calldata ipfsId, bytes32 key) external returns(address);
}
