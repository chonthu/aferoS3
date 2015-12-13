package aferoS3

import (
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"github.com/spf13/afero"
	"testing"
)

func TestS3OpenFile(t *testing.T) {

	auth, err := aws.EnvAuth()
	if err != nil {
		t.Fatal(err)
	}

	client := s3.New(auth, aws.USEast)
	bucket := client.Bucket("bucket_name")

	var AppFs afero.Fs = S3Fs{
		Bucket: bucket,
	}

	file, err := AppFs.Open("path/to/file")

	if err != nil {
		t.Error(err)
	}

	if e, _ := file.Stat(); e == nil {
		t.Error("Corrupted file read")
	}

	var OsFs afero.Fs = afero.OsFs{}
	newFile, err := OsFs.Create("output.jpg")
	io.Copy(newFile, file)
}
