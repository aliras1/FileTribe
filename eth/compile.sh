truffle compile --reset --all

compile_contracts() {
    IN=$1
    OUT=$2

    mkdir ./${OUT}

    for f in ./${IN}/*.sol
    do
        file_name_with_ext=$(basename ${f})
        filename="${file_name_with_ext%.*}"
        go_filename=$(echo ${filename} | tr '[:upper:]' '[:lower:]')

        echo "./build/contracts/$filename.json"

        mkdir ./${OUT}/${filename}

        cat "./build/contracts/$filename.json" | jq -r '.abi' > .contract.abi
        cat "./build/contracts/$filename.json" | jq -r '.bytecode' > .contract.bin
        abigen -abi .contract.abi -pkg ${filename} --out ./${OUT}/${filename}/${go_filename}.go --bin .contract.bin
        rm .contract.abi
        rm .contract.bin
    done
}

compile_contracts "contracts" "gen"
compile_contracts "contracts/factory" "gen/factory"
compile_contracts "contracts/interfaces" "gen/interfaces"