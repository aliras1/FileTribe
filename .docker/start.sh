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

echo "[*] Program started"

./main -stderrthreshold=INFO -log_dir=./$last_num