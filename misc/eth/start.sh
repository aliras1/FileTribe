echo "pwd" > /password

mkdir /devnet

echo "[*] Initializing..."
geth --datadir /devnet --keystore /ethnet/ethkeystore/ init /ethnet/ethgenesis/genesis.json

echo "[*] Creating network node..."
nohup geth --datadir /devnet --keystore /ethnet/ethkeystore/ --networkid 15 --ws --wsaddr "0.0.0.0" --wsport "8001" --wsapi "db,eth,net,web3" --wsorigins "*" --rpc --rpcaddr "0.0.0.0" --rpcport "8000" --rpccorsdomain "*" --port "30304" --rpcapi "db,eth,net,web3" --nat "any" --nodiscover --password /password --unlock "c4f45f1822b614116ea5b68d4020f3ae1a0179e5" --mine --minerthreads=4 &

echo "[*] Network is up"

echo "[*] Deyploying contract..."

cd /ethcode
truffle compile
truffle migrate --network development

echo "[*] Contract deployed"

bash