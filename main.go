package main

import (
	"flag"
	"fmt"
	"github.com/marcochilese/Go-Trie"
	"github.com/marcochilese/negawordfixer/src/fsutils"
	"github.com/marcochilese/negawordfixer/src/processing"
	"io/ioutil"
	"os"
	"path"
)

func buildTrieAndReplacementDict(pathToDict string) (*trie.Trie, *map[string]string){
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
		if f.Name() == ".DS_Store"{
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


func main() {

	//pathToDict := flag.String("dict", "./dictionary_data/en.txt", "Path to language dictionary")
	tarPathPtr := flag.String("tar", "", "Path to negapedia-LANG.tar.gz")
	langPtr := flag.String("lang", "en", "Negapedia language")
	verbosePtr := flag.Bool("verbose", false, "Negapedia language")
	flag.Parse()

	*tarPathPtr = getNewestFileInDir(*tarPathPtr)
	pathToDict := path.Join("./dictionary_data/", *langPtr+".txt")

	logger := ioutil.Discard
	if *verbosePtr {
		logger = os.Stdout
	}
	
	fmt.Println("Run with config:\n\tLang: "+*langPtr+
		"\n\tDict: "+pathToDict+
		"\n\tTar: "+*tarPathPtr+"\n\t")

	mytrie, replacementDict := buildTrieAndReplacementDict(pathToDict)

	tmpDir, err := fsutils.ExtractTarGz(*tarPathPtr)
	if err != nil {
		fmt.Println(err)
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
	os.Mkdir("out", 0755)
	err = fsutils.CompressTarGz(tmpDir, path.Join(*tarPathPtr))
	if err != nil {
		fmt.Println(err)
	}
	os.RemoveAll(tmpDir)
	fmt.Fprintln(logger,"Compression end.")
	fmt.Fprintln(logger,"Done.")
}
