pragma solidity ^0.5.0;

contract Invitable {
    address[] private _pendingInvitations;

    constructor () internal {

    }

    modifier onlyInvitations() {
        require(isInvitation(msg.sender));
        _;
    }

    function isInvitation(address inv) internal view returns(bool) {
        for (uint256 i = 0; i < _pendingInvitations.length; i++) {
            if (_pendingInvitations[i] == inv) {
                return true;
            }
        }

        return false;
    }

    function removeInvitation(address inv) internal {
        for (uint256 i = 0; i < _pendingInvitations.length; i++) {
            if (inv == _pendingInvitations[i]) {
                _pendingInvitations[i] = _pendingInvitations[_pendingInvitations.length - 1];
                _pendingInvitations.length--;
            }
        }
    }

    function addInvitation(address inv) internal {
        _pendingInvitations.push(inv);
    }

    function invitations() public view returns(address[] memory) {
        return _pendingInvitations;
    }
}
