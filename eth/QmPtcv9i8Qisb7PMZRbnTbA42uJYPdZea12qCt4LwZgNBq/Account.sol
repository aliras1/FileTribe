pragma solidity ^0.5.0;

import "./interfaces/IGroup.sol";
import "./interfaces/IAccount.sol";
import "./interfaces/IFileTribeDApp.sol";
import "./common/Ownable.sol";

contract Account is Ownable, IAccount {
    IFileTribeDApp _fileTribe;
    string private _name;
    string private _ipfsPeerId;
    bytes32 private _boxingKey;
    IGroup[] private _groups;
    IGroup[] private _invitations;

    event NewInvitation(IAccount account, IGroup group);
    event InvitationAccepted(IAccount account, IGroup group);
    event InvitationDeclined(IAccount account, IGroup group);
    event GroupCreated(IAccount account, IGroup group);
    event GroupLeft(IAccount account, IGroup group);
    event Debug(string msg);
    event DebugBytesAcc(bytes msg);

    constructor (
        IFileTribeDApp fileTribe,
        address owner,
        string memory name,
        string memory ipfsPeerId,
        bytes32 boxingKey)
    public Ownable(owner) {
        _fileTribe = fileTribe;
        _name = name;
        _ipfsPeerId = ipfsPeerId;
        _boxingKey = boxingKey;
    }

    function createGroup(string memory name) public onlyOwner {
        IGroup group = _fileTribe.createGroup(name);

        _groups.push(group);

        emit GroupCreated(this, group);
    }

    // TODO: onlyGroups modifier
    function invite() external {
        IGroup group = IGroup(msg.sender);
        _invitations.push(group);

        emit NewInvitation(this, group);
    }

    function onInvitationAccepted() external {
        uint256 i;
        IGroup group = IGroup(msg.sender);

        for (i = 0; i < _invitations.length; i++) {
            if (_invitations[i] == group) {
                _groups.push(group);
                _invitations[i] = _invitations[_invitations.length-1];
                _invitations.length--;

                emit InvitationAccepted(this, group);
            }
        }
    }

    function onInvitationDeclined() external {
        uint256 i;
        IGroup group = IGroup(msg.sender);

        for (i = 0; i < _invitations.length; i++) {
            if (_invitations[i] == group) {
                _groups.push(group);
                _invitations[i] = _invitations[_invitations.length-1];
                _invitations.length--;

                emit InvitationDeclined(this, group);
            }
        }
    }

    function onGroupLeft(IGroup group) external {
        for(uint256 i = 0; i < _groups.length; i++) {
            if (_groups[i] == group) {
                _groups[i] = _groups[_groups.length - 1];
                _groups.length--;

                emit GroupLeft(this, group);

                return;
            }
        }
    }

    function groups() public view returns(IGroup[] memory) {
        return _groups;
    }

    modifier onlyParent() {
        require(IFileTribeDApp(msg.sender) == _fileTribe, "non parent address tried to call");
        _;
    }

    function name() public view returns(string memory) {
        return _name;
    }

    function ipfsId() public view returns(string memory) {
        return _ipfsPeerId;
    }
}
