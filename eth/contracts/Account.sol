pragma solidity ^0.5.0;

import "./interfaces/IDipfshare.sol";
import "./common/Ownable.sol";

contract Account is Ownable {
    address _parent; // Dipfshare
    string private _name;
    string private _ipfsPeerId;
    bytes32 private _boxingKey;
    address[] private _groups;
    address[] private _invitations;

    event NewInvitation(address account, address group);
    event InvitationAccepted(address account, address group);
    event InvitationDeclined(address account, address group);
    event GroupCreated(address account, address group);
    event GroupLeft(address account, address group);
    event Debug(string msg);
    event DebugBytesAcc(bytes msg);

    constructor (
        address parent,
        address owner,
        string memory name,
        string memory ipfsPeerId,
        bytes32 boxingKey)
    public Ownable(owner) {
        _parent = parent;
        _name = name;
        _ipfsPeerId = ipfsPeerId;
        _boxingKey = boxingKey;
    }

    function createGroup(string memory name) public onlyOwner {
        address group = IDipfshare(_parent).createGroup(name);

        _groups.push(group);

        emit GroupCreated(address(this), group);
    }

    // TODO: onlyGroups modifier
    function invite() external {
        _invitations.push(msg.sender);

        emit NewInvitation(address(this), msg.sender);
    }

    function onInvitationAccepted() external {
        uint256 i;

        for (i = 0; i < _invitations.length; i++) {
            if (_invitations[i] == msg.sender) {
                _groups.push(msg.sender);
                _invitations[i] = _invitations[_invitations.length-1];
                _invitations.length--;

                emit InvitationAccepted(address(this), msg.sender);
            }
        }
    }

    function onInvitationDeclined() external {
        uint256 i;

        for (i = 0; i < _invitations.length; i++) {
            if (_invitations[i] == msg.sender) {
                _groups.push(msg.sender);
                _invitations[i] = _invitations[_invitations.length-1];
                _invitations.length--;

                emit InvitationDeclined(address(this), msg.sender);
            }
        }
    }

    function onGroupLeft(address group) external {
        for(uint256 i = 0; i < _groups.length; i++) {
            if (_groups[i] == group) {
                _groups[i] = _groups[_groups.length - 1];
                _groups.length--;

                emit GroupLeft(address(this), group);

                return;
            }
        }
    }

    function groups() public view returns(address[] memory) {
        return _groups;
    }

    modifier onlyParent() {
        require(msg.sender == _parent, "non parent address tried to call");
        _;
    }

    function name() public view returns(string memory) {
        return _name;
    }

    function ipfsId() public view returns(string memory) {
        return _ipfsPeerId;
    }
}
