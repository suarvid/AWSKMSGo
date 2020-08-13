package filehandler

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// FileHandler stores filepaths for files when
// encrypting/decrypting and uploading/downloading
// the same files repeatedly.
// Defines behaviour for reading and writing files.
type FileHandler struct {
	plaintextPath string
	encryptedPath string
	downloadPath  string
	decryptedPath string
}

// NewHandler returns a new FileHandler with the
// specified filepaths for ease of use.
func NewHandler(plaintextPath string, encryptedPath string, downloadPath string, decryptedPath string) FileHandler {
	handler := new(FileHandler)
	handler.plaintextPath = plaintextPath
	handler.encryptedPath = encryptedPath
	handler.downloadPath = downloadPath
	handler.decryptedPath = decryptedPath
	return *handler
}

// GetFileHandle returns a file handle for reading and
// writing file in a less naive way
func (f *FileHandler) GetFileHandle(path string) *os.File {
	if !f.FileExists(path) {
		os.Create(path)
	}
	fileHandle, err := os.OpenFile(path, os.O_RDWR, os.ModeAppend)
	if err != nil {
		fmt.Printf("Error getting handle for file %s ", path)
		log.Fatal(err)
	}
	return fileHandle
}

// ReadFile returns the content of the file with the given
// path as a byte slice.
func (f *FileHandler) ReadFile(path string) []byte {
	data, err := ioutil.ReadFile(path)
	f.checkError(err)
	return data
}

// WriteFile writes a byte slice to the given filepath.
func (f *FileHandler) WriteFile(path string, data []byte) {
	file, err := os.Create(path)
	f.checkError(err)
	defer file.Close()
	file.Write(data)
}

// FileExists returns a boolean over wheter the file
// with the given path exists.
func (f *FileHandler) FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (f *FileHandler) checkError(err error) {
	if err != nil {
		fmt.Println("Error in fileHandler.go: ")
		log.Fatal(err)
	}
}

// GetPlaintextPath returns the path to the plaintext file
// used by the FileHandler
func (f *FileHandler) GetPlaintextPath() string {
	return f.plaintextPath
}

// GetDecryptedPath returns the path to the decrypted file
// used by the FileHandler
func (f *FileHandler) GetDecryptedPath() string {
	return f.decryptedPath
}

// GetDownloadPath returns the path to the downloaded file
// used by the FileHandler
func (f *FileHandler) GetDownloadPath() string {
	return f.downloadPath
}

// GetEncryptedPath returns the path to the encrypted file
// used by the FileHandler
func (f *FileHandler) GetEncryptedPath() string {
	return f.encryptedPath
}
