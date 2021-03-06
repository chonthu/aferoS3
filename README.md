# afero-s3

Afero S3 is a Afero FS interface for Amazon s3

## Install

	go get github.com/spf13/afero
	go get github.com/chonthu/aferoS3

## How to use

	// import afero, aferos3, and the aws sdk
	import (
		"github.com/aws/aws-sdk-go/aws"
		"github.com/aws/aws-sdk-go/aws/session"
		"github.com/chonthu/aferos3"
	)

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")
	}))

	S3Fs, err := aferoS3.NewS3Fs(sess,"bucket_name")
	if err != nil {
		fmt.Println(err)
		return
	}

	file, err := S3Fs.Open("path/to/file.jpg")

	if err != nil {
		fmt.Println(err)
		return
	}

## License

This project uses the https://github.com/aws/aws-sdk-go, thanks amazon web services for the great library

This SDK is distributed under the
[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0),
see LICENSE.txt and NOTICE.txt for more information.