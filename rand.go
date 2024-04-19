package kmsrand

type RandSource interface {
	Read(data []byte) (n int, err error)
}
