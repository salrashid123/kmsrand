## KMS backed crypto/rand Reader   

A [crypto.rand](https://pkg.go.dev/crypto/rand) reader that uses a `Key Management System` (KMS) as the source of randomness.

Basically, its just a source of randomness used to create RSA keys or just get bits for use anywhere else but where you're asking a KMS system's hardware for the random bits.

As background, the default rand generator with golang uses the following sort-of random sources by default in [rand.go](https://go.dev/src/crypto/rand/rand.go)

but if you want HSM backed sources, you can use KMS or a TPM:

- [TPM and PKCS-11 backed crypto/rand Reader](https://github.com/salrashid123/tpmrand)
- [GCP KMS GenerateRandomBytes](https://cloud.google.com/kms/docs/samples/kms-generate-random-bytes)
- [AWS GenerateRandom](https://docs.aws.amazon.com/kms/latest/APIReference/API_GenerateRandom.html)

This repo implements just GCP's KMS...

KMS api operations to get random bytes has an associated consumption $ costs.

Just note that asking for random stuff isnt' free...nothing is free

>> this repo is *not* supported by google

---

From there, the usage is simple:

```golang
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

	//"time"

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

	defer kmsClient.Close()

	randomBytes := make([]byte, 32)
	r, err := gcpkms.NewGCPRand(&gcpkms.GCPReader{
		Client:          kmsClient,
		Location:        *location,
		ProtectionLevel: kmspb.ProtectionLevel_HSM,
		//Scheme:    backoff.NewConstantBackOff(time.Millisecond * 10),
	})

	// Rand read
	_, err = r.Read(randomBytes)


	fmt.Printf("Random String :%s\n", base64.StdEncoding.EncodeToString(randomBytes))


	fmt.Println()

	// /// RSA keygen
	privkey, err := rsa.GenerateKey(r, 512)

	keyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privkey),
		},
	)
	fmt.Printf("RSA Key: \n%s\n", keyPEM)
}
```

which gives:

```bash
$ go run main.go --location=projects/srashid-test2/locations/us-central1
Random String :H2ttxPK3Dz+CeHVTmjSa84rVj6n/0nY5Ib8EOWa00W0=

RSA Key: 
-----BEGIN RSA PRIVATE KEY-----
MIIBOwIBAAJBANSsm5YF/zZgMf0TlAY2dN+dnHGbXwyQLO2I3LdMRYM46kCKnZDn
UYPPbbNBfcs/Y/ENTmKKi/EZW9udaqTo05ECAwEAAQJBAIsKJdXZKdcE4OmRyS6e
n54qTsM/Ts7J23WYCqSTasa0VplBVAqEaCLPAnTMp+fG5ta+8o2PTScscW4Vaj/r
khkCIQD//n49OKMq4XHk22JTBLcuPjUwhqit4smru6kEwOctYwIhANSt3BFPVklp
qzZkFXBSE6LTd+a4OZw+47faafU4X3d7AiEAlnyGvXqUANsy1vRYorD89kQ/hF1E
v6O4JipVO6Qiwj0CICZRrPTxdnqDr3V9Ut+J6j/MGi5XwwmDy0O09qJYJdtBAiBQ
HGiiBYRz88f8Yd5qCdUhkL4s1fVpwC+57MRFgxVVYw==
-----END RSA PRIVATE KEY-----
```

---

While you're here, some other references on TPMs and usage:

* [Trusted Platform Module (TPM) recipes with tpm2_tools and go-tpm](https://github.com/salrashid123/tpm2)

---
