package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	util "github.com/urento/shoppinglist/pkg"
)

var jwtSecret []byte

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

func GenerateToken(username, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	passwordHash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Fatal(err)
	}

	claims := Claims{
		Encrypt(username, getKeyString()),
		passwordHash,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "shoppinglist",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

func getKeyString() (keyString string) {
	var err error
	if util.PROD {
		err = godotenv.Load()
	} else {
		err = godotenv.Load("../.env")
	}
	if err != nil {
		log.Fatal(err.Error())
	}
	return os.Getenv("ENCRYPTION_KEYSTRING")
}

func Encrypt(stringToEncrypt string, keyString string) (encryptedString string) {
	key, _ := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}

func Decrypt(encryptedString, keyString string) (decryptedString string) {
	key, _ := hex.DecodeString(keyString)
	enc, _ := hex.DecodeString(encryptedString)

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err.Error())
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		log.Fatal(err.Error())
	}

	nonceSize := aesGCM.NonceSize()

	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	return fmt.Sprintf("%s", plaintext)
}
