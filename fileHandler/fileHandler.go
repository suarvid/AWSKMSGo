package FileHandler

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type FileHandler struct {
	plaintextPath string
	encryptedPath string
	downloadPath  string
	decryptedPath string
}

func NewHandler(plaintextPath string, encryptedPath string, downloadPath string, decryptedPath string) FileHandler {
	handler := new(FileHandler)
	handler.plaintextPath = plaintextPath
	handler.encryptedPath = encryptedPath
	handler.downloadPath = downloadPath
	handler.decryptedPath = decryptedPath
	return *handler
}

func (self *FileHandler) GetFileHandle(path string) *os.File {
	if !self.FileExists(path) {
		os.Create(path)
	}
	fileHandle, err := os.OpenFile(path, os.O_RDWR, os.ModeAppend)
	if err != nil {
		fmt.Printf("Error getting handle for file %s ", path)
		log.Fatal(err)
	}
	return fileHandle
}

func (self *FileHandler) ReadFile(path string) []byte {
	data, err := ioutil.ReadFile(path)
	self.checkError(err)
	return data
}

func (self *FileHandler) WriteFile(path string, data []byte) {
	file, err := os.Create(path)
	self.checkError(err)
	defer file.Close()
	file.Write(data)
}

func (self *FileHandler) FileExists(fileName string) bool {
	info, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (self *FileHandler) checkError(err error) {
	if err != nil {
		fmt.Println("Error in fileHandler.go: ")
		log.Fatal(err)
	}
}

func (self *FileHandler) GetPlaintextPath() string {
	return self.plaintextPath
}

func (self *FileHandler) GetDecryptedPath() string {
	return self.decryptedPath
}

func (self *FileHandler) GetDownloadPath() string {
	return self.downloadPath
}

func (self *FileHandler) GetEncryptedPath() string {
	return self.encryptedPath
}
