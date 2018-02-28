package validate

import (
	"testing"
	"fmt"
)

func TestIsPhone(t *testing.T) {
	fmt.Println(IsPhone("15555111700"))
	fmt.Println(IsPhone("23555111700"))
}

func TestIsEmail(t *testing.T) {
	fmt.Println(IsEmail("ymex.cn@gmail.com"))
	fmt.Println(IsEmail("ymex.mail.com"))
}