pragma solidity ^0.5.0;

import "./interfaces/IAccount.sol";
import "./interfaces/IGroup.sol";

contract Invitation {
    address _group;
    address _account;

    constructor (address account, address group) public {
        _group = group;
        _account = account;
    }

    function accept() public onlyInvitee {
        IAccount(_account).onInvitationAccepted(_group);
        IGroup(_group).onInvitationAccepted(_account);
    }

    function decline() public onlyInvitee {
        IAccount(_account).onInvitationDeclined();
        IGroup(_group).onInvitationDeclined(_account);
    }

    modifier onlyInvitee() {
        require(msg.sender == IAccount(_account).owner());
        _;
    }

    function group() public view returns(address) {
        return _group;
    }
}
