const truffleAssert = require('truffle-assertions');

const Account = artifacts.require('Account');
const AccountFactory = artifacts.require('AccountFactory');
const Group = artifacts.require('Group');
const GroupFactory = artifacts.require('GroupFactory');
const Consensus = artifacts.require('Consensus');
const ConsensusFactory = artifacts.require('ConsensusFactory');
const FileTribeDApp = artifacts.require('FileTribeDApp');

var ipfsHashChange = false;

var fs = require('fs');
var util = require('util');
var log_file = fs.createWriteStream(__dirname + '/debug.log', {flags : 'w'});
var log_stdout = process.stdout;

// console.log = function(d) { //
//     log_file.write(util.format(d) + '\n');
//     log_stdout.write(util.format(d) + '\n');
// };

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
        let group = await Group.at(event.args.group);

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

        let result = await group.join({from: inviteeAddress});
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
    var result = await group.commit(ipfsHash, {from: fromAddress});
    console.log("change ipfs hash: " + result.receipt.gasUsed);
    truffleAssert.eventEmitted(result, 'NewConsensus', (ev) => {
        consensusAddress = ev.consensus;
        return true;
    }, 'No NewConsensus event!');
    assert.notEqual(consensusAddress, undefined);

    let consensus = await Consensus.at(consensusAddress);
    assert.notEqual(consensus, undefined);

    // register IpfsHashChanged event listener
    keyDirty = false;
    await group.IpfsHashChanged().on('data', async event => {
        keyDirty = true;
    });

    approverAddresses.forEach(async addr => {
        let idx = approverAddresses.findIndex(elem => {
            return elem === addr;
        });

        await sleep(100 * idx);

        console.log("approving cons " + consensus.address + " from addr: " + addr);
        try {
            result = await consensus.approve({from: addr});
        } catch (e) {
            console.log(e + ": " + addr);
        }

        console.log("aprroval: " + result.receipt.gasUsed);
    });

    await sleep(10000);
    assert(keyDirty);

    assert.equal(await group._ipfsHash(), ipfsHash);
}

async function changeKeyDecline(fromAddress, ipfsHash, approverAddresses, group) {
    let consensusAddress;
    result = await group.changeKey(ipfsHash, {from: fromAddress});
    console.log("change ipfs hash: " + result.receipt.gasUsed);
    truffleAssert.eventEmitted(result, 'NewConsensus', (ev) => {
        consensusAddress = ev.consensus;
        return true;
    }, 'No NewConsensus event!');
    assert.notEqual(consensusAddress, undefined);

    let consensus = await Consensus.at(consensusAddress);
    assert.notEqual(consensus, undefined);

    // register KeyDirty event listener
    keyDirty = false;
    await group.KeyDirty().on('data', async event => {
        keyDirty = true;
    });

    approverAddresses.forEach(async addr => {
        let idx = approverAddresses.findIndex(elem => {
            return elem === addr;
        });

        await sleep(100 * idx);

        console.log("approving cons " + consensus.address + " from addr: " + addr);
        try {
            result = await consensus.decline("0x00", "0x00", 1, {from: addr});
        } catch (e) {
            console.log(e + ": " + addr);
        }

        console.log("decline: " + result.receipt.gasUsed);
    });

    await sleep(10000);
    assert(keyDirty);

    assert.notEqual(await group.ipfsHash(), ipfsHash);
}

async function changeKey(fromAddress, ipfsHash, approverAddresses, group) {
    let consensusAddress;
    result = await group.changeKey(ipfsHash, {from: fromAddress});
    console.log("change ipfs hash: " + result.receipt.gasUsed);
    truffleAssert.eventEmitted(result, 'NewConsensus', (ev) => {
        consensusAddress = ev.consensus;
        return true;
    }, 'No NewConsensus event!');
    assert.notEqual(consensusAddress, undefined);

    let consensus = await Consensus.at(consensusAddress);
    assert.notEqual(consensus, undefined);

    // register KeyDirty event listener
    isConsensus = false;
    await group.ConsensusReached().on('data', async event => {
        isConsensus = true;
    });

    approverAddresses.forEach(async addr => {
        let idx = approverAddresses.findIndex(elem => {
            return elem === addr;
        });

        await sleep(100 * idx);

        console.log("approving cons " + consensus.address + " from addr: " + addr);
        try {
            result = await consensus.approve("0x00", "0x00", 1, {from: addr});
        } catch (e) {
            console.log(e + ": " + addr);
        }

        console.log("approve key: " + result.receipt.gasUsed);
    });

    await sleep(10000);
    assert(isConsensus);

    assert.equal(await group.ipfsHash(), ipfsHash);
}

contract('FileTribeDApp', accounts => {
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
        dipfshare = await FileTribeDApp.deployed();
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

        // var acc;
        //
        // for (let i = 4; i < 10; i++) {
        //     acc = await createAccount(dipfshare, ethAccounts[i], "Acc"+i, "ipfs"+i, "0x06");
        //     await invite(group, alice, aliceAddress, acc, ethAccounts[i]);
        // }

        // check members
        let members = await group.memberOwners();
        // console.log(members);
        // console.log(bob.address);
        // console.log(await bob.owner());

        // members commit their changes
        // await commit(ethAccounts[8], "0x01", ethAccounts, group);
        await commit(aliceAddress, "0x01", [bobAddress, charlesAddress], group);
        await sleep(1000);

        let consensusAddress;
        result = await group.commit("0x02", {from: bobAddress});
        console.log("change ipfs hash: " + result.receipt.gasUsed);
        truffleAssert.eventEmitted(result, 'NewConsensus', (ev) => {
            consensusAddress = ev.consensus;
            return true;
        }, 'No NewConsensus event!');

        await commit(aliceAddress, "0x03", [bobAddress, charlesAddress], group);
        await sleep(1000);
        await commit(bobAddress, "0x04", [aliceAddress, charlesAddress], group);
        await sleep(1000);

        // await commit(bobAddress, "0x03", [ethAccounts[4], charlesAddress], group);
        // await sleep(1000);
        // await commit(bobAddress, "0x05", [ethAccounts[4], charlesAddress], group);

        // result = await group.kick(acc.address, {from: aliceAddress});
        // console.log("kick: " + result.receipt.gasUsed);

        // charles leaves the group
        result = await group.leave({from: bobAddress});
        console.log("leave: " + result.receipt.gasUsed);

        truffleAssert.eventEmitted(result, 'MemberLeft', (ev) => {
            return group.address === ev.group;
        }, 'Contract should emit a valid MemberLeft event');

        await sleep(1000);
        // let numMembers = await group.members();
        // console.log(numMembers.length());
        // assert.equal(numMembers, 2, "expected number of members is 2, got " + numMembers)
    });
});