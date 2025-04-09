package tests

import (
	"encoding/hex"
	"testing"

	"github.com/gogf/gf/v2/crypto/gdes"
)

func TestDesEncrypt(t *testing.T) {
	const (
		// Maintain device PassId
		ENCRYPTED_DATA = "0cf23c53f0f22469ddb69c262b53c604b482e8f1b736528dc316e9064940cc51c7ed7bbb4af60a2a6981c5ad80600c9762667d4e74f2293fb7ac43fd12a43a7ca9dcfdde32727a0b89f010e04a4625aaef3cf9ad5adfed0e908b0f2e96124652d0df9b3dc6dc22d4ffeebdfb8cc8f694"

		DES_KEY     = "0123456789abcd0123456789"
		DES_IV      = "12345679"
		DES_PADDING = 1
	)

	data, err := hex.DecodeString(ENCRYPTED_DATA)
	if err != nil {
		t.Fatalf("Failed to decode string: %s\n", err.Error())
	}

	data, err = gdes.DecryptCBCTriple(data, []byte(DES_KEY), []byte(DES_IV), DES_PADDING)
	if err != nil {
		t.Fatalf("Failed to decrypt data: %s\n", err.Error())
	}

	t.Logf("Decrypted data: %s\n", data)
}
