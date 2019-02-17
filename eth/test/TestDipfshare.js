const truffleAssert = require('truffle-assertions');

const Account = artifacts.require('Account');
const AccountFactory = artifacts.require('AccountFactory');
const Group = artifacts.require('Group');
const GroupFactory = artifacts.require('GroupFactory');
const Consensus = artifacts.require('Consensus');
const ConsensusFactory = artifacts.require('ConsensusFactory');
const Dipfshare = artifacts.require('Dipfshare');
const Invitation = artifacts.require('Invitation');

var ipfsHashChange = false;

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

async function createAccount(dipfshare, address, name, ipfsId, publicKey) {
    let result = await dipfshare.createAccount(
        name,
        ipfsId,
        publicKey,
        {from: address}
    );
    console.log("register: " + result.receipt.gasUsed);

    truffleAssert.eventEmitted(result, 'AccountCreated', (ev) => {
        return ev.owner === address;
    }, 'Contract should return the correct message.');

    const accountAddress = await dipfshare.getAccount(address);
    let account = await Account.at(accountAddress);

    const accountOwner = await account.owner();
    const accountName = await account.name();
    const accountIpfsId = await account.ipfsId();

    assert(accountOwner === address);
    assert(accountName === name);
    assert(accountIpfsId === ipfsId);

    return account;
}

async function invite(group, inviter, inviterAddress, invitee, inviteeAddress) {
    await invitee.NewInvitation().on('data', async event =>  {
        let inv = await Invitation.at(event.args.inv);

        var inviteeAccepted = false;
        var groupAccepted = false;
        invitee.InvitationAccepted().on('data', async event => {
            if (event.args.group === group.address) {
                inviteeAccepted = true;
            }
        });
        // group.InvitationAccepted().on('data', async event => {
        //     console.log(event.args.account);
        //     console.log(invitee.address);
        //     if (event.args.account === invitee.address) {
        //         groupAccepted = true;
        //     }
        // });

        let result = await inv.accept({from: inviteeAddress});
        console.log("inv accept: " + result.receipt.gasUsed);
        await sleep(100);

        // assert(inviteeAccepted && groupAccepted);
        assert(inviteeAccepted);
    });

    result = await group.invite(invitee.address, {from: inviterAddress});
    console.log("invite:" + result.receipt.gasUsed);
    truffleAssert.eventEmitted(result, 'InvitationSent', (ev) => {
        return ev.group === group.address && ev.account === invitee.address;
    }, 'Contract should return the correct message.');

    await sleep(200);
}

async function commit(fromAddress, ipfsHash, approverAddresses, group) {
    let consensusAddress;
    result = await group.changeIpfsHash(ipfsHash, {from: fromAddress});
    console.log("change ipfs hash: " + result.receipt.gasUsed);
    truffleAssert.eventEmitted(result, 'NewConsensus', (ev) => {
        consensusAddress = ev.consensus;
        return true;
    }, 'No NewConsensus event!');
    assert.notEqual(consensusAddress, undefined);

    let consensus = await Consensus.at(consensusAddress);
    assert.notEqual(consensus, undefined);

    // register IpfsHashChanged event listener
    ipfsHashChanged = false;
    await group.IpfsHashChanged().on('data', async event => {
        ipfsHashChanged = true;
    });

    approverAddresses.forEach(async addr => {
        if (approverAddresses.findIndex(elem => {
            return elem === addr;
        }) === 1) {
            console.log(1);
            await sleep(100);
        }
        console.log("approving cons " + consensus.address + " from addr: " + addr);
        result = await consensus.approve("0x00", "0x00", 1, {from: addr});
        console.log("aprroval: " + result.receipt.gasUsed);
    });

    await sleep(500);
    assert(ipfsHashChanged);

    assert.equal(await group.ipfsHash(), ipfsHash);
}

contract('Dipfshare', accounts => {
    let dipfshare;
    let accountFactory;
    let groupFactory;
    let consensusFactory;

    const ethAccounts = accounts;
    const creator = accounts[0];
    const aliceAddress = accounts[0];
    const bobAddress = accounts[1];
    const charlesAddress = accounts[2];

    beforeEach(async function () {
        dipfshare = await Dipfshare.new({ from: creator });
        let receipt = await web3.eth.getTransactionReceipt(dipfshare.transactionHash);
        console.log("app: " + receipt.gasUsed);
        accountFactory = await AccountFactory.new({ from: creator });
        receipt = await web3.eth.getTransactionReceipt(accountFactory.transactionHash);
        console.log("acc factory: " + receipt.gasUsed);
        groupFactory = await GroupFactory.new({ from: creator });
        receipt = await web3.eth.getTransactionReceipt(groupFactory.transactionHash);
        console.log("group factory: " + receipt.gasUsed);
        consensusFactory = await ConsensusFactory.new({ from: creator });
        receipt = await web3.eth.getTransactionReceipt(consensusFactory.transactionHash);
        console.log("cons factory: " + receipt.gasUsed);

        await dipfshare.setAccountFactory(accountFactory.address);
        await dipfshare.setGroupFactory(groupFactory.address);
        await dipfshare.setConsensusFactory(consensusFactory.address);

        await accountFactory.setParent(dipfshare.address);
        await groupFactory.setParent(dipfshare.address);
        await consensusFactory.setParent(dipfshare.address);
    });

    // it('check owner', async function () {
    //     const owner = await dipfshare.owner();
    //     assert(owner === creator);
    // });
    //
    // it('register user when it does not exists expect event', async function () {
    //     let username = "Alice";
    //     let publicKey = "0x0000000000000000000000000000000000000000000000000000000000000002";
    //     let ipfsAddress = "rweiro34n3453uz4grsrd";
    //
    //     let result = await dipfshare.createAccount(
    //         username,
    //         ipfsAddress,
    //         publicKey,
    //         {from: aliceAddress}
    //     );
    //
    //     truffleAssert.eventEmitted(result, 'AccountCreated', (ev) => {
    //         return ev.owner === aliceAddress;
    //     }, 'Contract should return the correct message.');
    //
    //     const aliceAccountAddress = await dipfshare.getAccount(aliceAddress);
    //     const aliceAccount = await Account.at(aliceAccountAddress);
    //
    //     const aliceOwner = await aliceAccount.owner();
    //     const aliceName = await aliceAccount.name();
    //     const aliceIpfsId = await aliceAccount.ipfsId();
    //
    //     assert(aliceOwner === aliceAddress);
    //     assert(aliceName === username);
    //     assert(aliceIpfsId === ipfsAddress);
    // });
    //
    // it('register user when it already exists expect error', async function () {
    //     let username = "Alice";
    //     let publicKey = "0x0000000000000000000000000000000000000000000000000000000000000002";
    //     let ipfsAddress = "rweiro34n3453uz4grsrd";
    //
    //     await dipfshare.createAccount(
    //         username,
    //         ipfsAddress,
    //         publicKey,
    //         {from: aliceAddress}
    //     );
    //
    //     let error;
    //     try {
    //         await dipfshare.createAccount(
    //             "asd",
    //             ipfsAddress,
    //             publicKey,
    //             {from: aliceAddress}
    //         );
    //     } catch (exception) {
    //         error = exception;
    //     }
    //
    //     assert.notEqual(error, undefined, 'Error must be thrown');
    //     assert.isAbove(error.message.search('VM Exception while processing transaction: revert'), -1, 'Error: VM Exception while processing transaction: revert');
    //     assert(error.message.search("Account already exists") !== -1);
    // });
    //
    // it('create group from account expect event', async  function () {
    //     await dipfshare.createGroup("group", {from: charlesAddress});
    // });

    it('integration test', async function () {
        // create accounts
        let alice = await createAccount(dipfshare, aliceAddress, "Alice", "ipfs1", "0x01");
        let bob = await createAccount(dipfshare, bobAddress, "Bob", "ipfs2", "0x02");
        let charles = await createAccount(dipfshare, charlesAddress, "Charles", "ipfs3", "0x03");

        // create group
        let groupAddress;
        let result = await alice.createGroup("group", {from: aliceAddress});
        console.log("group created: " + result.receipt.gasUsed);
        truffleAssert.eventEmitted(result, 'GroupCreated', (ev) => {
            groupAddress = ev.group;
            return ev.account === alice.address;
        }, 'Contract should return the correct message.');

        assert.notEqual(groupAddress, undefined, "GroupAddress does not exist");
        let group = await Group.at(groupAddress);
        assert.notEqual(group, undefined, "Group does not exist");

        // invite members
        await invite(group, alice, aliceAddress, bob, bobAddress);
        await invite(group, alice, aliceAddress, charles, charlesAddress);

        var acc;

        for (let i = 4; i < 5; i++) {
            acc = await createAccount(dipfshare, ethAccounts[i], "Acc"+i, "ipfs"+i, "0x06");
            await invite(group, alice, aliceAddress, acc, ethAccounts[i]);
        }

        // check members
        let members = await group.members();
        console.log(members);
        console.log(bob.address);
        console.log(await bob.owner());

        // members commit their changes
        await commit(aliceAddress, "0x01", [bobAddress, charlesAddress], group);
        await sleep(1000);
        await commit(bobAddress, "0x03", [ethAccounts[4], charlesAddress], group);
        await sleep(1000);
        await commit(bobAddress, "0x05", [ethAccounts[4], charlesAddress], group);

        result = await group.kick(acc.address, {from: aliceAddress});
        console.log("kick: " + result.receipt.gasUsed);

        // charles leaves the group
        result = await group.leave({from: charlesAddress});
        console.log("leave: " + result.receipt.gasUsed);
        truffleAssert.eventEmitted(result, 'KeyDirty', (ev) => {
            return group.address === ev.group;
        }, 'Contract should emit a valid KeyDirty event');

        // let numMembers = await group.members();
        // console.log(numMembers.length());
        // assert.equal(numMembers, 2, "expected number of members is 2, got " + numMembers)
    });
});