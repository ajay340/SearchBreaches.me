package database

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

func GenerateToken(someString string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(someString), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hash)
}

func GenerateHash(someHash string) string {
	hash := md5.Sum([]byte(someHash))
	return hex.EncodeToString(hash[:])
}

func GenerateID() int {
	return int(rand.Int31())
}
