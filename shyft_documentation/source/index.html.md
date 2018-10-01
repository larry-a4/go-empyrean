---
title: API Reference

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

Before running any CLI options ensure you run `make geth` in the root directory.

``.shyft-geth.sh --setup`` This sets up postgres and the shyft chain db

``./shyft-geth.sh --start`` This starts GETH

At this point you should see GETH running in the terminal and if you opened your postgres instance you should see data being populated into the tables.

To stop Geth, `crtl+C` in the terminal window, if you proceed with the start script mentioned above the Shyft chain will begin from the last block height, if you wish to start the chain fresh from genesis follow the below steps:

``./shyft-geth.sh --reset`` This drops postgres and chaindb data

``./shyft-geth.sh --start`` Starts GETH

To see transactions being submitted on the network see the sendTransactions command in the CLI section of this readme.
### Docker Images

Docker Images are available for ShyftGeth and the Postgresql Database which can be used for development and testing. To launch these containers you will need to have docker-compose installed on your computer. Installation instructions for docker-compose are available [here](https://docs.docker.com/install/).

<aside class="warning">
There is currently an issue with respect to starting and stopping the docker containers. Currently docker stop does not allow the shyft geth instance to undergo an orderly shutdown and corrupts the chaindata, resulting in a complete rewind of the chain and the chain being rebuilt from block 0. We are currently researching a solution whereby the docker stop command will be captured by the shyft geth container allowing for an orderly shutdown and not corrupting the chain data.
</aside>

To launch ShyftGeth, PG, the ShyftBlock Explorer Api and UI - issue the following commands from the root of the project directory:

`./shyft-geth.sh --setup # clears persisted directories prior to docker build`

`docker-compose up`

If you would like to reinitialize/rebuild the docker images you can just add the build flag to the docker command:


`./shyft-geth.sh --setup # clears persisted directories prior to docker build`

`docker-compose up --build`

To stop/pause mining - enter:

`docker-compose stop`

And then just issue `docker-compose up` to continue mining.
### Docker Postgresql - DB Connection
From your local machine you can view the database by connecting to the database in the container at 
``127.0.0.1:8001``

Use User: 'postgres' Password: 'docker' Database: 'shyftdb'
### Docker Block Explorer Api 
To access the shyftBlockExplorer open a browser and visit 

``http://localhost:3000``

To rebuild any one of the services - issue the following commands:

``
docker-compose up -d --no-deps --build <docker compose file service name> 
``

ie. for shyftBlockExplorerApi:

``docker-compose up -d --no-deps --build shyft_block_api``

The Postgresql Database Container will persist the database data to the directory ``./pg-data`` _. So if you do want to reinitialize the database you should delete this directory as well as the blockchain data directories ``(./shyftData ./privatenet)`` prior to launching the docker containers. There is a shell script available to delete these folders to run it execute the following command:

``./shyft-cli/resetShyftGeth.sh``

Blockchain data is persisted to ``./ethash/.ethash and ./shyftData__``. If you would like to reset the test blockchain you will need to delete the ``__./ethash ./shyftData & ./privatenet__`` directories.

The docker container for the ShyftBlockExplorerApi utilizes govendor to minimize its image size. __If you would like the docker image for this container to reflect any uncommitted changes which may have occurred in the go-empyrean repository, ie. changes with respect to go-empyrean core (ie. cryptographic functions and database). Prior to launching the docker containers you should rebuild the vendor directory for the shyftBlockExplorerApi - by executing the following steps:__

Remove existing shyftBlockExplorerApi vendor.json and vendored components:

``rm -rf shyftBlockExplorerApi/vendor``

reinitialize vendor.json

``cd shyftBlockExplorerApi && govendor init``

rebuild vendor.json using latest uncommitted changes

``govendor add +external``

Due to a bug in govendor and it not being able to pull in some dependencies that are c-header files 
you should execute the following commands - see these issues - which whilst closed
appears to have not been fixed: https://github.com/kardianos/govendor/issues/124 && https://github.com/kardianos/govendor/issues/61

``govendor remove github.com/ShyftNetwork/go-empyrean/crypto/secp256k1/^``

``govendor fetch github.com/ShyftNetwork/go-empyrean/crypto/secp256k1/^``

NB: The Shyft Geth docker image size is 1+ GB so make sure you have adequate space on your disk drive/

# Shyft BlockExplorer API

In order to store the block explorer database, a custom folder was created `./shyft_schema` that contains all the necessary functions to read and write to the explorer database.

The main functions exist in `./core/shyft_database_util.go` and `./core/shyft_get_utils.go`

To run the block explorer rest api that queries the postgres instance and returns a json body, open a new terminal window, navigate to the root directory of the project and run the following command:

``go run blockExplorerApi/*.go``

This will start a go server on port 8080 and allow you to either run the pre-existing block explorer or query the api endpoints. Its important to note, that if you have nothing in your postgres database the API will return nothing.

Below is an API map containing the different endpoints you can query. If you are running locally and example request would be like so:

This would return the block data for block number 10, like so: 

```
GET http://localhost:8080/api/get_block/10
Responses 200 
Headers:
"Content-Type": "application/json"
Body:
```
```json
{
    "Hash":"0xb6f0906a276d992e9dc82f82e3be5487251ff6e7b8ff6b0e5e1603092f534799",
    "Coinbase":"0x43EC6d0942f7fAeF069F7F63D0384a27f529B062",
    "Number":"10",
    "GasUsed":"189000",
    "GasLimit":"26863872",
    "TxCount":"9",
    "UncleCount":"0",
    "Age":"2018-05-10T16:26:02Z"
}

```
| GET ENDPOINTS  | Description | Type                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
| ----------- | -------------- | ------------
|  `/api/get_block/{blockNumber} `  | Returns block data by block height/number | Blocks
|  `/api/get_all_blocks`            | Returns block data for all blocks | Blocks
|  `/api/get_recent_block`         | Returns block data for the most recent block mined | Blocks
|  `/api/get_blocks_mined/{coinbase}` | Returns block data by miner address | Blocks
|  `/api/get_transaction/{txHash}`                          | Returns tx data by transaction hash | Transactions
|  `/api/get_all_transactions`                              | Returns tx data for all transactions | Transactions
|  `/api/get_all_transactions_from_block/{blockNumber}`     | Returns tx data by block height/number | Transactions
|  `/api/get_internal_transactions/{address}`               | Returns internal tx data by address | Transactions
|  `/api/get_internal_transactions_hash/{transactions_hash}`|  Returns internal tx data by transaction hash | Transactions
|  `/api/get_account/{address}`                             | Returns account data by address | Accounts
|  `/api/get_account_txs/{address}`                         | Returns tx data by address | Accounts
|  `/api/get_all_accounts`                                  | Returns account data from all accounts | Accounts

The above endpoints will respond with a json payload for the given request, each of these endpoints are subject to change in the future.

# Shyft Block Explorer UI Example

To demonstrate the ability to create your own block explorer, a custom folder was created `./shyftBlockExplorerUI` that contains an example block explorer using react!

To run the Block Explorer UI, ensure that you have the API running as mentioned above. Then run the following command in

``./shyftBlockExplorerUI``

``npm run start``

This will start a development server on port 3000 and spin up an example block explorer that uses the API to query the postgres database.
                                      
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

