package fsutils

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"
)

func WriteZip(zipFile string, content string) error {
	handle, err := OpenFile(zipFile)
	if err != nil {
		fmt.Println("[ERROR] Opening file:", err)
		return err
	}

	zipWriter, err := gzip.NewWriterLevel(handle, 9)
	if err != nil {
		fmt.Println("[ERROR] New gzip writer:", err)
		return err
	}
	_, err = zipWriter.Write([]byte(content))
	if err != nil {
		fmt.Println("[ERROR] Writing:", err)
		return err
	}
	err = zipWriter.Close()
	if err != nil {
		fmt.Println("[ERROR] Closing zip writer:", err)
		return err
	}
	//fmt.Println("[INFO] Number of bytes written:", numberOfBytesWritten)

	CloseFile(handle)
	return nil
}

func ReadZip(zipFile string) []byte {
	handle, err := OpenFile(zipFile)
	if err != nil {
		fmt.Println("[ERROR] Opening file:", err)
	}

	zipReader, err := gzip.NewReader(handle)
	if err != nil {
		fmt.Println("[ERROR] New gzip reader:", err)
	}
	defer zipReader.Close()

	fileContents, err := ioutil.ReadAll(zipReader)
	if err != nil {
		fmt.Println("[ERROR] ReadAll:", err)
	}

	//fmt.Printf("[INFO] Uncompressed contents: %s\n", fileContents)

	// ** Another way of reading the file **
	//
	// fileInfo, _ := handle.Stat()
	// fileContents := make([]byte, fileInfo.Size())
	// bytesRead, err := zipReader.Read(fileContents)
	// if err != nil {
	//     fmt.Println("[ERROR] Reading gzip file:", err)
	// }
	// fmt.Println("[INFO] Number of bytes read from the file:", bytesRead)

	CloseFile(handle)
	return fileContents
}

func OpenFile(fileToOpen string) (*os.File, error) {
	return os.OpenFile(fileToOpen, OpenFileOptions, OpenFilePermissions)
}

func CloseFile(handle *os.File) {
	if handle == nil {
		return
	}

	err := handle.Close()
	if err != nil {
		fmt.Println("[ERROR] Closing file:", err)
	}
}

const OpenFileOptions int = os.O_CREATE | os.O_RDWR
const OpenFilePermissions os.FileMode = 0660
