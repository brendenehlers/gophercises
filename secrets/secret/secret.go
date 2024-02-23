package secret

import (
	"fmt"
	"gophercises/secrets/encrypt"
	"io"
	"os"
	"strings"
)

func NewVault(key, filepath string) Vault {
	return Vault{
		key:      []byte(key),
		filepath: filepath,
	}
}

type Vault struct {
	key      []byte
	filepath string
}

type VaultSecrets map[string]string

func (v Vault) Get(key string) (string, bool, error) {
	secrets, err := v.readSecrets()
	if err != nil {
		return "", false, err
	}

	secret, ok := secrets[key]
	if !ok {
		return "", false, nil
	}

	return secret, true, nil
}

func (v Vault) Set(key, val string) error {
	secrets, err := v.readSecrets()
	if err != nil {
		return err
	}

	secrets[key] = val

	err = v.writeSecrets(secrets)
	if err != nil {
		return err
	}

	return nil
}

func (v Vault) readSecrets() (VaultSecrets, error) {
	file, err := openFile(v.filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cr, err := encrypt.DecryptReader(v.key, file)
	if err != nil {
		return nil, err
	}

	secrets, err := decodeSecrets(cr)
	if err != nil {
		return nil, err
	}

	return secrets, nil
}

func decodeSecrets(r io.Reader) (VaultSecrets, error) {
	plaintext, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	secrets := make(VaultSecrets)
	if len(plaintext) == 0 {
		return secrets, nil
	}

	records := strings.Split(string(plaintext), "\n")
	for _, record := range records {
		kv := strings.Split(record, "=")
		fmt.Println(kv)
		// if len(kv) != 2 {
		// 	return nil, fmt.Errorf("invalid format when decoding secrets")
		// }
		secrets[kv[0]] = kv[1]
	}

	return secrets, nil
}

func (v Vault) writeSecrets(secrets VaultSecrets) error {
	file, err := openFile(v.filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	encrypter, err := encrypt.EncryptWriter(v.key, file)
	if err != nil {
		return err
	}

	contents := encodeSecrets(secrets)
	encrypter.Write([]byte(contents))

	return nil
}

func encodeSecrets(secrets VaultSecrets) string {
	var records []string
	for k, v := range secrets {
		records = append(records, fmt.Sprintf("%s=%s", k, v))
	}
	contents := strings.Join(records, "\n")

	return contents
}

func openFile(filepath string) (*os.File, error) {
	return os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0666)
}
