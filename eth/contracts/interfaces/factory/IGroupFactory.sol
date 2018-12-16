pragma solidity ^0.5.0;

interface IGroupFactory {

    function create(address owner, string calldata name) external returns(address);
}
