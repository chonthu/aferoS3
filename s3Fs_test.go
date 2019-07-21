package aferoS3

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// hold aws sesssion in mem
var sess *session.Session

// setup return a aws session, doesnt recreate if exisits
func setup() *session.Session {
	if sess == nil {
		return session.Must(session.NewSession(&aws.Config{
			Region: aws.String("us-east-1"),
		}))
	}

	return sess
}

// TestChmod(name string, mode os.FileMode) : error
// TestChtimes(name string, atime time.Time, mtime time.Time) : error
// TestCreate(name string) : File, error
// TestMkdir(name string, perm os.FileMode) : error
// TestMkdirAll(path string, perm os.FileMode) : error
// TestName() : string
// TestOpen(name string) : File, error
// TestOpenFile(name string, flag int, perm os.FileMode) : File, error
func TestS3OpenFile(t *testing.T) {
	sess := setup()
	appFs := NewS3Fs(sess, "bucket_name")

	file, err := appFs.Create("output.jpg")
	if err != nil {
		t.Fatal(err)
	}

	info, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}

	if info.Name() != "output.jpg" {
		t.Fatal("couldnt open file")
	}
}

// TestRemove(name string) : error
// TestRemoveAll(path string) : error
// TestRename(oldname, newname string) : error
// TestStat(name string) : os.FileInfo, error
