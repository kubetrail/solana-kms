# solana-kms
`solana-kms` is a Google KMS backed Solana token management CLI utility.
The main purpose of the tool is to ensure that the private key is never
written to the disk in plaintext format. This is achieved by persisting
the private key on the disk after encrypting using Google KMS service.

All actions requiring use of private key will go through KMS roundtrip
to decrypt the private key at runtime ensuring that the private key
plaintext data only exists in memory.


## installation
Clone this repo and `cd` to the folder and run `go install`

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

### Use the key
The key can now be used with other Solana CLI tools by piping via STDIN.
```
└─ $ ▶ solana-kms key pubkey 
B9g4B79PHmyCcRQnuAxmzXK1PriVGqmxT7wo4DT7QRUP
```
```
└─ $ ▶ solana-kms key decrypt | solana-keygen pubkey
B9g4B79PHmyCcRQnuAxmzXK1PriVGqmxT7wo4DT7QRUP
```
Similarly, use against an endpoint to airdrop sols and check balance
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
