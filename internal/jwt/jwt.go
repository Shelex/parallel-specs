package jwt

import (
	"crypto/rsa"
	"io/ioutil"
	"log"
	"time"

	"github.com/Shelex/parallel-specs/internal/entities"
	"github.com/Shelex/parallel-specs/internal/users"
	"github.com/dgrijalva/jwt-go"
)

const (
	privKeyPath = "keys/app.rsa"     // openssl genrsa -out keys/app.rsa keysize
	pubKeyPath  = "keys/app.rsa.pub" // openssl rsa -in keys/app.rsa -pubout > keys/app.rsa.pub
)

var (
	verifyKey *rsa.PublicKey
	SignKey   *rsa.PrivateKey
)

func init() {
	signBytes, err := ioutil.ReadFile(privKeyPath)
	fatal(err)

	SignKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	fatal(err)

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	fatal(err)

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	fatal(err)
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// data we save in each token
type Claims struct {
	email  string //nolint
	id     string //nolint
	exp    uint64 //nolint
	entity string //nolint

	jwt.StandardClaims
}

// GenerateToken generates a jwt token and assign an email to it's claims and return it
func GenerateToken(user users.User) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	/* Create a map to store our claims */
	claims := token.Claims.(jwt.MapClaims)
	/* Set token claims */
	claims["email"] = user.Email
	claims["id"] = user.ID
	claims["entity"] = "user"
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, err := token.SignedString(SignKey)
	if err != nil {
		log.Fatal("Error in Generating key")
		return "", err
	}
	return tokenString, nil
}

// GenerateApiKey generates a jwt token and assign an user with customized expiry
func GenerateApiKey(user users.User, apiKey entities.ApiKey) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	/* Create a map to store our claims */
	claims := token.Claims.(jwt.MapClaims)
	/* Set token claims */
	claims["email"] = user.Email
	claims["id"] = user.ID
	claims["entity"] = apiKey.ID
	claims["exp"] = apiKey.ExpireAt
	tokenString, err := token.SignedString(SignKey)
	if err != nil {
		log.Fatal("Error in Generating key")
		return "", err
	}
	return tokenString, nil
}
