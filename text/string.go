package text

import "strings"

type String struct {
	text string
}

func NewString(str string) *String {
	return &String{text:str}
}

//判断字串是否为空
func (t *String)IsEmpty() bool {
	return t.text == "" && len(t.text) == 0
}

func (t *String)HasPostfix(suffix string)  bool {
	return strings.HasSuffix(t.text, suffix)
}