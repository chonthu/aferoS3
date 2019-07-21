package aferoS3

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/afero"
	"github.com/spf13/afero/mem"
)

// S3Fs main struct
type S3Fs struct {
	Bucket  string
	client  *s3.S3
	session *session.Session
}

//NewS3Fs create a new instance
func NewS3Fs(sess *session.Session, bucketName string) *S3Fs {
	return &S3Fs{
		Bucket:  bucketName,
		client:  s3.New(sess),
		session: sess,
	}
}

// Name the name of the module
func (S3Fs) Name() string {
	return "S3Fs"
}

// CreateBucket create a s3 bucket
func (s *S3Fs) CreateBucket(name string) error {
	_, err := s.client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(name),
	})
	return err
}

// Create a new file in memory, doesnt persits until pushing to s3
func (S3Fs) Create(name string) (afero.File, error) {
	return mem.NewFileHandle(mem.CreateFile(name)), nil
}

// Open from s3, and bring down whole file
func (s S3Fs) Open(name string) (afero.File, error) {

	memFile, err := s.Create(getNameFromPath(name))
	if err != nil {
		return nil, err
	}

	res, err := s.client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(name),
	})

	if err != nil {
		return memFile, err
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return memFile, err
	}

	memFile.Write(b)

	return memFile, err
}

// Push object to s3
func (s S3Fs) Push(f afero.File, path string) error {

	body, err := ioutil.ReadAll(f)

	if err != nil {
		return err
	}

	s.client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(path),
		Body:   bytes.NewReader(body),
	})

	return err
}

func getNameFromPath(fileName string) string {
	var name string
	tokens := strings.Split(fileName, ".")
	ext := tokens[len(tokens)-1]

	if len(tokens) > 2 {
		name = strings.Join(tokens[:len(tokens)-1], ".")
	} else {
		name = tokens[0]
	}

	return fmt.Sprintf("%s.%s", name, ext)
}

// S3FileInfo info
type S3FileInfo struct {
	os.FileInfo
	file *afero.File
}

// OpenFile different between Open is its torrent, and this is the actuall file
func (s S3Fs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	file, err := s.Open(name)
	s.Chmod(name, perm)
	return file, err
}

// Chmod not implemented
func (s S3Fs) Chmod(name string, mode os.FileMode) error {
	return nil
}

// Chtimes not implemented
func (S3Fs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return nil
}

// Stat get s3 stats
func (s S3Fs) Stat(name string) (os.FileInfo, error) {
	f, err := s.Open(name)
	return S3FileInfo{file: &f}, err
}

// Rename a file
func (s S3Fs) Rename(oldname, newname string) error {
	_, err := s.client.CopyObject(&s3.CopyObjectInput{
		CopySource: aws.String(fmt.Sprintf("%s/%s", s.Bucket, oldname)),
		Bucket:     aws.String(s.Bucket),
		Key:        aws.String(newname),
	})

	if err != nil {
		return err
	}

	// remove old one, if the copy worked
	return s.Remove(oldname)
}

// Remove a file
func (s S3Fs) Remove(name string) error {
	_, err := s.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(name),
	})
	return err
}

// ReadDir reads a current file and returns contents
func (s S3Fs) ReadDir(dirname string) ([]string, error) {
	keys := []string{}
	res, err := s.client.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(s.Bucket),
		Prefix: aws.String(dirname),
	})

	if err != nil {
		return nil, err
	}

	for _, obj := range res.Contents {
		keys = append(keys, *obj.Key)
	}

	for res.NextMarker != nil {
		res, err = s.client.ListObjects(&s3.ListObjectsInput{
			Bucket: aws.String(s.Bucket),
			Prefix: aws.String(dirname),
			Marker: res.NextMarker,
		})

		if err != nil {
			return nil, err
		}

		for _, obj := range res.Contents {
			keys = append(keys, *obj.Key)
		}
	}

	return keys, nil
}

// ReadDirNames(n int) : []string, error
// DirExists(path string) (bool, error)
// Exists(path string) (bool, error)

// IsDir check if this is a valid directory
func (s S3Fs) IsDir(path string) (bool, error) {

	_, err := s.client.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(s.Bucket),
		Prefix: aws.String(path),
	})

	if err != nil {
		return false, err
	}

	return true, nil
}

// IsEmpty(path string) (bool, error)
// Walk(root string, walkFn filepath.WalkFunc) error

// RemoveAll removes all files
func (s S3Fs) RemoveAll(dirname string) error {

	isdir, err := s.IsDir(dirname)

	if err != nil {
		return err
	}

	if isdir {
		return errors.New("invalid directory")
	}

	keys, err := s.ReadDir(dirname)

	// delete each of them
	for _, key := range keys {
		err = s.Remove(key)
		if err != nil {
			return err
		}
	}

	return nil
}

// Mkdir Dont think we can do much here
func (S3Fs) Mkdir(name string, perm os.FileMode) error { return nil }

// MkdirAll Dont think we can do much here
func (S3Fs) MkdirAll(path string, perm os.FileMode) error { return nil }
