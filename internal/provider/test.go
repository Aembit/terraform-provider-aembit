package provider

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var maxRand = big.NewInt(10000000)

const (
	tagsCount    = "tags.%"
	tagsAllCount = "tags_all.%"
	tagsColor    = "tags.color"
	tagsDay      = "tags.day"
	tagsAllName  = "tags_all.Name"
	tagsAllOwner = "tags_all.Owner"
)

func randomizeFileConfigs(newConfig, modifyConfig, startValue string) (string, string, string) {
	randID, _ := rand.Int(rand.Reader, maxRand)

	endValue := fmt.Sprintf("%s%d", startValue, randID)
	return strings.ReplaceAll(
			newConfig,
			startValue,
			endValue,
		), strings.ReplaceAll(
			modifyConfig,
			startValue,
			endValue,
		), endValue
}

func checkValidClientID(resourceName, attributeName, arnType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		clientID := rs.Primary.Attributes[attributeName]
		if len(clientID) <= 0 {
			return fmt.Errorf("empty client id: %s %s", resourceName, attributeName)
		}

		if len(strings.Split(clientID, ":")) != 6 {
			return fmt.Errorf("clientID does not have the correct number of blocks: %s", clientID)
		}
		if !strings.HasPrefix(clientID, "aembit:") {
			return fmt.Errorf("clientID does not have the correct prefix: %s", clientID)
		}
		if !strings.Contains(clientID, arnType) {
			return fmt.Errorf("clientID does not have the expected ARN type: %s", clientID)
		}
		if len(strings.Split(strings.Split(clientID, ":")[5], "-")) != 5 {
			return fmt.Errorf("clientID identifier is not a valid GUID: %s", clientID)
		}
		return nil
	}
}

func generateRandomPEMCertificate(t *testing.T, commonName string) string {
	var priv *rsa.PrivateKey
	var err error
	priv, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Errorf("Failed to generate private key: %v", err)
	}

	keyUsage := x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign
	var notBefore = time.Now()
	notAfter := notBefore.Add(time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		t.Errorf("Failed to generate serial number: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: commonName,
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              keyUsage,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template,
		&priv.PublicKey, priv)
	if err != nil {
		t.Errorf("Failed to create certificate: %v", err)
	}

	// Encode the pem.Block to a byte slice
	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	return string(pemBytes)
}
