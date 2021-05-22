package main

import (
	"github.com/marcochilese/NegaWordFixer/src/FsUtils"
	"github.com/marcochilese/NegaWordFixer/src/Processing"
	"flag"
	"fmt"
	"github.com/marcochilese/Go-Trie"
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
	tarPathPtr := flag.String("tar", "test_data/core-2020-01-24-negapedia-en.tar.gz", "Path to negapedia-LANG.tar.gz")
	langPtr := flag.String("lang", "en", "Negapedia language")
	verbosePtr := flag.Bool("verbose", false, "Negapedia language")

	logger := ioutil.Discard
	if *verbosePtr {
		logger = os.Stdout
	}
	
	fmt.Fprintln(logger, "Run with config:\n\tLang: %s\n\tDict: %s\n\tTar: %s\n\t", *langPtr, *pathToDictPtr, *tarPathPtr)

	mytrie, replacementDict := buildTrieAndReplacementDict(*pathToDictPtr)

	tmpDir, err := FsUtils.ExtractTarGz(*tarPathPtr)
	if err != nil {
		fmt.Fprintln(logger,err)
	}

	filesToProcess := FsUtils.GetFilesList(tmpDir, false)

	fmt.Fprintln(logger,"To process: ", len(filesToProcess))
	fmt.Fprintln(logger,"Processing start")
	for _, file := range filesToProcess {
		err := Processing.ProcessPage(file, *mytrie, replacementDict, &logger)
		if err != nil {
			os.RemoveAll(tmpDir)
		}

	}
	fmt.Fprintln(logger,"Processing end.")
	fmt.Fprintln(logger,"Compression start")
	err = FsUtils.CompressTarGz(tmpDir, "test_data/FIXED_core-2020-01-24-negapedia-en.tar.gz")
	if err != nil {
		fmt.Fprintln(logger,err)
	}
	os.RemoveAll(tmpDir)
	fmt.Fprintln(logger,"Compression end.")
	fmt.Fprintln(logger,"Done.")
}
