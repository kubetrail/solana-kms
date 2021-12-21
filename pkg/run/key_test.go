package run

import (
	"testing"

	"github.com/portto/solana-go-sdk/types"
)

// TestKeyDecrypt tests key creation from base58 string.
// {"account_id":"4387925efbc0659ac749f9d6e6c42f1db174e03f9d9a07325db84060937fae53","public_key":"ed25519:5YcDXyNdRuKVZ8aoPjy2uCGyHPJUKfxmWHXM4JneuCsc","private_key":"ed25519:39fVpgen8BDGfryixELCVEV51D2CNUJG6DREeRAQ7Qn564rzarkBMeQb6HxdLyZw1xKqhNEqwNMAxuFr24xiX6yG"}
// Above key was generated using near protocol cli
func TestKeyDecrypt(t *testing.T) {
	privateKey := "39fVpgen8BDGfryixELCVEV51D2CNUJG6DREeRAQ7Qn564rzarkBMeQb6HxdLyZw1xKqhNEqwNMAxuFr24xiX6yG"
	publicKey := "5YcDXyNdRuKVZ8aoPjy2uCGyHPJUKfxmWHXM4JneuCsc"
	account, err := types.AccountFromBase58(privateKey)
	if err != nil {
		t.Fatal(err)
	}

	if account.PublicKey.ToBase58() != publicKey {
		t.Fatal("public keys do not match")
	}
}
