package flags

const (
	GoogleProjectID              = "google-project-id"              // Google KMS project ID
	KmsLocation                  = "kms-location"                   // KMS location for the key and keyring
	KmsKeyring                   = "kms-keyring"                    // KMS keyring name
	KmsKey                       = "kms-key"                        // KMS key name
	GoogleApplicationCredentials = "google-application-credentials" // Google service account with KMS encrypter/decrypter role
	Config                       = "config"                         // Solana config file
	KeyFile                      = "keyfile"                        // Solana private keypair file
	SeedFile                     = "seedfile"                       // Seedfile associated with private keypair
	PubKey                       = "pubkey"                         // Public key aka Solana address
	Url                          = "url"                            // Solana validator endpoint
)
