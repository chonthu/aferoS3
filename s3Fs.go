package remote

import (
	"fmt"
	"github.com/mitchellh/goamz/s3"
	"github.com/spf13/afero"
	"github.com/spf13/afero/mem"
	"io"
	"io/ioutil"
	"mime"
	"os"
	"strings"
	"time"
)

/*

I think here I should create a s3 to file system mapping?

Example:

	"private":                   600,
	"public-read":               664,
	"public-read-write":         666,
	"authenticated-read":        660,
	"bucket-owner-read":         660,
	"bucket-owner-full-control": 666,

*/

type S3Fs struct {
	Bucket *s3.Bucket
}

type S3File struct {
	*mem.File
}

func (S3Fs) Name() string {
	return "S3Fs"
}

func (S3Fs) Create(name string) (afero.File, error) {
	return S3File{}, nil
}

// Read from s3, and bring down whole file? or torrent?
func (s S3Fs) Open(name string) (afero.File, error) {

	memFile := mem.CreateFile(getNameFromPath(name))

	torrentReader, err := s.Bucket.GetTorrentReader(name)
	if err != nil {
		return memFile, err
	}

	if err != nil {
		if _, err = io.Copy(memFile, torrentReader); err != nil {
			return memFile, err
		}
	}
	return memFile, err
}

func (s S3Fs) Push(f afero.File, path string) error {

	body, err := ioutil.ReadAll(f)

	s.Bucket.Put(path, body, mime.TypeByExtension(path), "")

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

type S3FileInfo struct {
	os.FileInfo
	file *afero.File
}

// Maybe different between Open is its torrent, and this is the actuall file
func (s S3Fs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	file, err := s.Open(name)
	s.Chmod(name, perm)
	return file, err
}

// Set ACL Perms
func (s S3Fs) Chmod(name string, mode os.FileMode) error {
	return nil
}

func (S3Fs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return nil
}

func (s S3Fs) Stat(name string) (os.FileInfo, error) {
	f, err := s.Open(name)
	return S3FileInfo{file: &f}, err
}

// Renames a file
func (s S3Fs) Rename(oldname, newname string) error {
	return s.Bucket.Copy(oldname, newname, s3.ACL(""))
}

// Removes a file
func (s S3Fs) Remove(name string) error {
	return s.Bucket.Del(name)
}

// Dont think we can do much here
func (S3Fs) Mkdir(name string, perm os.FileMode) error    { return nil }
func (S3Fs) MkdirAll(path string, perm os.FileMode) error { return nil }
func (S3Fs) RemoveAll(path string) error                  { return nil }
