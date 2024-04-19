package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	awskms "github.com/salrashid123/kmsrand/aws"
)

const ()

func main() {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2")},
	)
	if err != nil {
		fmt.Printf("Error creating session:  %v", err)
		os.Exit(-1)
	}
	svc := kms.New(sess)

	randomBytes := make([]byte, 32)
	r, err := awskms.NewAWSRand(&awskms.AWSReader{
		Service: svc,
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
	fmt.Printf("Random String: %s", base64.StdEncoding.EncodeToString(randomBytes))

	fmt.Println()

	/// RSA keygen
	privkey, err := rsa.GenerateKey(r, 2048)
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
