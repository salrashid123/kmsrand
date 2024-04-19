package main

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"flag"
	"fmt"
	"log"

	//"github.com/cenkalti/backoff/v4"
	gcpkms "github.com/salrashid123/kmsrand/gcp"

	cloudkms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	//kmspb "cloud.google.com/go/kms/apiv1/kmspb"
)

var (
	location = flag.String("location", "projects/srashid-test2/locations/us-central1", "location used to generate random")
)

func main() {

	ctx := context.Background()
	kmsClient, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		log.Fatalf("Error getting kms client %v", err)
	}
	defer kmsClient.Close()

	randomBytes := make([]byte, 32)
	r, err := gcpkms.NewGCPRand(&gcpkms.GCPReader{
		Client:          kmsClient,
		Location:        *location,
		ProtectionLevel: kmspb.ProtectionLevel_HSM,
		//Scheme:    backoff.NewConstantBackOff(time.Millisecond * 10),
	})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	// Rand read
	_, err = r.Read(randomBytes)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("Random String :%s\n", base64.StdEncoding.EncodeToString(randomBytes))

	fmt.Println()

	// /// RSA keygen
	privkey, err := rsa.GenerateKey(r, 512)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	keyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privkey),
		},
	)
	fmt.Printf("RSA Key: \n%s\n", keyPEM)
}
