package WikiPage

import (
	"../Utils"
	"strings"
)

type JSVariables struct {
	Tfidf string
	Badw  string
	Word2Occur string
}

func NegaJSVariables() JSVariables {
	return JSVariables{
		Tfidf: "Word2TFIDF",
		Badw:  "BWord2Occur",
		Word2Occur: "Word2Occur",
	}
}

func ParseJSMap(jsmap string) *[]Utils.VarCouple {
	/**
	Returns a slice made like [(word, value), ...] from a variable formatted like:
		Word2TFIDF = new Map([[\"2015\",  0.0047 ], ... ,])"
	*/
	var variable []Utils.VarCouple

	mapStart := strings.Index(jsmap, "Map([") + 5
	values := strings.Split(jsmap[mapStart:], "],")

	for _, keyvalue := range values {
		if keyvalue == "]" {
			continue
		}
		fields := strings.Split(keyvalue, ",")

		variable = append(variable, Utils.VarCouple{Word: strings.ToLower(fields[0][2:len(fields[0])-1]), Value: fields[1]})
	}

	return &variable
}

func GetJSMapFromSlice(dataSlice *[]Utils.VarCouple, varName string) (newVar string) {
	newVar = varName + " = new Map(["

	for _, couple := range *dataSlice {
		newVar += couple.GetStrList()
	}

	newVar += "])"
	return
}
