package utils

type VarCouple struct {
	Word, Value string
}

func (v VarCouple) GetStrList() string {
	return "[\"" + v.Word + "\"," + v.Value + "],"
}
