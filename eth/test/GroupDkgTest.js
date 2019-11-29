const truffleAssert = require('truffle-assertions');
const helper = require("./helpers/truffleTestHelper");
let testData = require('./test_data.json');

const GroupDkg = artifacts.require('GroupDkg');

var fs = require('fs');
var util = require('util');
var log_file = fs.createWriteStream(__dirname + '/debug.log', {flags : 'w'});
var log_stdout = process.stdout;

console.log = function(d) { //
    log_file.write(util.format(d) + '\n');
    log_stdout.write(util.format(d) + '\n');
};

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

contract('GroupDkg', accounts => {
    let groupDkg;    

    beforeEach(async function () {
        groupDkg = await GroupDkg.deployed();        
    });

    it('integration test', async function () {
        for (let i = 0; i < testData.length; i++) {
            const testCase = testData[i];
            let t = testCase.t;
            let n = testCase.n;
            
            console.log("N=" + n)
            await groupDkg.init(t, n);

            // Join
            console.log("JOINING")
            for (let i = 0; i < n; i++) {
                let result = await groupDkg.join(testCase.pks[i], {from:accounts[i]})
                console.log(result);
            }
            
            // Commit
            console.log("COMMITTING")
            for (let i = 0; i < n; i++) {
                let commit = testCase.commits[i];
                let result = await groupDkg.commit(
                    commit.senderIndex,
                    commit.rootPubCommit,
                    commit.rootEncPrvCommit,
                    commit.yG1,
                    commit.commitIpfsHash,
                    {from:accounts[i]})
                console.log(result);
                console.log(result.receipt.gasUsed);
            }
            
            console.log("WAIT");
            for (let j = 0; j < 40; j++) {
                await helper.advanceTimeAndBlock(1);                
            }
            console.log("FINALIZE");
            let result = await groupDkg.postCommitTimedOut(1, [1, 2], {from:accounts[0]});
            console.log(result);
            console.log(result.receipt.gasUsed);

            // console.log("CALC GROUP VK")
            // let result = await groupDkg.calculateGroupPK()
            // console.log(result)
            // console.log(result.receipt.gasUsed);

            // result = await groupDkg.getGroupPK()
            // console.log(result)
        }
    });
});