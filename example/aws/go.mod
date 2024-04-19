module main

go 1.21

toolchain go1.21.0

require (
	github.com/aws/aws-sdk-go v1.51.25
	github.com/salrashid123/kmsrand/aws v0.0.0-00010101000000-000000000000
)

require (
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/salrashid123/kmsrand v0.0.0-20240419172308-e58279d05f07 // indirect
)

replace (
	github.com/salrashid123/kmsrand => ../../
	github.com/salrashid123/kmsrand/aws => ../../aws
)
