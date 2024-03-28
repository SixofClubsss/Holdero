package holdero

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"strconv"

	"github.com/dReam-dApps/dReams/rpc"
)

var handKey string
var handKeyLock bool

// Gets local cards with local key
func findCard(hash string) int {
	for i := 1; i < 53; i++ {
		finder := strconv.Itoa(i)
		add := handKey + finder + round.seed
		sha := sha256.Sum256([]byte(add))
		str := hex.EncodeToString(sha[:])

		if str == hash {
			return i
		}

	}
	return 0
}

// Generate a new Holdero key
func generateKey() string {
	random, _ := rand.Prime(rand.Reader, 128)
	shasum := sha256.Sum256([]byte(random.String()))
	str := hex.EncodeToString(shasum[:])
	handKeyLock = true
	rpc.PrintLog("[Holdero] Round Key: %s", str)

	return str
}

// Create pass hash
func createHash(key string) string {
	sha := sha256.Sum256([]byte(key))
	md5 := md5.New()
	md5.Write([]byte(hex.EncodeToString(sha[:])))
	return hex.EncodeToString(md5.Sum(nil))
}

// Encrypt plaintext data with pass
func Encrypt(data []byte, tag, pass, add string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(pass)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		logger.Printf("[%s] Encrypt %s\n", tag, err)
		return nil
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		logger.Printf("[%s] Encrypt %s\n", tag, err)
		return nil
	}

	extra := []byte(add)

	return gcm.Seal(nonce, nonce, data, extra)
}

// Decrypt ciphertext with pass
func Decrypt(data []byte, tag, pass, add string) []byte {
	key := []byte(createHash(pass))
	block, err := aes.NewCipher(key)
	if err != nil {
		logger.Printf("[%s] Decrypt %s\n", tag, err)
		return nil
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		logger.Printf("[%s] Decrypt %s\n", tag, err)
		return nil
	}

	extra := []byte(add)

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, extra)
	if err != nil {
		logger.Printf("[%s] Decrypt %s\n", tag, err)
		return nil
	}

	return plaintext
}

// Write encrypted file
func EncryptFile(data []byte, tag, filename, pass, add string) {
	if data != nil {
		if file, err := os.Create(filename); err == nil {
			defer file.Close()
			file.Write(Encrypt(data, tag, pass, add))
		}
	}
}

// Decrypt a file
func DecryptFile(tag, filename, pass, add string) []byte {
	if data, err := os.ReadFile(filename); err == nil {
		return Decrypt(data, tag, pass, add)
	}
	return nil
}

// Sets a nil handKey string
func setHandKey() {
	if handKey == "" {
		shasum := sha256.Sum256([]byte("nil"))
		handKey = hex.EncodeToString(shasum[:])
	} else {
		handKeyLock = true
	}
}
