package gcpkms

import (
	"context"
	"fmt"
	"sync"

	"github.com/cenkalti/backoff/v4"
	kmsrand "github.com/salrashid123/kmsrand"

	kms "cloud.google.com/go/kms/apiv1"
	kmspb "cloud.google.com/go/kms/apiv1/kmspb"
)

const (
	MIN_BYTES = 8 // https://cloud.google.com/kms/docs/reference/rest/v1/projects.locations/generateRandomBytes
	MAX_BYTES = 1024
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
			resp, err := r.Client.GenerateRandomBytes(context.Background(), &kmspb.GenerateRandomBytesRequest{
				Location:        r.Location,
				LengthBytes:     int32(len(chunk)),
				ProtectionLevel: r.ProtectionLevel,
			})
			if err != nil {
				fmt.Printf("%v\n", err)
				return err
			}
			result = append(result, resp.Data...)
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
