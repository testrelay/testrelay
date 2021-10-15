package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

// TokenVerifier verifies JWT tokens with provided config.
type TokenVerifier interface {
	VerifyToken(token *jwt.Token) (interface{}, error)
	Parse(token string) error
}

// Auth0Verifier is satisfies teh token verifier using auth0 config
type Auth0Verifier struct {
	Aud     string
	Iss     string
	JwksURL string
}

// Parse parses the token string to a jwt token and returns error if failed.
func (a Auth0Verifier) Parse(token string) error {
	_, err := jwt.Parse(token, a.VerifyToken)

	return err
}

// VerifyToken takes a token a verifies it against a auth0 configuration
func (a Auth0Verifier) VerifyToken(token *jwt.Token) (interface{}, error) {
	// Verify 'aud' claim
	checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(a.Aud, false)
	if !checkAud {
		return nil, errors.New("Invalid audience.")
	}
	// Verify 'iss' claim
	checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(a.Iss, false)
	if !checkIss {
		return nil, errors.New("Invalid issuer.")
	}

	cert, err := a.getPemCert(token)
	if err != nil {
		panic(err.Error())
	}

	t, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
	return t, err
}

func (a Auth0Verifier) getPemCert(token *jwt.Token) (string, error) {
	var cert string
	resp, err := http.Get(a.JwksURL)

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("Unable to find appropriate key.")
		return cert, err
	}

	return cert, nil
}


// CertsAPIEndpoint is endpoint of getting Public Key.
var CertsAPIEndpoint = "https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com"

// GetCertificate is useful for testing.
var GetCertificate = getCertificate

func getCertificates() (map[string]string, error) {
	var certs map[string]string
	res, err := http.Get(CertsAPIEndpoint)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &certs)

	return certs, err
}

// GetCertificate returns certificate.
func getCertificate(kid string) ([]byte, error) {
	certs, err := getCertificates()
	if err != nil {
		return nil, err
	}
	certString := certs[kid]
	cert := []byte(certString)
	return cert,  err
}

// GetCertificateFromToken returns cert from token.
func GetCertificateFromToken(token *jwt.Token) ([]byte, error) {
	// Get kid
	kid, ok := token.Header["kid"]
	if !ok {
		return []byte{}, errors.New("kid not found")
	}
	kidString, ok := kid.(string)
	if !ok {
		return []byte{}, errors.New("kid cast error to string")
	}
	return GetCertificate(kidString)
}


func readPublicKey(cert []byte) (*rsa.PublicKey, error) {
	publicKeyBlock, _ := pem.Decode(cert)
	if publicKeyBlock == nil {
		return nil, errors.New("invalid public key data")
	}
	if publicKeyBlock.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("invalid public key type: %s", publicKeyBlock.Type)
	}
	c, err := x509.ParseCertificate(publicKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}
	publicKey, ok := c.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not RSA public key")
	}
	return publicKey, nil
}

type FirebaseVerifier struct {
	ProjectID string
}

func (f *FirebaseVerifier) VerifyToken(t *jwt.Token) (interface{}, error) {
	cert, err := GetCertificateFromToken(t)
	if err != nil {
		return "", err
	}
	publicKey, err := readPublicKey(cert)
	if err != nil {
		return "", err
	}

	return publicKey, nil
}

func (f *FirebaseVerifier) Parse(token string) error {
	parsed, err := jwt.Parse(token, f.VerifyToken)
	if err != nil {
		return fmt.Errorf("could not parse token %w", err)
	}

	if !parsed.Valid {
		return fmt.Errorf("invalid token")
	}
	// Verify header.
	if parsed.Header["alg"] != "RS256" {
		return fmt.Errorf("invalid algorithm %s", parsed.Header["alg"])
	}

	return f.verifyPayload(parsed)
}

// Verify the token payload.
func (f *FirebaseVerifier) verifyPayload(t *jwt.Token) error {
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("could not map jwt claims")
	}
	// Verify User
	claimsAud, ok := claims["aud"].(string)
	if !ok || claimsAud != f.ProjectID {
		return fmt.Errorf("incorect jwt audience")
	}

	iss := "https://securetoken.google.com/" + f.ProjectID
	claimsIss, ok := claims["iss"].(string)
	if !ok || claimsIss != iss {
		return fmt.Errorf("incorect issuer")
	}

	return nil
}
