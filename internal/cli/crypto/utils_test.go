package crypto

import (
	"iwakho/gopherkeep/internal/model"
	"testing"
)

func TestCryptoManager(t *testing.T) {
	p := model.Pair{
		Login:    "alice",
		Password: "12345",
	}
	cm := newCryptoManager(p)
	sample := "Шифрование очень длинного сообщения, которое будет зашифровано, а затем расшифровано, чтобы проверить, что шифрование работает правильно."

	encrypted, err := cm.Encrypt(sample)
	if err != nil {
		t.Errorf("CryptoManager.Encrypt() error = %v", err)
		return
	}
	decrypted, err := cm.Decrypt(encrypted)
	if err != nil {
		t.Errorf("CryptoManager.Decrypt() error = %v", err)
		return
	}
	if decrypted != sample {
		t.Errorf("CryptoManager.Decrypt() = %v, want %v", decrypted, sample)
	}
	t.Logf("decrypted: %s", decrypted)
}
