package auth

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/golang-jwt/jwt/v5"
	"math/big"
	"os"
	"time"
)

func ValidateToken(token string) error {
	verifyToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		fmt.Println(err)
		return err
	}
	if verifyToken.Valid == false {
		return errors.New("invalid token")
	}
	return nil
}

func GenToken(data interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"data": data,
			"exp":  time.Now().Add(time.Hour * 24).Unix(),
		})
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func GenSalt() ([]byte, error) {
	salt := make([]byte, 25)
	_, randErr := rand.Read(salt[:])
	if randErr != nil {
		return nil, randErr
	}
	return salt, nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789!@#$%^&*()_"

func randPassword(size int) (string, error) {
	b := make([]byte, size)
	index, randErr := rand.Int(rand.Reader, big.NewInt(int64(len(letterBytes))))
	if randErr != nil {
		return "", randErr
	}
	for i := range b {
		b[i] = letterBytes[index.Int64()]
		index, randErr = rand.Int(rand.Reader, big.NewInt(int64(len(letterBytes))))
		if randErr != nil {
			return "", randErr
		}
	}
	return string(b), randErr
}

func GenKey() (string, []byte, error) {
	salt, saltError := GenSalt()
	if saltError != nil {
		return "", nil, saltError
	}
	key, keyErr := randPassword(25)
	if keyErr != nil {
		return "", nil, keyErr
	}
	return key, salt, nil
}

func EncryptString(str string, salt []byte) (string, error) {
	stringBytes := []byte(str)
	hasher := sha512.New()
	saltedString := append(stringBytes, salt...)
	hasher.Write(saltedString)
	hashedStringBytes := hasher.Sum(nil)
	HashedSaltedString := hex.EncodeToString(hashedStringBytes)
	return HashedSaltedString, nil
}
