# afero-s3

Afero S3 is a Afero FS interface for Amazon s3

## How to use

	auth, err := aws.EnvAuth()
	if err != nil {
		t.Fatal(err)
	}

	client := s3.New(auth, aws.USEast)
	bucket := client.Bucket("spin.com")

	var AppFs afero.Fs = S3Fs{
		Bucket: bucket,
	}

	file, err := AppFs.Open("files/00-01-spin-cover.jpg")

	if err != nil {
		t.Error(err)
	}

	fmt.Println(file)