## KMS backed crypto/rand Reader   

A [crypto.rand](https://pkg.go.dev/crypto/rand) reader that uses a `Key Management System` (KMS) as the source of randomness.

Basically, its just a source of randomness used to create RSA keys or just get bits for use anywhere else but where you're asking a KMS system's hardware for the random bits.

As background, the default rand generator with golang uses the following sort-of random sources by default in [rand.go](https://go.dev/src/crypto/rand/rand.go)

but if you want HSM backed sources, you can use KMS or a TPM:

- [TPM and PKCS-11 backed crypto/rand Reader](https://github.com/salrashid123/tpmrand)
- [GCP KMS GenerateRandomBytes](https://cloud.google.com/kms/docs/samples/kms-generate-random-bytes)
- [AWS GenerateRandom](https://docs.aws.amazon.com/kms/latest/APIReference/API_GenerateRandom.html)

This repo implements just GCP and AWS's KMS...

KMS api operations to get random bytes has an associated consumption $ costs.

Just note that asking for random stuff isnt' free...nothing is free

>> this repo is *not* supported by google


update 4/22:  i realized after i wrote this on friday that the gcp kms backed on already exists [sethvargo/gcpkms-rand](https://github.com/sethvargo/gcpkms-rand).  The difference with this implementation is the interface and support for AWS.

---

From there, the usage is simple:

#### GCP 

```golang
package main

import (
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
$ gcloud auth application-default login

$ cd example/gcp

$ go run main.go --location=projects/srashid-test2/locations/us-central1
Random String :H2ttxPK3Dz+CeHVTmjSa84rVj6n/0nY5Ib8EOWa00W0=

RSA Key: 
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAy3WBNli/nAB8gyJeIFI70xHw2cm7CbUJGfGiSyJQt/HXJkMf
yAvMdp4vpaBnJagKmhFOHD8V3h3GUMyDnxHd8yF74T+CK1ODjwhoZJVlzJmmSA7L
MpBlSfqmcO0umyiu3qrqfelwsZab28P4GJN1ipNSqGVYNKF0NAokgZ5Y7W3tMw+H
G7sKvUsRbVasx79YiVyHCxzCnz+xEpBKkGD+L4dk8+xhkDiCs53ZevoOg8tm4Bhd
3AGoUJptKJaSfMNxNVkAryHDxaGNV6XOsoffWSY5LYrzk3vPeodwMm7iFJadyThS
kgdzOxE2GFRxjpZsKiAHp5ivYVEnI25G53wC8QIDAQABAoIBAHYcP6dp+8m3KpEB
uXyv4FTWfGghyLeI5cCu2lUdlZhDB3AJ1YBPASH3EJfotxhQJd9snlidcrdft4me
P+Zu+9axoHWRZaJ7N8snyVpitBcDN1lrZSB0XKiGnmq99alTA7j1pWz0wFwHn3ED
oZm6uKh6f6iMNJlRBOFU5f5tCxjBB47hmkw5QAN/kqMNJw5Zo+mnBmaxSpred1fK
UacIznJL/IiwqtuCUX1ckANTeU5tCFFAu1XtQqylNdBZI5SZwTG2dcVOJL2zbkXO
sVObQdYB/2gk4gDs1dMHiYL7CDK+ONdqHBq6C32+g3EkT6rGIZuHFLj8tmerdnr9
cVDQT7ECgYEA7omtMTv9SDAQ4cdq4NXGN8EjHw96VkUtH5NRxdDao2XvwGdeqQd+
8ioFiS3ZWubYFs7WQlO5eEj+EyNv+NPw6fCA5URpwPlM/fFwSgf8rYy7x0ROKyB6
eUBKAQqjhOFB+x7dEghhhGr+d/GbtTmRsO9jU1IzTXfNieHRCUsvKM0CgYEA2lpu
4+vzUfj6wkQeShaIfd4NmdGiiR+ZH8wrxNhTCXmd5SDSrm7Pj7EZjzJiSJ/qZFph
bS7MIwRH3hc31w/rgCj7c1uscZVY974HHyp0qBOs+i7rcYG6pRxqSQT5vP6IdQ9B
mR/ZXe2XvpPWtRWj/N8005vigNQqDwZqAv5c0rUCgYEAi2YDw4D2PGhyhS8/w1LK
eqywtKcb7CyS+R/jqsGp89FPcdY22Hrb8fMitw8HNXswDuwjBDHfcm7dpBuShQx+
foghG1qGntJR7xlYcLsIK/fRiNre/48EY7VxSfiIpM/q+jEIKlChhHvuZ/PW9epF
vOu41Ol1t7Dqechwm4jHb4UCgYAnhJhvLaPi4RHZGOT2ea+IQCjr/tHQyWQ4KgZ9
4LzeiSE3d8JJiYqNMfszPGYnSLHuKaFaVk7hw4OSQVd818fCcShZD21dPS9V3xGA
5Xkpdi4nNVitOVJjUYo23uyn9NUTgohXwzje1AJTnoQMT/dW67qu1ZafxEY8Y+fJ
1OlNxQKBgQCk+T5QFN5OaReXEa35WMIHKyyCozbhgRcxpOF8oLkcRoY58nrJ22mv
4OH0CvyFNd4m4epErkgp5G6WwOI2rQTHcTHqCGk0ksZDf4Qz+0aW7XPDJRpN/050
/IauDTmcwbVKPgKLF4Gqh6sN8mm9pkLSdtjnw9NcUzCEYI3F1vsbsg==
-----END RSA PRIVATE KEY-----
```

### AWS

```golang
package main

import (
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

	svc := kms.New(sess)

	randomBytes := make([]byte, 32)
	r, err := awskms.NewAWSRand(&awskms.AWSReader{
		Service: svc,
		//Scheme:    backoff.NewConstantBackOff(time.Millisecond * 10),
	})

	// Rand read
	_, err = r.Read(randomBytes)

	fmt.Printf("Random String: %s", base64.StdEncoding.EncodeToString(randomBytes))

	fmt.Println()

	/// RSA keygen
	privkey, err := rsa.GenerateKey(r, 2048)


	keyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privkey),
		},
	)
	fmt.Printf("RSA Key: \n%s\n", keyPEM)
}

```

```bash
export AWS_ACCESS_KEY_ID=AKIAUH3H6...
export AWS_SECRET_ACCESS_KEY=FZ2HR...

$ cd example/aws
$ go run main.go 

Random String: UoEasiXvmV81BCOALLqIZMZ063iYvdo8urmZP8K4kjc=
RSA Key: 
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA+VFdK/GC75gtYq0MBGHMowmZhApXhxDRyx8mpeC5g2DBAgrz
zK7Qy3CxWBkMytNjQi7c8iRHUPIy3V5P2i/I4aybQckw/a7JHs1FiB0M1w4itcKv
404126KEEkhYnrRgBlrHyHJTcYkhCcegFugWw3/fPxFEdic1+A90+J1SluqSedEW
pqha+I9jZwbrbNFH8vMtvK+BJAaqP+Vwkpko+uMukmkONQ8EfFrm8gAAay3x1bwg
Ag4Cik4OHS2hTn4MYgmYwkFBy+2md9Llur0s4Z/1fp46UT/OeaDgxXZjOx+vxA8r
8UctNpT3z7i9bqkRRSmOtB4/5B52qRYV3y6qfwIDAQABAoIBAQCDa8JDUbGFfqAd
7b3x6WOnZX4IvjLZPaJ5Adirg8QGXtAetYtCD7x8INE68SlvGPKvhmhtM3ZsUt9B
FV/eUWYAn63PhbBPaP0XQXkvgLCuBAOD8DYrCaUWO5qG0J/2OHqNnvjEzo7xwCks
MJBQwtKNBzC02/NMnOqz8eHk03kfl0h2NNuKToWeckwU1Zbu4KRb0nZBda+/KinK
S3K6rdY9JLBgXDukjXndDgq+iHpvmvWNtaFtBeN+Z/aoBywZWV/Olf2q7skV/JKf
U2FIS3GEqnun6tF2iNTwB7uiEHcDnuBSTgBeUgx9/mSOj5ep6CkLmOQlC+5PxEIF
lPyQQZmxAoGBAP+IccPJxUHdZIv8NXxlHrjNWoGxr/Wn+DEoypgA9pinyllJ1VkR
dYzFJarI4T3rhHi6XsL8EzgJe3LQQBwfj1I1cF23NikZlQMIOoEgkGxL1LNl59l9
QH0FJx0R3Ngi+dEDuqx6HSbvNzHAihKGckqC8i5SYl2RRbV+5hAwEC5dAoGBAPnG
Av3qhZRKIpjhPSc0daLt4Pt0F9PGNKyLdDU6Sq8waM9RWK0dfCnum97jkmJ1izri
Vxm/L4jUebQYhc2fphTzt76sv/wDPaUrs3seEyDj1yRguSeBiDjUyopYs/X7P4Z/
16ClnAj50TV9buoIXJPG3tqJEVxJUAyD2JfWd5aLAoGAFRsK8nXm4gLMPDevnz+m
4vKrKA0qEGs4N6871IQ32fH555gOlBW6FM9vxgRjfj7GqUYTb51sZPN7i8chlHES
4GJjjooEYi6nvSFf26x54Uf+IHcpSDBtNCZJzb/c8skowxfAwmAvqjiV4Xkarl8G
b5sTL7pEP6AxFsWNcQbXP00CgYBtVdBZdh+jGhCq+23Zi40zFQ43BEqp2UmVfjYQ
VsP6jCZVGjbHEPEZKenxV4zsrKeVzx5xls8oBlqAC3wG1qvM4CK+xMAFgSWq98ZJ
TpDxBMtYkT57nKgUuJEwnkOomaLlLXEmUVhMVY7O62lx6NcdmSBUaUvAKhdwYwac
8LTIoQKBgQDujtpdfoV17Y4yfA0W/1KyCxjbLS0XhaNn12/vRNCAo5s3VuXvo2nk
8krMiRQR9p3D1QEUcbyvjBUH+thUjAWWLwmNJZGhJrxFGMrGac7Kcp0CsZpJr262
clUulFOR0yB0Em6S8sg99tH8ZrlSefLNWnKRH188qTrV0bHX2K/3Og==
-----END RSA PRIVATE KEY-----

```

---
