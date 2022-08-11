package sign

import (
	"fmt"
	"testing"
)

func TestVUserExists(t *testing.T) {
	ok, err := vUserExistsCheck("13980824611")
	fmt.Println(ok, err)
}

func TestVUserLogin(t *testing.T) {
	ok, err := vUserLoginCheck("13980824611", "123456")
	fmt.Println(ok, err)
	ok, err = vUserLoginCheck("13980824611", "12345")
	fmt.Println(ok, err)
}
