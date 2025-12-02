package provider

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	trustProviderSecretPath = "aembit_trust_provider_secret.secret"
)

func TestAccTrustProviderSecretResource(t *testing.T) {
	createFile, _ := os.ReadFile(
		"../../tests/trust_provider_secret/TestAccTrustProviderSecretResource.tf",
	)
	modifyFile, _ := os.ReadFile(
		"../../tests/trust_provider_secret/TestAccTrustProviderSecretResource.tfmod",
	)
	pemCertificate := generateRandomPEMCertificate(t, "Aembit Unit Test")
	updatedPemCertificate := generateRandomPEMCertificate(t, "Updated Aembit Unit Test")

	// Replace placeholder with generated certificate
	createFileString := strings.Replace(string(createFile), "PEM_CERTIFICATE", pemCertificate, -1)
	modifyFileString := strings.Replace(string(modifyFile), "PEM_CERTIFICATE", updatedPemCertificate, -1)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileString,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(trustProviderSecretPath, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(trustProviderSecretPath, "id"),
					// Verify Subject
					resource.TestCheckResourceAttr(
						trustProviderSecretPath,
						"subject",
						"CN=Aembit Unit Test",
					),
				),
			},
			// ImportState testing
			{ResourceName: trustProviderSecretPath, ImportState: true, ImportStateVerify: false},
			// Update and Read testing
			{
				Config: modifyFileString,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Subject updated
					resource.TestCheckResourceAttr(
						trustProviderSecretPath,
						"subject",
						"CN=Updated Aembit Unit Test",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func generateRandomPEMCertificate(t *testing.T, commonName string) string {
	var priv *rsa.PrivateKey
	var err error
	priv, err = rsa.GenerateKey(rand.Reader, 2048)

	keyUsage := x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign
	var notBefore time.Time = time.Now()
	notAfter := notBefore.Add(time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("Failed to generate serial number: %v", err)
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
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %v", err)
	}

	block := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	}

	// Encode the pem.Block to a byte slice
	pemBytes := pem.EncodeToMemory(block)

	return string(pemBytes)
}
