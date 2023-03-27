package token

// Acknowledgements to inanzzz for
// http://www.inanzzz.com/index.php/post/kdl9/creating-and-validating-a-jwt-rsa-token-in-golang

import (
	"crypto"
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	privateKey *rsa.PrivateKey
	publicKey  crypto.PublicKey
}

func NewJWT(privateKeyPath string) (*JWT, error) {

	// Load the RSA 4096-bit pem private key file
	keyBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load RSE pem private file: %w", err)
	}

	// Parse the pem bytes into an RSA key
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSE private key: %w", err)
	}

	// Return our JWT instance
	return &JWT{
		privateKey: privateKey,
		publicKey:  privateKey.Public(),
	}, nil
}

func (j *JWT) Create(ttl time.Duration, content interface{}) (string, error) {
	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["dat"] = content             // Our custom data.
	claims["exp"] = now.Add(ttl).Unix() // The expiration time after which the token must be disregarded.
	claims["iat"] = now.Unix()          // The time at which the token was issued.
	claims["nbf"] = now.Unix()          // The time before which the token must be disregarded.

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(j.privateKey)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

func (j *JWT) Validate(token string) (interface{}, error) {

	// Parse the JWT and validate its signature
	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {

		// Confirm that the JWT header signals that the JWT was signed with the RSA algorithm
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		// All is well with the header, return our public key
		return j.publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return nil, fmt.Errorf("validate: invalid")
	}

	return claims["dat"], nil
}
