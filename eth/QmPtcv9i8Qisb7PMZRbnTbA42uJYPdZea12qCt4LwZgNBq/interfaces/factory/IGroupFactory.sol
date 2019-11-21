pragma solidity ^0.5.0;

import "../IAccount.sol";
import "../IGroup.sol";

interface IGroupFactory {

    function create(IAccount owner, string calldata name) external returns(IGroup);
}
