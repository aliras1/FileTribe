module.exports = {
    // See <http://truffleframework.com/docs/advanced/configuration>
    // to customize your Truffle configuration!
    networks: {
        development: {
            // host: "localhost",
            // port: 8001,
            // port: 8545,
            network_id: "*", // Match any network id
            websockets: true,
            gasLimit: 4700000,
            // from: "0xc4f45f1822b614116ea5b68d4020f3ae1a0179e5",
            // provider:function () {
            //   let Web3 = require("web3");
            //   let web3 = new Web3();
            //   return new web3.providers.WebsocketProvider("ws://localhost:8001");
            // }

            provider:function () {
              let Web3 = require("web3");
              let web3 = new Web3();
              return new web3.providers.HttpProvider("http://localhost:8000");
            }
        }
    },

    solc: {
        optimizer: {
            enabled: true,
            runs: 200
        }
    },

    mocha: {
        enableTimeouts: false
    }
};
