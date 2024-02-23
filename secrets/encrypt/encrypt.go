package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
)

func Encrypt(key, plaintext []byte) ([]byte, error) {
	block, err := newCipherBlock(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

func Decrypt(key, ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("cipher text too short")
	}

	block, err := newCipherBlock(key)
	if err != nil {
		return nil, err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}

func newCipherBlock(key []byte) (cipher.Block, error) {
	hasher := md5.New()
	fmt.Fprint(hasher, key)
	cipherKey := hasher.Sum(nil)
	return aes.NewCipher(cipherKey)
}

func EncryptWriter(key []byte, w io.Writer) (io.Writer, error) {
	block, err := newCipherBlock(key)
	if err != nil {
		return nil, err
	}

	var iv [aes.BlockSize]byte
	stream := cipher.NewCFBEncrypter(block, iv[:])

	writer := &cipher.StreamWriter{S: stream, W: w}

	return writer, nil
}

func DecryptReader(key []byte, r io.Reader) (io.Reader, error) {
	block, err := newCipherBlock(key)
	if err != nil {
		return nil, err
	}

	var iv [aes.BlockSize]byte
	stream := cipher.NewCFBDecrypter(block, iv[:])

	reader := &cipher.StreamReader{S: stream, R: r}

	return reader, nil
}
