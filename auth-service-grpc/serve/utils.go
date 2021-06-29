package serve

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
)

func PEMStringToRSAPublicKey(pemStr string) (*rsa.PublicKey, error) {
	p, _ := pem.Decode([]byte(pemStr))

	pub, err := x509.ParsePKIXPublicKey(p.Bytes)

	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break
	}

	return nil, fmt.Errorf("public key type is not RSA")
}

func StructToJSON(obj interface{}) (string, error) {
	jsonBytes, err := json.Marshal(obj)

	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}
