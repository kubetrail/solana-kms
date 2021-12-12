# solana-kms
`solana-kms` is a Google KMS backed Solana token management CLI utility.
The main purpose of the tool is to ensure that the private key is never
written to the disk in plaintext format. This is achieved by persisting
the private key on the disk after encrypting using Google KMS service.

All actions requiring use of private key will go through KMS roundtrip
to decrypt the private key at runtime ensuring that the private key
plaintext data only exists in memory.

## Installation
Clone this repo and `cd` to the folder and run `go install`. This assumes you
have Go (Golang) toolchain installed. Also install Solana CLI tools.

* https://go.dev/dl/
* https://docs.solana.com/cli/install-solana-cli-tools

## Setup
### Solana Config
Setup config file used by `Solana` CLI tools to look something as follows
```
└─ $ ▶ solana config get
Config File: /home/username/.config/solana/cli/config.yml
RPC URL: http://localhost:8899 
WebSocket URL: ws://localhost:8900/ (computed)
Keypair Path: stdin:/home/username/.config/solana/id 
Commitment: confirmed 
```
Note that the `Keypair Path` has `stdin:` prefix, which triggers the behavior
by Solana CLI tools to accept keypair from STDIN. The actual path is ignored
by the Solana CLI tools, however, this CLI makes use of that path to store
the KMS encrypted data.

Make sure that the file `/home/username/.config/solana/id` does not exist when
creating a new key.

### Setup KMS
Please follow Google KMS instructions to create a KMS keyring and a key. In addition,
you will also need a service account with KMS encrypter/decrypter role. Setup following
environment variables:
```bash
export GOOGLE_APPLICATION_CREDENTIALS="${HOME}/.config/service-accounts/service-account.json"
export LOCATION=<kms location>
export KEYRING=<keyring name>
export KEY=<key name>
export PROJECT=<project id>
```

## Usage
### Create a new key
Run `solana-kms key new` to generate a new keypair. It will write the encrypted keydata
to the filepath set previously (`/home/username/.config/solana/id`) and also write
encrypted seed for it.

```
└─ $ ▶ solana-kms key new
```

### Use the key
The key can now be used with other Solana CLI tools by piping via STDIN. You can verify
that is working by fetching public key address from `solana-kms` directly or via
`solana-keygen` CLI as shown below:
```
└─ $ ▶ solana-kms key pubkey 
B9g4B79PHmyCcRQnuAxmzXK1PriVGqmxT7wo4DT7QRUP
```
```
└─ $ ▶ solana-kms key decrypt | solana-keygen pubkey
B9g4B79PHmyCcRQnuAxmzXK1PriVGqmxT7wo4DT7QRUP
```
Similarly, this setup can be used against all Solana CLI tools that support
keypair input via STDIN
```
└─ $ ▶ solana-kms key decrypt | solana airdrop 100 
Requesting airdrop of 100 SOL

Signature: 5Z8s3BUStxym3NyXT7zE2fxpgVi6VnL66F78ZyMkr6DdQ1v65ZgsW74xT2oJrDuv1kvGwfp8tjYiSvNdEDfSMgxm

100 SOL
```

```
└─ $ ▶ solana-kms key decrypt | solana balance
100 SOL
```

## Key Rotation
It is possible to regenerate the keypair from the seed. The newly created ecrypted
file will differ from the original, however, they both would map to the same
plaintext bytes.
```
└─ $ ▶ solana-kms key new --keyfile=/path/to/new/id --seedfile=/path/to/old/id.seed
```
This will generate two new files `/path/to/new/id` and `/path/to/new/id.seed`. At this
point the previous keypair and seed files can be deleted and any associated KMS
key version that was used for encryption of those old keypair files,
but is no longer in use, can also be disabled. In other words, we are not rotating
the Solana keys, the roation implies to the encrypted content using different
versions of the KMS keys.

## Security Concerns
* https://unix.stackexchange.com/questions/156859/is-the-data-transiting-through-a-pipe-confidential

