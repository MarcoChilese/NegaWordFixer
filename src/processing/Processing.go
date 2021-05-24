package processing

import (
	"fmt"
	trie "github.com/marcochilese/Go-Trie"
	"github.com/negapedia/negawordfixer/src/fsutils"
	"github.com/negapedia/negawordfixer/src/utils"
	"github.com/negapedia/negawordfixer/src/wikipage"
	"io"
	"strings"
)

func PerformCorrection(jsMap *[]utils.VarCouple, trie trie.Trie, replacementDict *map[string]string, logger *io.Writer) {
	for i, _ := range *jsMap {
		fmt.Fprint(*logger, (*jsMap)[i].Word+" replaced by ")

		newWord, exists := (*replacementDict)[(*jsMap)[i].Word]

		if exists { // if the replacement is already known, then use it and avoid the search on the trie
			if newWord == (*jsMap)[i].Word {
				fmt.Fprintln(*logger, "--SAME--")
			} else {
				(*jsMap)[i].Word = newWord
				fmt.Fprintln(*logger, newWord)
			}
		} else {
			// search for the alternatives on the trie
			alternatives := trie.PrefixSearch((*jsMap)[i].Word)
			if alternatives == nil || len(alternatives) == 0 { // no replacements available, skip
				fmt.Fprintln(*logger, "--NO RES--")
				continue
			}
			oldWord := (*jsMap)[i].Word
			//otherwise, replace the word with the shortest alternative, if the word is
			//a legal word, then it won't be replaced
			if alternatives[0] == oldWord {
				fmt.Fprintln(*logger, "--SAME--")
			} else {
				(*jsMap)[i].Word = alternatives[0] // it will be in position 0 since the alternatives are ordered by length
				fmt.Fprintln(*logger, (*jsMap)[i].Word)
			}
			(*replacementDict)[oldWord] = alternatives[0] // store the replacement for the future
		}
	}
}

func replaceJSVariable(pageData string, variableName string, trie trie.Trie, replacementDict *map[string]string, logger *io.Writer) (string, error) {
	if !strings.Contains(pageData, variableName) {
		return pageData, nil
	}

	startIdx := strings.Index(pageData, variableName)
	endIdx := strings.Index(pageData[startIdx:], "])") + 1

	varMapData := wikipage.ParseJSMap(pageData[startIdx : startIdx+endIdx])

	// deal with dictionary_data and do the words correction
	// call here
	PerformCorrection(varMapData, trie, replacementDict, logger)

	newJsMap := wikipage.GetJSMapFromSlice(varMapData, variableName) // get back JS map format

	// do the replacement in pageData[startIdx:startIdx+endIdx]
	//newJsMap = "x = new Map([[\"PROVA\", 1234]])"
	//fmt.Fprintln(*logger,newJsMap)

	pageData = pageData[:startIdx] + newJsMap + pageData[startIdx+endIdx+1:]

	return pageData, nil
}

func ProcessPage(gzPagePath string, trie trie.Trie, replacementDict *map[string]string, logger *io.Writer) error {
	data, err := fsutils.ReadGzPage(gzPagePath)
	if err != nil {
		return err
	}
	data, err = replaceJSVariable(data, wikipage.NegaJSVariables().Tfidf, trie, replacementDict, logger)
	if err != nil {
		return err
	}
	data, err = replaceJSVariable(data, wikipage.NegaJSVariables().Badw, trie, replacementDict, logger)
	if err != nil {
		return err
	}
	data, err = replaceJSVariable(data, wikipage.NegaJSVariables().Word2Occur, trie, replacementDict, logger)
	if err != nil {
		return err
	}
	err = fsutils.WriteGzPage(gzPagePath, data)
	if err != nil {
		return err
	}

	return nil
}
