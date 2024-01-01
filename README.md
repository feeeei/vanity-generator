English | [中文](./README_ZH-CN.md)

# vanity-generator

`vanity-generator` is a blockchain wallet vanity number generator implemented in Go. It is highly efficient, capable of generating over 200,000 vanity addresses per second on a single core and millions per second on multiple cores.

Wallet generation is currently supported for three networks: `Ethereum`、`Tron`、`Polkadot`.

![preview](images/preview.jpg)

### Features:

- Purely local generation, no internet connection required, can be used as a cold wallet
- 100% based on the official project code of [go-ethereum](https://github.com/ethereum/go-ethereum), providing better security
- Implemented in Go language, generating efficiency is <strong>several tens of times</strong> higher than the [JS version](https://vanity-eth.tk/)
- Supports specifying prefix and suffix, supports specifying both at the same time
- Provides estimated time, default provides estimated time of 50%, 70%, 90% probability
- Native high concurrency support, default uses all CPU cores

### Native Usage:
Command: `./vanity {eth/tron/dot} --prefix=xxxx --suffix=xxx --concurrency=1`

It is recommended to use with commands such as `screen` or `nohup`.

### Usage with Docker:
Create a configuration file at /etc/vanity/config.json with the following format:
```json
{
  "wallet": "eth",
  "prefix": "xxxx",
  "suffix": "xxxx",
  "concurrency": 1
}
```
Command: 
```shell
docker pull feeeei/vanity-generator
docker run --name=vanity -v /etc/vanity:/etc/vanity feeeei/vanity-generator
```

### Parameter Description:

- `--prefix`: Specify prefix, ETH needs to start with `0x`, Tron needs to start with `T`.
- `--suffix`: Specify suffix, ETH needs to satisfy [0-9A-Fa-f], Tron needs to satisfy Base58.
- `--concurrency`: Number of concurrent processes, if not specified, it is equal to the number of CPU cores by default.