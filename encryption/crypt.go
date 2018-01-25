package encryption

import (
    "io/ioutil"
    "io"
    "errors"

    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
)

type Encryption struct {
    Key []byte
}

var instance *Encryption = nil

func CreateInstance() *Encryption {
    instance = &Encryption{
    }

    return instance
}

func GetInstance() *Encryption {
    return instance
}

// https://gist.github.com/manishtpatel/8222606
func Encrypt(plaintext []byte) ([]byte, error){
    key := GetInstance().Key

    block, err := aes.NewCipher(key)
    if err != nil {
        return []byte(""), err
    }

    // The IV needs to be unique, but not secure. Therefore it's common to
    // include it at the beginning of the ciphertext.
    ciphertext := make([]byte, aes.BlockSize+len(plaintext))
    iv := ciphertext[:aes.BlockSize]
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return []byte(""), err
    }

    stream := cipher.NewCFBEncrypter(block, iv)
    stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

    return ciphertext, nil
}

func Decrypt(ciphertext []byte) ([]byte, error){
    key := GetInstance().Key

    block, err := aes.NewCipher(key)
    if err != nil {
        return []byte(""), err
    }

    // The IV needs to be unique, but not secure. Therefore it's common to
    // include it at the beginning of the ciphertext.
    if len(ciphertext) < aes.BlockSize {
        return []byte(""), errors.New("ciphertext too short")
    }
    iv := ciphertext[:aes.BlockSize]
    ciphertext = ciphertext[aes.BlockSize:]

    stream := cipher.NewCFBDecrypter(block, iv)

    // XORKeyStream can work in-place if the two arguments are the same.
    stream.XORKeyStream(ciphertext, ciphertext)

    return ciphertext, nil
}

func ExportKeyToFile(key []byte, filename string){
    ioutil.WriteFile(filename, key, 0644)
}

func ImportKeyFromFile(filename string) ([]byte, error){
    k, err := ioutil.ReadFile(filename)
    if err == nil{
        GetInstance().Key = k
    }
    return k, err
}

func GenerateRandomKey() ([]byte, error){
    key := make([]byte, 32)

    _, err := rand.Read(key)
    return key, err
}

