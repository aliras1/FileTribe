pragma solidity ^0.5.0;

import "./IAccount.sol";

interface IGroup {
    function isMember(address owner) external view returns(bool);

    function threshold() external view returns(uint256);
}
