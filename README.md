# KAVA WORKER

This repository contains a worker part dedicated for kava transactions.

## Worker
Stateless worker is responsible for connecting with the chain, getting information, converting it to a common format and sending it back to manager.
Worker can be connected with multiple managers but should always answer only to the one that sent request.

## API
Implementation of bare requests for network.

### Client
Worker's business logic wiring of messages to client's functions.


## Installation
This system can be put together in many different ways.
This readme will describe only the simplest one worker, one manager with embedded scheduler approach.

### Compile
To compile sources you need to have go 1.14.1+ installed.

```bash
    make build
```

### Running
Worker also need some basic config:

```bash
    MANAGERS=0.0.0.0:8085
    TENDERMINT_RPC_ADDR=https://kava-5--rpc--archive.datahub.figment.io
    TENDERMINT_LCD_ADDR=https://kava-5--rpc--archive.datahub.figment.io
    CHAIN_ID=kava-5
```

Where
    - `TENDERMINT_RPC_ADDR` is a http address to node's RPC endpoint
    - `TENDERMINT_LCD_ADDR` is a http address to node's LCD endpoint
    - `MANAGERS` a comma-separated list of manager ip:port addresses that worker will connect to. In this case only one

After running both binaries worker should successfully register itself to the manager.

If you wanna connect with manager running on docker instance add `HOSTNAME=host.docker.internal` (this is for OSX and Windows). For linux add your docker gateway address taken from ifconfig (it probably be the one from interface called docker0).

## Developing Locally

First, you will need to set up a few dependencies:

1. [Install Go](https://golang.org/doc/install)
2. A Kava network node with both RPC and LCD APIs (in this example, we assume it's running at http://127.0.0.1)
3. A running [manager](https://github.com/figment-networks/indexer-manager) instance
4. A running datastore API instance (e.g. [search](https://github.com/figment-networks/indexer-search) - this is configured with `STORE_HTTP_ENDPOINTS`)

Then, run the worker with some environment config:

```
CHAIN_ID=kava-7 \
STORE_HTTP_ENDPOINTS=http://127.0.0.1:8986/input/jsonrpc \
TENDERMINT_RPC_ADDR=http://127.0.0.1:26657 \
TENDERMINT_LCD_ADDR=http://127.0.0.1:1317 \
go run ./cmd/worker-kava
```

Upon success, you should see logs that look like this:

```log
{"level":"info","time":"2021-06-30T13:35:59.634-0400","msg":"kava-worker  (git: ) - built at "}
{"level":"info","time":"2021-06-30T13:35:59.634-0400","msg":"Self-hostname (b2ea3a4e-05bd-48f1-8977-988ef093c1ff) is 0.0.0.0:3000 "}
{"level":"info","time":"2021-06-30T13:35:59.634-0400","msg":"Connecting to managers (127.0.0.1:8085)"}
{"level":"info","time":"2021-06-30T13:35:59.634-0400","msg":"[New client]","url":"http://127.0.0.1:26657"}
{"level":"info","time":"2021-06-30T13:35:59.635-0400","msg":"[New client]","url":"http://127.0.0.1:1317"}
{"level":"info","time":"2021-06-30T13:35:59.635-0400","msg":"[HTTP] Listening on 0.0.0.0:8087"}
{"level":"info","time":"2021-06-30T13:35:59.635-0400","msg":"[GRPC] Listening on 0.0.0.0:3000"}
```

Once the worker connects to a running [manager](https://github.com/figment-networks/indexer-manager), which runs by default at `127.0.0.1:8085`, you should see a stream registered in the logs:

```log
{"level":"debug","time":"2021-06-30T13:55:46.168-0400","msg":"[KAVA-CLIENT] Register Stream","streamID":"55bd3282-764f-4314-ae6a-a0936c18e04d"}
{"level":"debug","time":"2021-06-30T13:55:46.168-0400","msg":"[GRPC] Send started "}
```

## Debug with VSCode

The `.vscode` directory contains a launch config to debug the worker.


## Transaction Types
List of currently supporter transaction types in kava-worker are (listed by modules):
- auction:
   `place_bid`
- bank:
    `multisend` , `send`
- bep3:
    `create_atomic_swap`, `claim_atomic_swap`, `refund_atomic_swap`
- cdp:
    `create_cdp`, `deposit_cdp`, `withdraw_cdp`, `draw_cdp`, `repay_cdp`, `liquidate`
- committee:
    `commmittee_submit_proposal`, `committee_vote`
- crisis:
    `verify_invariant`
- distribution:
    `withdraw_validator_commission` , `set_withdraw_address` , `withdraw_delegator_reward` , `fund_community_pool`
- evidence:
    `submit_evidence`
- gov:
    `deposit` , `vote` , `submit_proposal`
- hard:
    `hard_deposit`, `hard_withdraw`,`hard_repay`,`hard_borrow`,`hard_liquidate`,
- incentive:
    `claim_hard_reward`,`claim_usdx_minting_reward`,
- issuance:
    `issue_tokens`, `redeem_tokens`, `block_address`, `unblock_address`, `change_pause_status`
- pricefeed:
    `post_price`
- slashing:
    `unjail`
- staking:
    `begin_unbonding` , `edit_validator` , `create_validator` , `delegate` , `begin_redelegate`
- internal:
    `error`
