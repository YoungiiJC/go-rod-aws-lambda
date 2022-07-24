package main

/*
	Getting this error:
	go: github.com/aws/aws-lambda-go@v1.33.0 requires
	github.com/stretchr/testify@v1.7.2: missing go.sum entry; to add it:
	go mod download github.com/stretchr/testify

	...so importing testify here to force the import
	see https://chromium.googlesource.com/external/github.com/GoogleCloudPlatform/google-cloud-go-testing/+/f550565525113d92f31e16591986917006f86eb9/tools.go

*/

import _ "github.com/stretchr/testify"
