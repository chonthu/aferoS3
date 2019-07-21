package aferoS3

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// hold aws sesssion in mem
var sess *session.Session

// setup return a aws session, doesnt recreate if exisits
func setup() *session.Session {

	if sess == nil {

		// server is the mock server that simply writes a 200 status back to the client
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		return session.Must(session.NewSession(&aws.Config{
			DisableSSL:  aws.Bool(true),
			Endpoint:    aws.String(server.URL),
			Credentials: credentials.NewStaticCredentials("AKID", "SECRET", "SESSION"),
			Region:      aws.String("mock-region"),
		}))
	}

	return sess
}

// TestChmod(name string, mode os.FileMode) : error
// TestChtimes(name string, atime time.Time, mtime time.Time) : error
// TestCreate(name string) : File, error
func TestCreate(t *testing.T) {
	appFs := NewS3Fs(setup(), "bucket_name")

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

// TestMkdir(name string, perm os.FileMode) : error
// TestMkdirAll(path string, perm os.FileMode) : error
// TestName() : string
func TestName(t *testing.T) {
	appFs := NewS3Fs(setup(), "bucket_name")
	if appFs.Name() != "S#Fs" {
		t.Fatal("unknown module name")
	}
}

// TestOpen(name string) : File, error
func TestOpen(t *testing.T) {
	appFs := NewS3Fs(setup(), "bucket_name")
	file, err := appFs.Open("output.jpg")
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

// TestOpenFile(name string, flag int, perm os.FileMode) : File, error
func TestS3OpenFile(t *testing.T) {
	appFs := NewS3Fs(setup(), "bucket_name")

	file, err := appFs.OpenFile("output.jpg", 0, 0777)
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
