pragma solidity ^0.5.0;

import "./Invitation.sol";
import "./interfaces/IDipfshare.sol";
import "./common/Ownable.sol";
import "./common/Invitable.sol";

contract Account is Ownable, Invitable {
    address _parent; // Dipfshare
    string private _name;
    string private _ipfsPeerId;
    bytes32 private _boxingKey;
    address[] private _groups;

    event NewInvitation(address account, address inv);
    event InvitationAccepted(address account, address group);
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
    public Ownable(owner) Invitable() {
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
    function invite() external returns(address invitation) {
        Invitation inv = new Invitation(address(this), msg.sender);
        addInvitation(address(inv));

        emit NewInvitation(address(this), address(inv));

        return address(inv);
    }

    function onInvitationAccepted(address group) external onlyInvitations {
        removeInvitation(msg.sender);
        _groups.push(group);

        emit InvitationAccepted(address(this), group);
    }

    function onInvitationDeclined() public onlyInvitations {
        removeInvitation(msg.sender);
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
