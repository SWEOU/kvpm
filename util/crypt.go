package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/chzyer/readline"
)

// hashTo32Bytes will compute a cryptographically useful hash of the input string.
func hashTo32Bytes(input string) []byte {

	data := sha256.Sum256([]byte(input))
	return data[0:]

}

// Takes two strings, cryptoText and keyString.
// cryptoText is the text to be decrypted and the keyString is the key to use for the decryption.
// The function will output the resulting plain text string with an error variable.
func DecryptString(cryptoText string, keyString string) (plainTextString string, err error) {

	encrypted, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}
	if len(encrypted) < aes.BlockSize {
		return "", fmt.Errorf("cipherText too short. It decodes to %v bytes but the minimum length is 16", len(encrypted))
	}

	decrypted, err := decryptAES(hashTo32Bytes(keyString), encrypted)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

func decryptAES(key, data []byte) ([]byte, error) {
	// split the input up in to the IV seed and then the actual encrypted data.
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(data, data)
	return data, nil
}

// Takes two string, plainText and keyString.
// plainText is the text that needs to be encrypted by keyString.
// The function will output the resulting crypto text and an error variable.
func EncryptString(plainText string, keyString string) (cipherTextString string, err error) {

	key := hashTo32Bytes(keyString)
	encrypted, err := encryptAES(key, []byte(plainText))
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(encrypted), nil
}

func encryptAES(key, data []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// create two 'windows' in to the output slice.
	output := make([]byte, aes.BlockSize+len(data))
	iv := output[:aes.BlockSize]
	encrypted := output[aes.BlockSize:]

	// populate the IV slice with random data.
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)

	// note that encrypted is still a window in to the output slice
	stream.XORKeyStream(encrypted, data)
	return output, nil
}

func ReadPass() (string, error) {
	rl, err := readline.New("")
	if err != nil {
		return "", err
	}
	defer rl.Close()
	setPasswordCfg := rl.GenPasswordConfig()
	setPasswordCfg.SetListener(func(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
		rl.SetPrompt(fmt.Sprintf("Enter password(%v): ", len(line)))
		rl.Refresh()
		return nil, 0, false
	})

	pswd, err := rl.ReadPasswordWithConfig(setPasswordCfg)
	if err != nil {
		return "", nil
	}
	return string(pswd), nil
}
