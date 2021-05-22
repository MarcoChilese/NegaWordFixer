package main

import (
	"flag"
	"fmt"
	"github.com/marcochilese/Go-Trie"
	"github.com/marcochilese/negawordfixer/src/fsutils"
	"github.com/marcochilese/negawordfixer/src/processing"
	"io/ioutil"
	"os"
)

func buildTrieAndReplacementDict(pathToDict string) (*trie.Trie, *map[string]string){
	mytrie := trie.BuildTrieFromDictionary(pathToDict)

	// in replacementDict are stored all the replacements in order to
	// speedup the replacement process when stored enough history
	replacementDict := make(map[string]string)
	return mytrie, &replacementDict
}


func main() {

	pathToDictPtr := flag.String("dict", "./dictionary_data/en.txt", "Path to language dictionary")
	tarPathPtr := flag.String("tar", "./test_data/core-2020-01-24-negapedia-en.tar.gz", "Path to negapedia-LANG.tar.gz")
	langPtr := flag.String("lang", "en", "Negapedia language")
	verbosePtr := flag.Bool("verbose", false, "Negapedia language")

	logger := ioutil.Discard
	if *verbosePtr {
		logger = os.Stdout
	}
	
	fmt.Fprintln(logger, "Run with config:\n\tLang: %s\n\tDict: %s\n\tTar: %s\n\t", *langPtr, *pathToDictPtr, *tarPathPtr)

	mytrie, replacementDict := buildTrieAndReplacementDict(*pathToDictPtr)

	tmpDir, err := fsutils.ExtractTarGz(*tarPathPtr)
	if err != nil {
		fmt.Fprintln(logger,err)
	}

	filesToProcess := fsutils.GetFilesList(tmpDir, false)

	fmt.Fprintln(logger,"To process: ", len(filesToProcess))
	fmt.Fprintln(logger,"processing start")
	for _, file := range filesToProcess {
		err := processing.ProcessPage(file, *mytrie, replacementDict, &logger)
		if err != nil {
			os.RemoveAll(tmpDir)
		}

	}
	fmt.Fprintln(logger,"processing end.")
	fmt.Fprintln(logger,"Compression start")
	err = fsutils.CompressTarGz(tmpDir, "./out/FIXED.tar.gz")
	if err != nil {
		fmt.Fprintln(logger,err)
	}
	os.RemoveAll(tmpDir)
	fmt.Fprintln(logger,"Compression end.")
	fmt.Fprintln(logger,"Done.")
}
