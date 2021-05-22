package main

import (
	"./FsUtils"
	"./Utils"
	"./WikiPage"
	"bufio"
	"fmt"
	"github.com/marcochilese/Go-Trie"
	"log"
	"os"
	"strings"
)

func BuildTrieFromDictionary(pathToDict string) *trie.Trie {
	gotrie := trie.NewTrie()

	file, err := os.Open(pathToDict)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := scanner.Text()
		gotrie.AddWord(word)
	}
	return gotrie
}

func PerformoCorrection(jsMap *[]Utils.VarCouple, trie trie.Trie, replacementDict *map[string]string) {
	for i, _ := range *jsMap {
		fmt.Print((*jsMap)[i].Word + " replaced by ")

		if newWord, exists := (*replacementDict)[(*jsMap)[i].Word]; exists { // if the replacement is already known, then use it and avoid the search on the trie
			if newWord == (*jsMap)[i].Word {
				fmt.Println("--SAME--")
			} else {
				(*jsMap)[i].Word = newWord
				fmt.Println(newWord)
			}
		} else {
			// search for the alternatives on the trie
			alternatives := trie.PrefixSearch((*jsMap)[i].Word)
			if alternatives == nil || len(alternatives) == 0 { // no replacements available, skip
				fmt.Println("--NO RES--")
				continue
			}
			oldWord := (*jsMap)[i].Word
			//otherwise, replace the word with the shortest alternative, if the word is
			//a legal word, then it won't be replaced
			if alternatives[0] == oldWord {
				fmt.Println("--SAME--")
			} else {
				(*jsMap)[i].Word = alternatives[0] // it will be in position 0 since the alternatives are ordered by length
				fmt.Println((*jsMap)[i].Word)
			}
			(*replacementDict)[oldWord] = alternatives[0] // store the replacement for the future
		}
	}
}

func replaceJSVariable(pageData string, variableName string, trie trie.Trie, replacementDict *map[string]string) (string, error) {
	if !strings.Contains(pageData, variableName) {
		return pageData, nil
	}

	startIdx := strings.Index(pageData, variableName)
	endIdx := strings.Index(pageData[startIdx:], "])") + 1

	varMapData := WikiPage.ParseJSMap(pageData[startIdx : startIdx+endIdx])

	// deal with dictionary_data and do the words correction
	// call here
	PerformoCorrection(varMapData, trie, replacementDict)

	newJsMap := WikiPage.GetJSMapFromSlice(varMapData, variableName) // get back JS map format

	// do the replacement in pageData[startIdx:startIdx+endIdx]
	//newJsMap = "x = new Map([[\"PROVA\", 1234]])"
	//fmt.Println(newJsMap)

	pageData = pageData[:startIdx] + newJsMap + pageData[startIdx+endIdx+1:]

	return pageData, nil
}

func ProcessPage(gzPagePath string, trie trie.Trie, replacementDict *map[string]string) error {
	data, err := FsUtils.ReadGzPage(gzPagePath)
	if err != nil {
		return err
	}
	data, err = replaceJSVariable(data, WikiPage.NegaJSVariables().Tfidf, trie, replacementDict)
	if err != nil {
		return err
	}
	data, err = replaceJSVariable(data, WikiPage.NegaJSVariables().Badw, trie, replacementDict)
	if err != nil {
		return err
	}
	data, err = replaceJSVariable(data, WikiPage.NegaJSVariables().Word2Occur, trie, replacementDict)
	if err != nil {
		return err
	}
	err = FsUtils.WriteGzPage(gzPagePath, data)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	mytrie := BuildTrieFromDictionary("dictionary_data/en_words_alpha.txt")

	// in replacementDict are stored all the replacements in order to
	// speedup the replacement process when stored enough history
	replacementDict := make(map[string]string)

	tarPath := "test_data/core-2020-01-24-negapedia-en.tar.gz"
	tmpDir, err := FsUtils.ExtractTarGz(tarPath)
	if err != nil {
		fmt.Println(err)
	}

	filesToProcess := FsUtils.GetFilesList(tmpDir, false)

	fmt.Println("To process: ", len(filesToProcess))
	fmt.Println("Processing start")
	for _, file := range filesToProcess {
		err := ProcessPage(file, *mytrie, &replacementDict)
		if err != nil {
			os.RemoveAll(tmpDir)
		}

	}
	fmt.Println("Processing end.")
	fmt.Println("Compression start")
	err = FsUtils.CompressTarGz(tmpDir, "test_data/FIXED_core-2020-01-24-negapedia-en.tar.gz")
	if err != nil {
		fmt.Println(err)
	}
	os.RemoveAll(tmpDir)
	fmt.Println("Compression end.")
	fmt.Println("Done.")
}
