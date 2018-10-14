hostname -i
address=$(hostname -i)
last_num="${address##*.}"
i=$((2040 + last_num))
echo $i

ipfs init -b $i
#ipfs daemon --enable-pubsub-experiment </dev/null &>/dev/null &
#sleep 15

#cd /root/go/src/ipfs-share/main
#/usr/lib/go-1.9/bin/go build main.go
#mkdir log
#
#echo "[*] Program started"
#
#./main -stderrthreshold=FATAL -log_dir=./log

/bin/sh