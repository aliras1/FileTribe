truffle compile
cat ./build/contracts/Dipfshare.json | jq -r '.abi' > .contract.abi
cat ./build/contracts/Dipfshare.json | jq -r '.bytecode' > .contract.bin
abigen -abi .contract.abi -pkg eth --out eth.go --bin .contract.bin
rm .contract.abi
rm .contract.bin