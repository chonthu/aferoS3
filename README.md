# afero-s3

Afero S3 is a Afero FS interface for Amazon s3

## Install

	go get github.com/spf13/afero
	go get github.com/chonthu/aferoS3

## How to use

	afero.Fs, err := aferoS3.GetBucket("bucket_name", aferoS3.USEast)
	if err != nil {
		fmt.Println(err)
		return
	}

	file, err := S3Fs.Open("path/to/file.jpg")

	if err != nil {
		fmt.Println(err)
		return
	}