package awskms

import (
	"fmt"

	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/cenkalti/backoff/v4"
	kmsrand "github.com/salrashid123/kmsrand"
)

const ()

type AWSReader struct {
	kmsrand.RandSource
	Scheme             backoff.BackOff
	Service            *kms.KMS
	CustomerKeyStoreId string
	mu                 sync.Mutex
}

func NewAWSRand(conf *AWSReader) (*AWSReader, error) {
	if conf.Service == nil {
		return &AWSReader{}, fmt.Errorf("kms client must be set")
	}

	if conf.Scheme == nil {
		conf.Scheme = backoff.NewExponentialBackOff()
	}
	return conf, nil
}

func (r *AWSReader) Read(data []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var result []byte
	operation := func() (err error) {
		input := &kms.GenerateRandomInput{
			NumberOfBytes: aws.Int64(int64(len(data))),
		}

		if r.CustomerKeyStoreId != "" {
			input.CustomKeyStoreId = aws.String(r.CustomerKeyStoreId)
		}

		randomBytes, err := r.Service.GenerateRandom(input)
		if err != nil {
			fmt.Printf("Error generating aws random:  %v", err)
			return err
		}

		result = randomBytes.Plaintext
		copy(data, result)

		return nil
	}

	// dont' know which scheme is better, probably the constant
	//err = backoff.Retry(operation, backoff.NewExponentialBackOff())
	err = backoff.Retry(operation, r.Scheme)
	if err != nil {
		return 0, err
	}

	return len(result), nil
}
