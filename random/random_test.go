package random

import (
	"testing"
	"fmt"
	"encoding/json"
)

func TestRandomChars(t *testing.T) {
	st32 := RandomChars(32)
	fmt.Println(st32)
	st4 := RandomChars(4)
	fmt.Println(st4)
	st64f := RandomChars(2, 'a', 'f')
	fmt.Println(st64f)
}
func TestRandomUUID(t *testing.T) {
	RandomUUID()
}

func TestJson(t *testing.T) {
	var jtext string = `{
		"password":"123456",
		"age":24,
		"account":{
			"email":"ymex.cn",
			"phone":"15555111700"
		}
	}`
	var v map[string]interface{} = make(map[string]interface{}, 3)
	if err := json.Unmarshal([]byte(jtext), &v); err == nil {
		if val, ok := v["age"]; ok {
			var st string = str(val)

			fmt.Println(st)
		} else {
			fmt.Println("key 不存在")
		}
	} else {
		fmt.Println(err)
	}

}

func str(reply interface{}) string {

	switch reply := reply.(type) {
	case []byte:
		return string(reply)
	case string:
		return reply
	case nil:
		return ""
	}
	return ""
}