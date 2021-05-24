package fsutils

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func ExtractTarGz(tarpath string) (string, error) {
	gzipStream, err := os.Open(tarpath)
	if err != nil {
		fmt.Println("error")
		return "", err
	}

	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		fmt.Println(err)
		log.Fatal("ExtractTarGz: NewReader failed")
	}
	tarReader := tar.NewReader(uncompressedStream)

	tmpdir := "tmp"
	//os.RemoveAll(tmpdir)

	err = os.Mkdir(tmpdir, 0755)
	if err != nil {
		return "", err
	}

	for true {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("ExtractTarGz: Next() failed: %s", err.Error())
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path.Join(tmpdir, header.Name), 0755); err != nil {
				log.Fatalf("ExtractTarGz: Mkdir() failed: %s", err.Error())
			}
		case tar.TypeReg:
			outFile, err := os.Create(path.Join(tmpdir, header.Name))
			if err != nil {
				log.Fatalf("ExtractTarGz: Create() failed: %s", err.Error())
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				log.Fatalf("ExtractTarGz: Copy() failed: %s", err.Error())
			}
			outFile.Close()

		default:
			log.Fatalf(
				"ExtractTarGz: uknown type: %s in %s",
				header.Typeflag,
				header.Name)
		}
	}
	return tmpdir, nil
}

func ExtractTarGz2(tarpath string) (string, error) {
	tmpDir := "./tmp"
	os.RemoveAll(tmpDir)

	err := os.Mkdir(tmpDir, 0755)
	if err != nil {
		log.Fatal(err)
	}

	commandArgs := []string{tarpath, tmpDir}
	extractionCmd := exec.Command("./extract.sh", commandArgs...)
	fmt.Println("\t" + extractionCmd.String())

	var cmdStderr bytes.Buffer
	extractionCmd.Stderr = &cmdStderr

	if err = extractionCmd.Run(); err != nil {
		fmt.Println(err)
		log.Fatal("Call to external command failed, with the following error stream:\n" + cmdStderr.String())
		return tmpDir, err
	}

	return tmpDir, err
}

func CompressTarGz(src string, outArchivePath string) error {

	// ensure the src actually exists before trying to tar it
	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("Unable to tar files - %v", err.Error())
	}

	//mw := io.MultiWriter(writers...)
	outF, err := os.Create(outArchivePath)
	defer outF.Close()
	if err != nil {
		return err
	}

	gzw := gzip.NewWriter(outF)
	defer gzw.Close()

	tw := tar.NewWriter(gzw)
	defer tw.Close()

	// walk path
	return filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {

		// return on any error
		if err != nil {
			return err
		}

		if !fi.Mode().IsRegular() {
			return nil
		}

		// create a new dir/file header
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		// update the name to correctly reflect the desired destination when untaring
		header.Name = strings.TrimPrefix(strings.Replace(file, src, "", -1), string(filepath.Separator))

		// write the header
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// open files for taring
		f, err := os.Open(file)
		if err != nil {
			return err
		}

		// copy file dictionary_data into tar writer
		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		// manually close here after each file operation; defering would cause each file close
		// to wait until all operations have completed.
		f.Close()

		return nil
	})
}



func CompressTarGz2(archiveFileName string, dirToCompress string) error {
	commandArgs := []string{archiveFileName, dirToCompress}
	extractionCmd := exec.Command("./compress.sh", commandArgs...)
	fmt.Println("\t" + extractionCmd.String())

	var cmdStderr bytes.Buffer
	extractionCmd.Stderr = &cmdStderr

	if err := extractionCmd.Run(); err != nil {
		fmt.Println(err)
		log.Fatal("Call to external command failed, with the following error stream:\n" + cmdStderr.String())
		return err
	}

	return nil
}


func GetFilesList(path string, includeDir bool) []string {
	var files []string

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !includeDir {
			if info.IsDir() {
				return nil
			}
		}

		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}

func ReadGzPage(gzPagePath string) (string, error) {
	f, err := os.Open(gzPagePath)
	if err != nil {
		return "", err
	}

	gz, err := gzip.NewReader(f)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	s, err := ioutil.ReadAll(gz)
	if err != nil {
		log.Fatal(err)
		return "nil", err
	}

	err = gz.Close()
	if err != nil {
		return "", err
	}
	err = f.Close()
	if err != nil {
		return "", err
	}

	return string(s), nil
}

func WriteGzPage(tarPath string, data string) error {
	//fmt.Println("gzing: ", tarPath)
	//_ = os.Remove(tarPath)

	/*fmt.Fprintf(os.Stdout,"%s\n", len(data))
	newGz, err := os.Create(tarPath)
	//defer newGz.Close()
	if err != nil {
		fmt.Println("Error writing page ", tarPath)
		return err
	}
	defer func(newGz *os.File) {
		err := newGz.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(newGz)

	writer := gzip.NewWriter(newGz)
	//defer writer.Close()
	_, err = writer.Write([]byte(data))
	if err != nil {
		return err
	}
	writer.Close()

	err = ioutil.WriteFile(tarPath, []byte(data), 0666)
	if err != nil {
		log.Fatal()
		return err
	}
	return err*/

	return WriteZip(tarPath, data)
}
