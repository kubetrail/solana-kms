package flags

const (
	Project                = "project"                 // Google KMS project ID
	Location               = "location"                // KMS location for the key and keyring
	Keyring                = "keyring"                 // KMS keyring name
	Key                    = "key"                     // KMS key name
	ApplicationCredentials = "application-credentials" // Google service account with KMS encrypter/decrypter role
	Config                 = "config"                  // Solana config file
	KeyFile                = "keyfile"                 // Solana private keypair file
	SeedFile               = "seedfile"                // Seedfile associated with private keypair
)
