var FileTribeDApp = artifacts.require("./FileTribeDApp.sol");
var AccountFactory = artifacts.require("./factory/AccountFactory.sol");
var GroupFactory = artifacts.require("./factory/GroupFactory.sol");
var ConsensusFactory = artifacts.require("./factory/ConsensusFactory.sol");
var DkgFactory = artifacts.require("./factory/DkgFactory.sol");
var ecOps = artifacts.require("./ecOps.sol");
var GroupDkg = artifacts.require("./GroupDkg.sol");

module.exports = function(deployer) {
    var app;
    var accFactory;
    var groupFactory;
    var consFactory;
    var dkgFactory;

    deployer.deploy(ecOps);
    deployer.link(ecOps, GroupDkg);
    deployer.link(ecOps, FileTribeDApp);
    deployer.link(ecOps, DkgFactory);
    deployer.link(ecOps, GroupFactory);

    deployer.deploy(FileTribeDApp).then(() => {
        return deployer.deploy(GroupFactory, {gas: 4700000});
    }).then(() => {
        return deployer.deploy(ConsensusFactory);
    }).then(() => {
        return deployer.deploy(AccountFactory);
    }).then(() => {
        return deployer.deploy(DkgFactory);
    }).then(() => {
        return FileTribeDApp.deployed();
    }).then((instance) => {
        app = instance;
        return AccountFactory.deployed();
    }).then((instance) => {
        accFactory = instance;
        return GroupFactory.deployed();
    }).then((instance) => {
        groupFactory = instance;
        return ConsensusFactory.deployed();
    }).then((instance) => {
        consFactory = instance;
        return DkgFactory.deployed();
    }).then((instance) => {
        dkgFactory = instance;
        return app.setAccountFactory(accFactory.address);
    }).then(() => {
        return accFactory.setParent(app.address);
    }).then(() => {
        return app.setGroupFactory(groupFactory.address);
    }).then(() => {
        return groupFactory.setParent(app.address);
    }).then(() => {
        return app.setConsensusFactory(consFactory.address);
    }).then(() => {
        return consFactory.setParent(app.address);
    }).then(() => {
        return app.setDkgFactory(dkgFactory.address);
    }).then(() => {
        return dkgFactory.setParent(app.address);
    });
};
