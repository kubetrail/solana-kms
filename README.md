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

