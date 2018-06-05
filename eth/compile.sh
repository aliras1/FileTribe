truffle compile
cat ./build/contracts/Storage.json | jq -r '.abi' > .contract.abi
abigen -abi .contract.abi -pkg eth --out eth.go
rm .contract.abi