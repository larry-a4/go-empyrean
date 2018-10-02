---
title: Shyft Block Explorer Documentation 

language_tabs: # must be one of https://git.io/vQNgJ
  - shell: cURL

toc_footers:
  - <a href='https://github.com/lord/slate'>Documentation Powered by Slate</a>

includes:
  - errors

search: true
---

# Introduction

### A Note From The Developers

Our goal was to avoid compromising the integrity of Geth and simply extend existing functionality to meet the specific needs of the Shyft Network. To our utmost ability we have documented, within the codebase, exactly where we have extended our functionality using the following notation:  NOTE:SHYFT. This document is meant to provide a high level overview of the changes made to Geth and to provide explanations, where needed, on the changes that were made. Another benefit of this document is to allow others to quickly see the changes that were made in order to get quicker feedback on a compromising line of code.

### Contributing To Shyft Geth

In order to successfully accept a PR the maintainers of the Shyft repositories require that this document must be updated, reflecting the changes made in the PR. Along with the documentation, we ask that contributors provide the NOTE:SHYFT. The tag could should contain a brief on the modified code. This will help with releases further down the road as we document what breaking changes have been made along the journey.

# Setup

### Dependencies
    
 - go 1.10
 - postgres 10
    
To install go please review the installation docs [here](https://golang.org/doc/install), but ensure you download version 1.10. If you would like to install go with a script please check out this repo [here](https://github.com/canha/golang-tools-install-script).
    
To install postgres please review the installation docs [here](https://www.postgresql.org/docs/10/static/tutorial-install.html).

### Govendor and Packages/Dependencies

> Download Go Vendor

```shell
go get -u github.com/kardianos/govendor
```

> To run govendor globally, have this in your bash_profile file:

```shell
export GOPATH=$HOME/go
export PATH=$PATH:$HOME/go/bin
```

> Then go_empyrean will need to be cloned to this directory:

```shell
$GOPATH/src/github.com/ShyftNetwork/
```

Geth uses govendor to manage packages/dependencies: [Go Vendor](https://github.com/kardianos/govendor)

This has some more information: [Ethereum Wiki](https://github.com/ethereum/go-ethereum/wiki/Developers'-Guide)

To add a new dependency, run govendor fetch <import-path> , and commit the changes to git. Then the deps will be accessible on other machines that pull from git.

<aside class="notice">
GOPATH is not strictly necessary however, for govendor it is much easier to use gopath as go will look for binaries in this directory ($GOPATH/bin). To set up GOPATH, read the govendor section.
</aside>

### Running Locally

To begin running locally, please ensure you have correctly installed go 1.10 and postgres (make sure postgres is running). 
Once cloned, in a terminal window run the following command:

Before running any CLI options ensure you run **`make geth`** in the root directory.

``.shyft-geth.sh --setup`` This sets up postgres and the shyft chain db

``./shyft-geth.sh --start`` This starts GETH

At this point you should see GETH running in the terminal and if you opened your postgres instance you should see data being populated into the tables. It might look something similiar to the image below.

<img src="./images/geth.png" alt="Geth example">

To stop Geth, **`crtl+C`** in the terminal window, if you proceed with the start script mentioned above the Shyft chain will begin from the last block height, if you wish to start the chain fresh from genesis follow the below steps:

``./shyft-geth.sh --reset`` This drops postgres and chaindb data

``./shyft-geth.sh --start`` Starts GETH

To see transactions being submitted on the network see the sendTransactions command in the CLI section of this readme.
### Docker Images

Docker Images are available for ShyftGeth and the Postgresql Database which can be used for development and testing. To launch these containers you will need to have docker-compose installed on your computer. Installation instructions for docker-compose are available [here](https://docs.docker.com/install/).

**To build the images for the first time please run the following command:**

`./shyft-geth.sh --setup # clears persisted directories prior to docker build`

`docker-compose up --build`

If you would like to reinitialize/rebuild the docker images you can run the above mentioned command as well.

To launch ShyftGeth, PG, the ShyftBlock Explorer Api and UI anytime after initial build - issue the following commands from the root of the project directory:

`./shyft-geth.sh --setup # clears persisted directories prior to docker build`

**`docker-compose up`**

To stop/pause mining - enter:

**`docker-compose stop`**

And then just issue `docker-compose up` to continue mining.
### Docker Postgresql - DB Connection
From your local machine you can view the database by connecting to the database in the container at 
**``127.0.0.1:8001``**

Use User: 'postgres' Password: 'docker' Database: 'shyftdb'
### Docker Block Explorer Api 
To access the shyftBlockExplorer open a browser and visit 

**``http://localhost:3000``**

To rebuild any one of the services- issue the following commands:

Services:

   - ShyftGeth
   - Postgres Instance
   - Shyft Explorer API
   - Shyft Example Explorer UI

**``
docker-compose up -d --no-deps --build <docker compose file service name> 
``**

ie. for shyftBlockExplorerApi:

**``docker-compose up -d --no-deps --build shyft_block_api``**

The Postgresql Database Container will persist the database data to the directory ``./pg-data`` _. So if you do want to reinitialize the database you should delete this directory as well as the blockchain data directories ``(./shyftData ./privatenet)`` prior to launching the docker containers. There is a shell script available to delete these folders to run it execute the following command:

**``./shyft-cli/resetShyftGeth.sh``**

Blockchain data is persisted to **``./ethash/.ethash and ./shyftData__``**. If you would like to reset the test blockchain you will need to delete the **``__./ethash ./shyftData & ./privatenet__``** directories.

The docker container for the ShyftBlockExplorerApi utilizes govendor to minimize its image size. **If you would like the docker image for this container to reflect any uncommitted changes which may have occurred in the go-empyrean repository, ie. changes with respect to go-empyrean core (ie. cryptographic functions and database). Prior to launching the docker containers you should rebuild the vendor directory for the shyftBlockExplorerApi - by executing the following steps:**

Remove existing shyftBlockExplorerApi vendor.json and vendored components:

**``rm -rf shyftBlockExplorerApi/vendor``**

reinitialize vendor.json

**``cd shyftBlockExplorerApi && govendor init``**

rebuild vendor.json using latest uncommitted changes

**``govendor add +external``**

Due to a bug in govendor and it not being able to pull in some dependencies that are c-header files 
you should execute the following commands - see these issues - which whilst closed
appears to have not been fixed: https://github.com/kardianos/govendor/issues/124 && https://github.com/kardianos/govendor/issues/61

**``govendor remove github.com/ShyftNetwork/go-empyrean/crypto/secp256k1/^``**

**``govendor fetch github.com/ShyftNetwork/go-empyrean/crypto/secp256k1/^``**

NB: The Shyft Geth docker image size is 1+ GB so make sure you have adequate space on your disk drive/

# Shyft BlockExplorer API

In order to store the block explorer database, a custom folder was created `./shyft_schema` that contains all the necessary functions to read and write to the explorer database.

The main functions exist in `./core/shyft_database_util.go` and `./core/shyft_get_utils.go`

To run the block explorer rest api that queries the postgres instance and returns a json body, open a new terminal window, navigate to the root directory of the project and run the following command:

**``go run blockExplorerApi/*.go``**

This will start a go server on port 8080 and allow you to either run the pre-existing block explorer or query the api endpoints. Its important to note, that if you have nothing in your postgres database the API will return nothing.

## Blocks

<span style="color:#fff; height: 25px; width: 30px; background: blue; padding: 7px; margin-right: 5px;">**GET**</span> ``/api/get_block/{blockNumber}`` 

<span style="color:#fff; height: 25px; width: 30px; background: red; padding: 7px; margin-right: 5px;">**Headers**</span> ``Content-Type: application/json``

**Description**

```json
{  
   "Hash":"0x72590ac6e7626b9b1f77452d83297c0361e6ff7fa011872289224ee02b9acc8f",
   "Coinbase":"0x43ec6d0942f7faef069f7f63d0384a27f529b062",
   "Age":"2018-10-01T11:15:54Z",
   "ParentHash":"0x5c1aa0559093d9e1e128b6d8ae63d5bf9a5fbf273f704afad9adac88e734c3cc",
   "UncleHash":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
   "Difficulty":"135005",
   "Size":"538.00 B",
   "Rewards":"5000000000000000000",
   "Number":63,
   "GasUsed":0,
   "GasLimit":25534458,
   "Nonce":7333650700872754740,
   "TxCount":0,
   "UncleCount":0
}
```

The above endpoint will respond with a the block data for that specific block number. The table below lists the attributes that will be included in the response from this endpoint.

| Parameter | Description |
| --- | --- |
| `blockNumber` | The block height or block number |

| Attributes | Description                                  |    Type  |
| :---       |     :---:                                    |    ---:  |
| `Hash`     | Block hash                                   | string   |
| `Coinbase` | Address of miner                             | string   |
| `Number`   | Block height                                 | number   |
| `GasUsed`  | Amount of gas used                           | number   |
| `GasLimit` | Maximum amount of gas willing to be spent    | number   |
| `TxCount`  | Amount of transactions included in the block | number   |
| `UncleCount`| Amount of uncle blocks                      | number   |
| `Age`       | Time stamp of block creation                | Timestamp|
| `ParentHash`| Hash of the prior block                     | string   |
| `UncleHash` | Hash of a the uncle block                   | string   |
| `Size`      | Size of block measured in Bytes             | string   |
| `Rewards`   | Block reward                                | string   |
| `Nonce`     | Value used for PoW                          | number   |
| `Difficulty`| The difficulty for this block               | string   |


<span style="color:#fff; height: 25px; width: 30px; background: blue; padding: 7px; margin-right: 5px;">**GET**</span> ``/api/get_all_blocks`` 

<span style="color:#fff; height: 25px; width: 30px; background: red; padding: 7px; margin-right: 5px;">**Headers**</span> ``Content-Type: application/json``

**Description**

```json
[{  
   "Hash":"0x5c1aa0559093d9e1e128b6d8ae63d5bf9a5fbf273f704afad9adac88e734c3cc",
   "Coinbase":"0x43ec6d0942f7faef069f7f63d0384a27f529b062",
   "Age":"2018-10-01T11:15:53Z",
   "ParentHash":"0xb0437e25e7cfb113bbea0f014883e8f441d5465bb628f9a112eb372e43055f1a",
   "UncleHash":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
   "Difficulty":"134940",
   "Size":"538.00 B",
   "Rewards":"5000000000000000000",
   "Number":62,
   "GasUsed":0,
   "GasLimit":25509548,
   "Nonce":2307164919004664188,
   "TxCount":0,
   "UncleCount":0
},
{  
   "Hash":"0x72590ac6e7626b9b1f77452d83297c0361e6ff7fa011872289224ee02b9acc8f",
   "Coinbase":"0x43ec6d0942f7faef069f7f63d0384a27f529b062",
   "Age":"2018-10-01T11:15:54Z",
   "ParentHash":"0x5c1aa0559093d9e1e128b6d8ae63d5bf9a5fbf273f704afad9adac88e734c3cc",
   "UncleHash":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
   "Difficulty":"135005",
   "Size":"538.00 B",
   "Rewards":"5000000000000000000",
   "Number":63,
   "GasUsed":0,
   "GasLimit":25534458,
   "Nonce":7333650700872754740,
   "TxCount":0,
   "UncleCount":0
}]
```

The above endpoint will respond with a the block data for all blocks in the postgres db. The table below lists the attributes that will be included in the response from this endpoint.

*No Parameter requirement*

| Attributes | Description                                  |    Type  |
| :---       |     :---:                                    |    ---:  |
| `Hash`     | Block hash                                   | string   |
| `Coinbase` | Address of miner                             | string   |
| `Number`   | Block height                                 | number   |
| `GasUsed`  | Amount of gas used                           | number   |
| `GasLimit` | Maximum amount of gas willing to be spent    | number   |
| `TxCount`  | Amount of transactions included in the block | number   |
| `UncleCount`| Amount of uncle blocks                      | number   |
| `Age`       | Time stamp of block creation                | Timestamp|
| `ParentHash`| Hash of the prior block                     | string   |
| `UncleHash` | Hash of a the uncle block                   | string   |
| `Size`      | Size of block measured in Bytes             | string   |
| `Rewards`   | Block reward                                | string   |
| `Nonce`     | Value used for PoW                          | number   |
| `Difficulty`| The difficulty for this block               | string   |


<span style="color:#fff; height: 25px; width: 30px; background: blue; padding: 7px; margin-right: 5px;">**GET**</span> ``/api/get_recent_block`` 

<span style="color:#fff; height: 25px; width: 30px; background: red; padding: 7px; margin-right: 5px;">**Headers**</span> ``Content-Type: application/json``

**Description**

```json
{  
   "Hash":"0x72590ac6e7626b9b1f77452d83297c0361e6ff7fa011872289224ee02b9acc8f",
   "Coinbase":"0x43ec6d0942f7faef069f7f63d0384a27f529b062",
   "Age":"2018-10-01T11:15:54Z",
   "ParentHash":"0x5c1aa0559093d9e1e128b6d8ae63d5bf9a5fbf273f704afad9adac88e734c3cc",
   "UncleHash":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
   "Difficulty":"135005",
   "Size":"538.00 B",
   "Rewards":"5000000000000000000",
   "Number":63,
   "GasUsed":0,
   "GasLimit":25534458,
   "Nonce":7333650700872754740,
   "TxCount":0,
   "UncleCount":0
}
```

The above endpoint will respond with the block data from the highest block height. The table below lists the attributes that will be included in the response from this endpoint.

*No Parameter requirement*

| Attributes | Description                                  |    Type  |
| :---       |     :---:                                    |    ---:  |
| `Hash`     | Block hash                                   | string   |
| `Coinbase` | Address of miner                             | string   |
| `Number`   | Block height                                 | number   |
| `GasUsed`  | Amount of gas used                           | number   |
| `GasLimit` | Maximum amount of gas willing to be spent    | number   |
| `TxCount`  | Amount of transactions included in the block | number   |
| `UncleCount`| Amount of uncle blocks                      | number   |
| `Age`       | Time stamp of block creation                | Timestamp|
| `ParentHash`| Hash of the prior block                     | string   |
| `UncleHash` | Hash of a the uncle block                   | string   |
| `Size`      | Size of block measured in Bytes             | string   |
| `Rewards`   | Block reward                                | string   |
| `Nonce`     | Value used for PoW                          | number   |
| `Difficulty`| The difficulty for this block               | string   |


<span style="color:#fff; height: 25px; width: 30px; background: blue; padding: 7px; margin-right: 5px;">**GET**</span> ``/api/get_blocks_mined/{coinbase}`` 

<span style="color:#fff; height: 25px; width: 30px; background: red; padding: 7px; margin-right: 5px;">**Headers**</span> ``Content-Type: application/json``

**Description**

```json
[{  
   "Hash":"0x5c1aa0559093d9e1e128b6d8ae63d5bf9a5fbf273f704afad9adac88e734c3cc",
   "Coinbase":"0x43ec6d0942f7faef069f7f63d0384a27f529b062",
   "Age":"2018-10-01T11:15:53Z",
   "ParentHash":"0xb0437e25e7cfb113bbea0f014883e8f441d5465bb628f9a112eb372e43055f1a",
   "UncleHash":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
   "Difficulty":"134940",
   "Size":"538.00 B",
   "Rewards":"5000000000000000000",
   "Number":62,
   "GasUsed":0,
   "GasLimit":25509548,
   "Nonce":2307164919004664188,
   "TxCount":0,
   "UncleCount":0
},
{  
   "Hash":"0x72590ac6e7626b9b1f77452d83297c0361e6ff7fa011872289224ee02b9acc8f",
   "Coinbase":"0x43ec6d0942f7faef069f7f63d0384a27f529b062",
   "Age":"2018-10-01T11:15:54Z",
   "ParentHash":"0x5c1aa0559093d9e1e128b6d8ae63d5bf9a5fbf273f704afad9adac88e734c3cc",
   "UncleHash":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
   "Difficulty":"135005",
   "Size":"538.00 B",
   "Rewards":"5000000000000000000",
   "Number":63,
   "GasUsed":0,
   "GasLimit":25534458,
   "Nonce":7333650700872754740,
   "TxCount":0,
   "UncleCount":0
}]
```

The above endpoint will respond with a the block data for all blocks which have been mined by the provided address. The table below lists the attributes that will be included in the response from this endpoint.

| Parameter | Description |
| --- | --- |
| `coinbase` | The miners address |

| Attributes | Description                                  |    Type  |
| :---       |     :---:                                    |    ---:  |
| `Hash`     | Block hash                                   | string   |
| `Coinbase` | Address of miner                             | string   |
| `Number`   | Block height                                 | number   |
| `GasUsed`  | Amount of gas used                           | number   |
| `GasLimit` | Maximum amount of gas willing to be spent    | number   |
| `TxCount`  | Amount of transactions included in the block | number   |
| `UncleCount`| Amount of uncle blocks                      | number   |
| `Age`       | Time stamp of block creation                | Timestamp|
| `ParentHash`| Hash of the prior block                     | string   |
| `UncleHash` | Hash of a the uncle block                   | string   |
| `Size`      | Size of block measured in Bytes             | string   |
| `Rewards`   | Block reward                                | string   |
| `Nonce`     | Value used for PoW                          | number   |
| `Difficulty`| The difficulty for this block               | string   |


## Transactions

<span style="color:#fff; height: 25px; width: 30px; background: blue; padding: 7px; margin-right: 5px;">**GET**</span> ``/api/get_transaction/{txHash}`` 

<span style="color:#fff; height: 25px; width: 30px; background: red; padding: 7px; margin-right: 5px;">**Headers**</span> ``Content-Type: application/json``

**Description**

```json
{  
   "TxHash":"0x5bd738164c61fb50eb12e227846cbaef2de965aa0x5bd738164c61fb50eb12e227846cbaef2de965aa",
   "To":"0x5bd738164c61fb50eb12e227846cbaef2de965aa",
   "From":"0x007622d84a234bb8b078230fcf84b67ae9a8acae",
   "BlockHash":"0xfa01942529ef3a4e543ef8c061c6e0cb69a61c489d3bb6891bb65651c02dafd4",
   "BlockNumber":"2",
   "Amount":"400000000000000000000",
   "GasPrice":1253,
   "Gas":12124,
   "GasLimit":24011655,
   "Cost":"53002",
   "Nonce":2,
   "Status":"SUCCESS",
   "IsContract":false,
   "Age":"2018-03-18T19:38:41Z",
   "Data":""
}
```

The above endpoint will respond with transaction data from the provided transaction hash. The table below lists the attributes that will be included in the response from this endpoint.

| Parameter | Description |
| --- | --- |
| `txHash` | The hash for that particular transaction |

| Attributes  | Description                                     |    Type  |
| :---        |     :---:                                       |    ---:  |
| `TxHash`    | Transaction hash                                | string   |
| `To`        | Address of transaction receiver                 | string   |
| `From`      | Address of transaction sender                   | string   |
| `BlockHash` | Hash of block                                   | string   |
| `BlockNumber`| Block height                                   | string   |
| `Amount`    | Amount of value being transferred               | string   |
| `GasPrice`  | Price of required gas                           | number   |
| `Gas`       | Required pricing value to process transaction   | number   |
| `GasLimit`  | Maximum amount of gas willing to be spent       | number   |
| `Cost`      | Hash of a the uncle block                       | string   |
| `Nonce`     | Number of transactions sent from a given address| number   |
| `Status`    | Whether the transaction was success or fail     | string   |
| `IsContract`| Whether the transaction was from a contract     | bool     |
| `Age`       | Time stamp of transaction creation              | timestamp|
| `Data`      | Contract data in byte code                      | byteArray|


<span style="color:#fff; height: 25px; width: 30px; background: blue; padding: 7px; margin-right: 5px;">**GET**</span> ``/api/get_all_transactions`` 

<span style="color:#fff; height: 25px; width: 30px; background: red; padding: 7px; margin-right: 5px;">**Headers**</span> ``Content-Type: application/json``

**Description**

```json
[{  
   "TxHash":"0x5bd738164c61fb50eb12e227846cbaef2de965aa0x5bd738164c61fb50eb12e227846cbaef2de965aa",
   "To":"0x5bd738164c61fb50eb12e227846cbaef2de965aa",
   "From":"0x007622d84a234bb8b078230fcf84b67ae9a8acae",
   "BlockHash":"0xfa01942529ef3a4e543ef8c061c6e0cb69a61c489d3bb6891bb65651c02dafd4",
   "BlockNumber":"2",
   "Amount":"400000000000000000000",
   "GasPrice":1253,
   "Gas":12124,
   "GasLimit":24011655,
   "Cost":"53002",
   "Nonce":2,
   "Status":"SUCCESS",
   "IsContract":false,
   "Age":"2018-03-18T19:38:41Z",
   "Data":""
},
{  
   "TxHash":"0x2f56g38164c634550eb12e222146cbaef2de965aa0x5bd738164c61fb50eb12e227846cbaef2de965aa",
   "To":"0x5bd738164c61fb50eb12e227846cbaef2de965aa",
   "From":"0x007622d84a234bb8b078230fcf84b67ae9a8acae",
   "BlockHash":"0xbg01234529ef3a4e543q23061c6e0cb69a61c489d3bb6891bb65651c02dafd4",
   "BlockNumber":"3",
   "Amount":"100000000000000000000",
   "GasPrice":1253,
   "Gas":12124,
   "GasLimit":24011655,
   "Cost":"53002",
   "Nonce":3,
   "Status":"SUCCESS",
   "IsContract":false,
   "Age":"2018-03-18T19:40:41Z",
   "Data":""
}]

```

The above endpoint will respond with transaction data for all transactions in the postgres database. The table below lists the attributes that will be included in the response from this endpoint.

*No parameters required*

| Attributes  | Description                                     |    Type  |
| :---        |     :---:                                       |    ---:  |
| `TxHash`    | Transaction hash                                | string   |
| `To`        | Address of transaction receiver                 | string   |
| `From`      | Address of transaction sender                   | string   |
| `BlockHash` | Hash of block                                   | string   |
| `BlockNumber`| Block height                                   | string   |
| `Amount`    | Amount of value being transferred               | string   |
| `GasPrice`  | Price of required gas                           | number   |
| `Gas`       | Required pricing value to process transaction   | number   |
| `GasLimit`  | Maximum amount of gas willing to be spent       | number   |
| `Cost`      | Hash of a the uncle block                       | string   |
| `Nonce`     | Number of transactions sent from a given address| number   |
| `Status`    | Whether the transaction was success or fail     | string   |
| `IsContract`| Whether the transaction was from a contract     | bool     |
| `Age`       | Time stamp of transaction creation              | timestamp|
| `Data`      | Contract data in byte code                      | byteArray|

<span style="color:#fff; height: 25px; width: 30px; background: blue; padding: 7px; margin-right: 5px;">**GET**</span> ``/api/get_all_transactions_from_block/{blockNumber}`` 

<span style="color:#fff; height: 25px; width: 30px; background: red; padding: 7px; margin-right: 5px;">**Headers**</span> ``Content-Type: application/json``

**Description**

```json
[{  
   "TxHash":"0x5bd738164c61fb50eb12e227846cbaef2de965aa0x5bd738164c61fb50eb12e227846cbaef2de965aa",
   "To":"0x5bd738164c61fb50eb12e227846cbaef2de965aa",
   "From":"0x007622d84a234bb8b078230fcf84b67ae9a8acae",
   "BlockHash":"0xfa01942529ef3a4e543ef8c061c6e0cb69a61c489d3bb6891bb65651c02dafd4",
   "BlockNumber":"3",
   "Amount":"400000000000000000000",
   "GasPrice":1253,
   "Gas":12124,
   "GasLimit":24011655,
   "Cost":"53002",
   "Nonce":3,
   "Status":"SUCCESS",
   "IsContract":false,
   "Age":"2018-03-18T19:38:41Z",
   "Data":""
},
{  
   "TxHash":"0x2f56g38164c634550eb12e222146cbaef2de965aa0x5bd738164c61fb50eb12e227846cbaef2de965aa",
   "To":"0x5bd738164c61fb50eb12e227846cbaef2de965aa",
   "From":"0x007622d84a234bb8b078230fcf84b67ae9a8acae",
   "BlockHash":"0xbg01234529ef3a4e543q23061c6e0cb69a61c489d3bb6891bb65651c02dafd4",
   "BlockNumber":"3",
   "Amount":"100000000000000000000",
   "GasPrice":1253,
   "Gas":12124,
   "GasLimit":24011655,
   "Cost":"53002",
   "Nonce":4,
   "Status":"SUCCESS",
   "IsContract":false,
   "Age":"2018-03-18T19:40:41Z",
   "Data":""
}]

```

The above endpoint will respond with transaction data for all transactions in the postgres database. The table below lists the attributes that will be included in the response from this endpoint.

| Parameter | Description |
| --- | --- |
| `blockNumber` | The block height or block number |

| Attributes  | Description                                     |    Type  |
| :---        |     :---:                                       |    ---:  |
| `TxHash`    | Transaction hash                                | string   |
| `To`        | Address of transaction receiver                 | string   |
| `From`      | Address of transaction sender                   | string   |
| `BlockHash` | Hash of block                                   | string   |
| `BlockNumber`| Block height                                   | string   |
| `Amount`    | Amount of value being transferred               | string   |
| `GasPrice`  | Price of required gas                           | number   |
| `Gas`       | Required pricing value to process transaction   | number   |
| `GasLimit`  | Maximum amount of gas willing to be spent       | number   |
| `Cost`      | Hash of a the uncle block                       | string   |
| `Nonce`     | Number of transactions sent from a given address| number   |
| `Status`    | Whether the transaction was success or fail     | string   |
| `IsContract`| Whether the transaction was from a contract     | bool     |
| `Age`       | Time stamp of transaction creation              | timestamp|
| `Data`      | Contract data in byte code                      | byteArray|

<span style="color:#fff; height: 25px; width: 30px; background: blue; padding: 7px; margin-right: 5px;">**GET**</span> ``/api/get_account_txs/{address}`` 

<span style="color:#fff; height: 25px; width: 30px; background: red; padding: 7px; margin-right: 5px;">**Headers**</span> ``Content-Type: application/json``

**Description**

```json
[  
   {  
      "TxHash":"0x7da80aaf6f7e382735310c725b81f790f84c75a541a5360ecba30eb2d7965395",
      "To":"0xfaeccae8c1af6bdfd71095e1b6a2f61c61c8a7e7",
      "From":"0x43ec6d0942f7faef069f7f63d0384a27f529b062",
      "BlockHash":"0x894459a52417bcf3e0eec81fce6cc0858813cbd2202711031e8c5185d9aa4d38",
      "BlockNumber":"2",
      "Amount":"0",
      "GasPrice":100000000000,
      "Gas":6721975,
      "GasLimit":24058572,
      "Cost":"672197500000000000",
      "Nonce":0,
      "Status":"SUCCESS",
      "IsContract":true,
      "Age":"2018-10-02T13:21:02Z",
      "Data":"YIBgQFI0gBVhABBXYAGCY/////8WfAEAAAAAAAAAAAAAAAAAAAAAAAYQVienpyMFS8PtYkAKQ=="
   },
   {  
      "TxHash":"0x2f49282ff117dd4d28f4b5a59a71fc0c3ab1cd595f6d503cd82caa88d1cc4897",
      "To":"0xfaeccae8c1af6bdfd71095e1b6a2f61c61c8a7e7",
      "From":"0x43ec6d0942f7faef069f7f63d0384a27f529b062",
      "BlockHash":"0x12b73a2ae27cb32193e8570c275b24d6a9aa9d637c61582e63afe589535f075a",
      "BlockNumber":"4",
      "Amount":"0",
      "GasPrice":100000000000,
      "Gas":6721975,
      "GasLimit":24105581,
      "Cost":"672197500000000000",
      "Nonce":1,
      "Status":"SUCCESS",
      "IsContract":false,
      "Age":"2018-10-02T13:21:05Z",
      "Data":"/azVdgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAB"
   }
]

```

The above endpoint will respond with transaction data for all transactions conducted by the provided address. The table below lists the attributes that will be included in the response from this endpoint.

| Parameter | Description |
| --- | --- |
| `address` | The address of an account |

| Attributes  | Description                                     |    Type  |
| :---        |     :---:                                       |    ---:  |
| `TxHash`    | Transaction hash                                | string   |
| `To`        | Address of transaction receiver                 | string   |
| `From`      | Address of transaction sender                   | string   |
| `BlockHash` | Hash of block                                   | string   |
| `BlockNumber`| Block height                                   | string   |
| `Amount`    | Amount of value being transferred               | string   |
| `GasPrice`  | Price of required gas                           | number   |
| `Gas`       | Required pricing value to process transaction   | number   |
| `GasLimit`  | Maximum amount of gas willing to be spent       | number   |
| `Cost`      | Hash of a the uncle block                       | string   |
| `Nonce`     | Number of transactions sent from a given address| number   |
| `Status`    | Whether the transaction was success or fail     | string   |
| `IsContract`| Whether the transaction was from a contract     | bool     |
| `Age`       | Time stamp of transaction creation              | timestamp|
| `Data`      | Contract data in byte code                      | byteArray|

<span style="color:#fff; height: 25px; width: 30px; background: blue; padding: 7px; margin-right: 5px;">**GET**</span> ``/api/get_internal_transactions/{hash}`` 

<span style="color:#fff; height: 25px; width: 30px; background: red; padding: 7px; margin-right: 5px;">**Headers**</span> ``Content-Type: application/json``

**Description**

```json
[  
   {  
      "ID":1,
      "Hash":"0x2f49282ff117dd4d28f4b5a59a71fc0c3ab1cd595f6d503cd82caa88d1cc4897",
      "BlockHash":"0x12b73a2ae27cb32193e8570c275b24d6a9aa9d637c61582e63afe589535f075a",
      "Action":"CALL",
      "From":"0x43ec6d0942f7faef069f7f63d0384a27f529b062",
      "To":"0xfaeccae8c1af6bdfd71095e1b6a2f61c61c8a7e7",
      "Value":"0",
      "Gas":6700511,
      "GasUsed":20544,
      "Input":"0xfdacd5760000000000000000000000000000000000000000000000000000000000000001",
      "Output":"0x",
      "Time":"6.950856ms"
   }
]
```

The above endpoint will respond with transaction data from the provided transaction hash. The table below lists the attributes that will be included in the response from this endpoint.

| Parameter | Description |
| --- | --- |
| `hash` | The transaction hash of the particular record |

| Attributes  | Description                                     |    Type  |
| :---        |     :---:                                       |    ---:  |
| `ID`        | Incrementing integer for unique db record       | number   |
| `Hash`      | Transaction hash                                | string   |
| `BlockHash` | Hash of block                                   | string   |
| `Action`    | Contract type, can be Call or Create            | string   |
| `To`        | Address of transaction receiver                 | string   |
| `From`      | Address of transaction sender                   | string   |
| `Value`     | Amount of value being transferred               | string   |
| `Gas`       | Required pricing value to process transaction   | number   |
| `GasUsed`   | Amount of gas used                              | number   |
| `Input`     | Parameters passed into contract in bytecode     | string   |
| `Output`    | Return values from contract in bytecode         | string   |
| `Time`      | Amount of time to trace transaction             | string   |

<span style="color:#fff; height: 25px; width: 30px; background: blue; padding: 7px; margin-right: 5px;">**GET**</span> ``/api/get_internal_transactions`` 

<span style="color:#fff; height: 25px; width: 30px; background: red; padding: 7px; margin-right: 5px;">**Headers**</span> ``Content-Type: application/json``

**Description**

```json
[  
   {  
      "ID":1,
      "Hash":"0x2f49282ff117dd4d28f4b5a59a71fc0c3ab1cd595f6d503cd82caa88d1cc4897",
      "BlockHash":"0x12b73a2ae27cb32193e8570c275b24d6a9aa9d637c61582e63afe589535f075a",
      "Action":"CALL",
      "From":"0x43ec6d0942f7faef069f7f63d0384a27f529b062",
      "To":"0xfaeccae8c1af6bdfd71095e1b6a2f61c61c8a7e7",
      "Value":"0",
      "Gas":6700511,
      "GasUsed":20544,
      "Input":"0xfdacd5760000000000000000000000000000000000000000000000000000000000000001",
      "Output":"0x",
      "Time":"6.950856ms"
   },
   {  
      "ID":2,
      "Hash":"0x16b218a28d48c5ef54cf09ce836e285cb23ba9c179b3c88269d61cb5bd5473db",
      "BlockHash":"0xbdb1aca7c9aeee9f9723ea84cb8d47ff46f49e29a85d41a92b7950a0970b1d25",
      "Action":"CALL",
      "From":"0x43ec6d0942f7faef069f7f63d0384a27f529b062",
      "To":"0xfaeccae8c1af6bdfd71095e1b6a2f61c61c8a7e7",
      "Value":"0",
      "Gas":6700511,
      "GasUsed":5544,
      "Input":"0xfdacd5760000000000000000000000000000000000000000000000000000000000000002",
      "Output":"0x",
      "Time":"1.370835ms"
   }
]
```

The above endpoint will respond with all internal transaction data from the postgres database. The table below lists the attributes that will be included in the response from this endpoint.

*No parameters requirements*

| Attributes  | Description                                     |    Type  |
| :---        |     :---:                                       |    ---:  |
| `ID`        | Incrementing integer for unique db record       | number   |
| `Hash`      | Transaction hash                                | string   |
| `BlockHash` | Hash of block                                   | string   |
| `Action`    | Contract type, can be Call or Create            | string   |
| `To`        | Address of transaction receiver                 | string   |
| `From`      | Address of transaction sender                   | string   |
| `Value`     | Amount of value being transferred               | string   |
| `Gas`       | Required pricing value to process transaction   | number   |
| `GasUsed`   | Amount of gas used                              | number   |
| `Input`     | Parameters passed into contract in bytecode     | string   |
| `Output`    | Return values from contract in bytecode         | string   |
| `Time`      | Amount of time to trace transaction             | string   |


## Accounts

<span style="color:#fff; height: 25px; width: 30px; background: blue; padding: 7px; margin-right: 5px;">**GET**</span> ``/api/get_account/{address}`` 

<span style="color:#fff; height: 25px; width: 30px; background: red; padding: 7px; margin-right: 5px;">**Headers**</span> ``Content-Type: application/json``

**Description**

```json
{  
   "Addr":"0xfaeccae8c1af6bdfd71095e1b6a2f61c61c8a7e7",
   "Balance":"100000000",
   "AccountNonce":"5"
}
```

The above endpoint will respond with account data from the provided address. The table below lists the attributes that will be included in the response from this endpoint.

| Parameter | Description |
| --- | --- |
| `address` | The address of an account to query|

| Attributes    | Description                                     |    Type  |
| :---          |     :---:                                       |    ---:  |
| `Addr`        | Account address                                 | string   |
| `Balance`     | Address balance in wei                          | string   |
| `AccountNonce`| Number of transactions sent from a given address| string   |

<span style="color:#fff; height: 25px; width: 30px; background: blue; padding: 7px; margin-right: 5px;">**GET**</span> ``/api/get_account/{address}`` 

<span style="color:#fff; height: 25px; width: 30px; background: red; padding: 7px; margin-right: 5px;">**Headers**</span> ``Content-Type: application/json``

**Description**

```json
{  
   "Addr":"0xfaeccae8c1af6bdfd71095e1b6a2f61c61c8a7e7",
   "Balance":"100000000",
   "AccountNonce":"5"
}
```

The above endpoint will respond with account data from the provided address. The table below lists the attributes that will be included in the response from this endpoint.

| Parameter | Description |
| --- | --- |
| `address` | The address of an account to query|

| Attributes    | Description                                     |    Type  |
| :---          |     :---:                                       |    ---:  |
| `Addr`        | Account address                                 | string   |
| `Balance`     | Address balance in wei                          | string   |
| `AccountNonce`| Number of transactions sent from a given address| string   |

<span style="color:#fff; height: 25px; width: 30px; background: blue; padding: 7px; margin-right: 5px;">**GET**</span> ``/api/get_all_accounts`` 

<span style="color:#fff; height: 25px; width: 30px; background: red; padding: 7px; margin-right: 5px;">**Headers**</span> ``Content-Type: application/json``

**Description**

```json
[  
   {  
      "Addr":"0xc04ee4131895f1d0c294d508af65d94060aa42bb",
      "Balance":"500000000000000000000",
      "AccountNonce":"1"
   },
   {  
      "Addr":"0x07d899c4ac0c1725c35c5f816e60273b33a964f7",
      "Balance":"100000000000000000000",
      "AccountNonce":"1"
   },
   {  
      "Addr":"0x5bd738164c61fb50eb12e227846cbaef2de965aa",
      "Balance":"400000000000000000000",
      "AccountNonce":"1"
   }
]
```

The above endpoint will respond with account data for all accounts in the postgres database. The table below lists the attributes that will be included in the response from this endpoint.

| Parameter | Description |
| --- | --- |
| `address` | The address of an account to query|

| Attributes    | Description                                     |    Type  |
| :---          |     :---:                                       |    ---:  |
| `Addr`        | Account address                                 | string   |
| `Balance`     | Address balance in wei                          | string   |
| `AccountNonce`| Number of transactions sent from a given address| string   |


**The above endpoints will respond with a json payload for the given request, each of these endpoints are subject to change in the future.**

# Shyft Block Explorer UI

To demonstrate the ability to create your own block explorer, a custom folder was created `./shyftBlockExplorerUI` that contains an example block explorer using react!

To run the Block Explorer UI, ensure that you have the API running as mentioned above. Then run the following command in a terminal:

``cd shyftBlockExplorerUI``

``npm install``

``npm run start``

This will start a development server on ``port 3000`` and spin up an example block explorer that uses the API to query the postgres database.

It should look like the below image.

<img src="./images/explorerUI.png" alt="Block Explorer Example">
                                      
# Command Line Options

Before running any CLI options ensure you run `make geth` in the root directory.

> In the root directory run `./shyft-geth.sh` with any of the following flags:

```shell
--setup              - Setups postgres and the shyft chain db.
--start              - Starts geth.
--reset              - Drops postgres and chain db, and instantiates both.
--js <web3 filename> - Executes web3 calls with a passed file name.
                       If the file name is sendTransactions.js:
                       ./shyft-geth.sh --js sendTransactions
```

For convenience a simple CLI was built using `shyft-geth.sh` as the executable file with some basic commands to get your environment setup.

This will create a new database for geth to use as well as all the necessary tables for the shyft blockexplorer.

# Custom Shyft Constants
### Block Rewards

``./consensus/ethash/consensus.go``

Shyft inflation is different than that of Ethereum, therefore the constants were changed in order to support this.

# Shyft Extended Functionality
### Database Functions

``./core/db.go``

``./shyft_schema``

### Database instanitation

The local database is instantiated where Geth generates and writes the genesis state/block.
``./core/genesis.go``

Specifically, the local database configuration and set up takes place in a custom database file
``./core/db.go``
### Writing Blocks

In our case, we use `SWriteBlock()` for writing all our data. So far, it contains all the data that we need to store to our local block explorer database. It invokes the `SWriteTransaction()` which writes the transactions and updates accounts in the local database. This may change in the future. This function is invoked in:
``./core/blockchain.go``

`SWriteBlock()` and `SWriteTransaction()` exist within:
``./core/shyft_database_util.go``

### Transaction Types Functions

``./core/types/transaction.go``

The existing transaction type in Geth did not allow the evm to call a helper function to retrieve the from address, essentially the sender. Therefore, we extended the functionality of the Transaction type to generate the from address through `*Transaction.From()`.

``./core/shyft_database_util.go``

### Chain Rollbacks

For development and testing purposes only, until a formal messaging system has been incorporated within go-empyrean, an endpoint is available and freely accessible to trigger a chain and postgresql database rollback.

To trigger a chain/pg database rollback the following command should be executed:

```
curl <node ip address>:8081/rollback_blocks/<block hashheader to rollback to>

ie. curl localhost:8081/rollback_blocks/0x6c7db5b09bda0277b480aece97d2efac70838cad4fe6ae45f68410c8cd7cd640
```
