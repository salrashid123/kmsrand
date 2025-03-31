package awskms

import (
	"fmt"

	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/cenkalti/backoff/v4"
	kmsrand "github.com/salrashid123/kmsrand"
)

const (
	MIN_BYTES = 1
	MAX_BYTES = 1024
)

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
	for i := 0; i < len(data); i += MAX_BYTES {
		end := i + MAX_BYTES
		if end > len(data) {
			end = len(data)
		}
		chunk := data[i:end]

		if len(chunk) < MIN_BYTES {
			chunk = make([]byte, MIN_BYTES)
		}

		operation := func() (err error) {

			input := &kms.GenerateRandomInput{
				NumberOfBytes: aws.Int64(int64(len(chunk))),
			}

			if r.CustomerKeyStoreId != "" {
				input.CustomKeyStoreId = aws.String(r.CustomerKeyStoreId)
			}

			randomBytes, err := r.Service.GenerateRandom(input)
			if err != nil {
				fmt.Printf("Error generating aws random:  %v", err)
				return err
			}

			result = append(result, randomBytes.Plaintext...)
			return nil
		}
		err = backoff.Retry(operation, r.Scheme)
		if err != nil {
			return 0, err
		}

	}
	copy(data, result)
	return len(result), err
}
