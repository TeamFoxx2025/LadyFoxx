# Rootchain helper command

Top level command for manipulating rootchain server.

## Start rootchain server

This command starts `ethereum/client-go` container which is Geth node started in dev mode.

```bash
$ foxx-chain rootchain server
```

## Fund initialized accounts

This command funds the initialized accounts via `foxx-chain polybft-secrets` command.

```bash
$ foxx-chain rootchain fund --data-dir data-dir- --num 2
```

or

```bash
$ foxx-chain rootchain fund --data-dir data-dir-1
```

## Deploy and initialize contracts

This command deploys and initializes rootchain contracts. Transactions are being sent to given `--json-rpc` endpoint and are signed by private key provided by `--adminKey` flag.

```bash
$ foxx-chain rootchain deploy \
    --genesis <chain_config_file> \
    --deployer-key <hex_encoded_rootchain_deployer_private_key> \
    --stake-manager <stake_manager_address> \
    --stake-token <stake_token_address> \
    --json-rpc <json_rpc_endpoint> \
    --proxy-contracts-admin <address of proxy contracts admin>
```

**Note:** In case `test` flag is provided, it engages test mode, which uses predefined test account private key to send transactions to the rootchain.
