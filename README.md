# side-chain

## cmd

```jsonc
This is a side chain application

Usage:
  side [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  init        Initialize side chain
  start       Start side chain node

Flags:
  -h, --help   help for side

Use "side [command] --help" for more information about a command.

```

## gen config

```jsonc
Initialize side chain

Usage:
  side init [flags]

Flags:
      --ceb                Create empty blocks (default true)
  -h, --help               help for init
  -l, --host-list string   Host list, specify hosts for different nodes, separated by semicolons. like 192.168.31.64;192.168.73.2 (default "127.0.0.1")
  -r, --root-dir string    Root directory, '.side-chain' will be generated in the directory you specified, like $HOME/.side-chain (default "./")
  -v, --validators int     Number of Validators (default 1)
```


`root-dir`: Specify the root directory where the configuration files are generated.  
`validators`: Specify how many validators need to be generated, which is related to the data in `genesis.json`

For example, the current instruction`./sc init --test true -v 5`  
All validators will be added to `genesis.json`
```jsonc

ubuntu@ubuntu:./sc init --test true -v 5

I[2024-10-28|17:46:39.540] Default Account                              
Address=0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 
Amount=1000eth 
PrivateKey=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
---------------------------------------------------------

ubuntu@ubuntu:~/.side-chain$ ll

drwx------  7 cloud cloud 4096 10月 25 17:22 ./
drwxr-x--- 50 cloud cloud 4096 10月 25 17:22 ../
drwx------  4 cloud cloud 4096 10月 25 17:22 0/
drwx------  4 cloud cloud 4096 10月 25 17:22 1/
drwx------  4 cloud cloud 4096 10月 25 17:22 2/
drwx------  4 cloud cloud 4096 10月 25 17:22 3/
drwx------  4 cloud cloud 4096 10月 25 17:22 4/
---------------------------------------------------------
ubuntu@ubuntu:~/.side-chain/0$ ll

drwx------ 4 cloud cloud 4096 10月 25 17:22 ./
drwx------ 7 cloud cloud 4096 10月 25 17:22 ../
drwx------ 2 cloud cloud 4096 10月 25 17:22 config/
drwx------ 2 cloud cloud 4096 10月 25 17:22 data/
---------------------------------------------------------
ubuntu@ubuntu:~/.side-chain/0/config$ ll

drwx------ 2 cloud cloud 4096 10月 25 17:22 ./
drwx------ 4 cloud cloud 4096 10月 25 17:22 ../
-rw-r--r-- 1 cloud cloud 3100 10月 25 17:22 config.toml
-rw-r--r-- 1 cloud cloud 1760 10月 25 17:22 genesis.json
-rw------- 1 cloud cloud  148 10月 25 17:22 node_key.json
-rw-r--r-- 1 cloud cloud   81 10月 25 17:22 node.toml
-rw------- 1 cloud cloud  345 10月 25 17:22 priv_validator_key.json
---------------------------------------------------------

{
  "genesis_time": "2024-10-28T09:46:39.570222676Z",
  "chain_id": "side-chain-4ilbdg",
  "initial_height": "0",
  "consensus_params": {
    "block": {
      "max_bytes": "10000000",
      "max_gas": "-1",
      "time_iota_ms": "1000"
    },
    "evidence": {
      "max_age_num_blocks": "100000",
      "max_age_duration": "172800000000000",
      "max_bytes": "1048576"
    },
    "validator": {
      "pub_key_types": [
        "ed25519"
      ]
    },
    "version": {}
  },
  "validators": [
    {
      "address": "8BD9302D565AC0680B8390FD89364F3511050CF2",
      "pub_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "jsBsimMYPVzLjg/Ua5Hp8t4LJ6Ij0dUZWaniFgnJ2n0="
      },
      "power": "10",
      "name": ""
    },
    {
      "address": "11CFC4D6441A4CE8DD6759F86F509FD6ADC0BBEF",
      "pub_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "W0OkPByPeSwKC3VKKgPK26O9W1itjNO2wIfZSsOq7+A="
      },
      "power": "10",
      "name": ""
    },
    {
      "address": "5FB8E69CDCDD713C84115FB7E8D3BC8F84690EC6",
      "pub_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "jNBWnCrEde9bkve2g0DNwUyN0X9QKEleZAcFdzzVTdI="
      },
      "power": "10",
      "name": ""
    },
    {
      "address": "B52D2C6156649D976D1865B231A4A07550C159AF",
      "pub_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "pJPp/k8sqYljOACcaM+399gLJ6YseWQRrlAGgtDbQrI="
      },
      "power": "10",
      "name": ""
    },
    {
      "address": "3F23DFE92DFB9954C2C388E3785B10E025D504A2",
      "pub_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "DXUo39abwMJjoQXhZl9XoxydrRiu4Wcorg7rEzJCc0c="
      },
      "power": "10",
      "name": ""
    }
  ],
  "app_hash": ""
}
```

## start
`sc start --validator-dir ./.side-chain/0` just start node

```jsonc
Start side chain node

Usage:
  side start [flags]

Flags:
  -h, --help                   help for start
  -v, --validator-dir string   Node directory (default "./.side-chain/0")
```

## mint

`./sc mint`: mint 1000ether to `0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266`

node log
```jsonc
D[2024-10-28|22:50:48.797] UpdateAccountNonce                           Address=0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 Nonce=1 file=/home/cloud/yunyc12345/side-chain/core/service/db.go line=53
D[2024-10-28|22:50:48.797] UpdateAccountBalance                         Address=0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 Balance=1000000000000000000000 file=/home/cloud/yunyc12345/side-chain/core/service/db.go line=39
```

## query


`./sc query balance`: query balance
```jsonc
I[2024-10-28|22:51:01.124] account balance                              balance=1000000000000000000000
```

`./sc query nonce`: query nonce
```jsonc
I[2024-10-28|22:51:04.983] account nonce                                nonce=1
```

## test tx

Calculating gas: `go test -v -run TestCheckTx ./core/types/test/tx_test.go -args -ltdp {filepath}`
```jsonc
=== RUN   TestCheckTx
    tx_test.go:190: flag value:  /home/cloud/Downloads/aa
    tx_test.go:197: data len:  5001553
Response Status: 200 OK
Response Body: {"jsonrpc":"2.0","id":1,"result":{"code":0,"data":null,"log":"","info":"","gas_wanted":"50015530","gas_used":"0","events":[],"codespace":"","sender":"","priority":"0","mempoolError":""}}
--- PASS: TestCheckTx (0.44s)
PASS
ok      command-line-arguments  0.523s

```

Send blob tx: `go test -v -run TestSendLargeTx ./core/types/test/tx_test.go -args -ltdp {filepath}`

```jsonc
=== RUN   TestSendLargeTx
    tx_test.go:126: flag value:  /home/cloud/Downloads/aa
    tx_test.go:133: data len:  5001553
Response Status: 200 OK
Response Body: {"jsonrpc":"2.0","id":1,"result":{"code":0,"data":"","log":"","codespace":"","hash":"F7F202123826CAF85DD8713BA45BA5D49BECF35021113B0DDDE0DD98766DEC9D"}}
--- PASS: TestSendLargeTx (0.46s)
PASS
ok      command-line-arguments  0.538s
```

## log

the current log settings are as follows,   
there is currently no cmd configurable operation available.  
If the root directory of the node is `./.side-chain/0`, the log is in `./.side-chain/0/log/node.log`
```jsonc
	LogMaxSize    = 128 // mb
	LogMaxAge     = 30  // day
	LogMaxBackups = 100 //
	LogPath       = "log/node.log"
	
```