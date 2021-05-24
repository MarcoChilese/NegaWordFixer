package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/marcochilese/Go-Trie"
	"github.com/negapedia/negawordfixer/src/fsutils"
	"github.com/negapedia/negawordfixer/src/processing"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

func buildTrieAndReplacementDict(pathToDict string) (*trie.Trie, *map[string]string) {
	mytrie := trie.BuildTrieFromDictionary(pathToDict)

	// in replacementDict are stored all the replacements in order to
	// speedup the replacement process when stored enough history
	replacementDict := make(map[string]string)
	return mytrie, &replacementDict
}

func getNewestFileInDir(dir string) string {
	if dir[len(dir)-1] != "/"[0] {
		dir += "/"
	}

	files, _ := ioutil.ReadDir(dir)
	var newestFile string
	var newestTime int64 = 0
	for _, f := range files {
		if f.Name() == ".DS_Store" {
			continue
		}
		if !strings.Contains(f.Name(), "tar.gz") {
			continue
		}
		fi, err := os.Stat(dir + f.Name())
		if err != nil {
			fmt.Println(err)
		}
		currTime := fi.ModTime().Unix()
		if currTime > newestTime {
			newestTime = currTime
			newestFile = f.Name()
		}
	}
	return path.Join(dir, newestFile)
}

func languageChecker(lang string) (string, error) {
	available_lang := map[string]string{
		"en":     "en",
		"simple": "en",
		"it":     "en",
		"fr":     "fr",
		"es":     "es",
		"de":     "de",
	}

	if _, available := available_lang[lang]; !available {
		return "", errors.New("Language " + lang + " not available")
	}
	return available_lang[lang], nil
}

func main() {
	langPtr := flag.String("lang", "en", "Negapedia language")
	verbosePtr := flag.Bool("verbose", false, "Negapedia language")
	exactFilePtr := flag.String("file", "", "Exact tar.gz to process. If not specified is automatically selected the newest .tar.gz")
	outFilenamePtr := flag.String("out", "", "Name of the tar.gz output file")
	flag.Parse()

	dictLang, err := languageChecker(*langPtr)
	if err != nil {
		log.Fatal(err)
	}
	tarPathPtr := "./out"
	if *exactFilePtr != "" {
		tarPathPtr = path.Join("./out", *exactFilePtr)
	} else {
		tarPathPtr = getNewestFileInDir(tarPathPtr)
	}

	if *outFilenamePtr == "" {
		name := "fixed-"+path.Base(tarPathPtr)
		*outFilenamePtr = path.Join("./out", name)
	} else {
		*outFilenamePtr = path.Join("./out", *outFilenamePtr)
	}

	pathToDict := path.Join("./dictionary_data/", dictLang+".txt")

	logger := ioutil.Discard
	if *verbosePtr {
		logger = os.Stdout
	}

	fmt.Println("--- NegaWordsFixer ---")
	fmt.Println("Run with config:\n\tLang: " + *langPtr +
		"\n\tDict: " + pathToDict +
		"\n\tTar: " + tarPathPtr + "\n\t")

	mytrie, replacementDict := buildTrieAndReplacementDict(pathToDict)

	fmt.Println("Extraction start")
	start := time.Now()
	tmpDir, err := fsutils.ExtractTarGz2(tarPathPtr)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Extraction done in ", time.Now().Sub(start))

	filesToProcess := fsutils.GetFilesList(tmpDir, false)

	fmt.Println("To process: ", len(filesToProcess))
	fmt.Println("Processing start")
	start = time.Now()

	for _, file := range filesToProcess {
		err := processing.ProcessPage(file, *mytrie, replacementDict, &logger)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Processing done in ", time.Now().Sub(start))

	fmt.Fprintln(logger, "Compression start")
	start = time.Now()
	err = fsutils.CompressTarGz2(*outFilenamePtr, path.Join(tmpDir, "html"))
	if err != nil {
		fmt.Println(err)
	}
	os.RemoveAll(tmpDir)
	fmt.Println("Compression done in ", time.Now().Sub(start))
	fmt.Println("Done.")
}
