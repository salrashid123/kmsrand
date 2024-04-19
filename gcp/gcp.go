package gcpkms

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/cenkalti/backoff/v4"
	kmsrand "github.com/salrashid123/kmsrand"

	kms "cloud.google.com/go/kms/apiv1"
	kmspb "cloud.google.com/go/kms/apiv1/kmspb"
)

const (
	MAX_BYTES = 1024 // https://cloud.google.com/kms/docs/reference/rest/v1/projects.locations/generateRandomBytes
)

type GCPReader struct {
	kmsrand.RandSource
	Scheme          backoff.BackOff
	Client          *kms.KeyManagementClient
	Location        string
	ProtectionLevel kmspb.ProtectionLevel
	mu              sync.Mutex
}

func NewGCPRand(conf *GCPReader) (*GCPReader, error) {
	if conf.Client == nil {
		return &GCPReader{}, fmt.Errorf("kms client must be set")
	}

	// todo, regex check the location is formatted
	if conf.Location == "" {
		return &GCPReader{}, fmt.Errorf("kms location must be set")
	}

	if conf.Scheme == nil {
		conf.Scheme = backoff.NewExponentialBackOff()
	}
	return conf, nil
}

func (r *GCPReader) Read(data []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(data) > MAX_BYTES {
		return 0, errors.New("kmsrand: Number of bytes to read exceeds cannot 1024")
	}
	originalLength := len(data) // sometimes len(data) < 8
	if len(data) < 8 {
		data = make([]byte, 8)
	}
	var result []byte
	operation := func() (err error) {
		resp, err := r.Client.GenerateRandomBytes(context.Background(), &kmspb.GenerateRandomBytesRequest{
			Location:        r.Location,
			LengthBytes:     int32(len(data)),
			ProtectionLevel: r.ProtectionLevel,
		})
		if err != nil {
			fmt.Printf("%v\n", err)
			return err
		}
		result = resp.Data[:originalLength] // reset to the original ask of random data
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
