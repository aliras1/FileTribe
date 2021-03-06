hostname -i
address=$(hostname -i)
last_num="${address##*.}"
i=$((2040 + last_num))
echo $i

ipfs init -b $i
ipfs daemon --enable-pubsub-experiment </dev/null &>/dev/null &
sleep 15
ipfs config --json Experimental.Libp2pStreamMounting true

cd /mounted
mkdir $last_num

sleep 10

echo "[*] Program started"

./main ./eth.key -stderrthreshold=INFO -log_dir=./$last_num