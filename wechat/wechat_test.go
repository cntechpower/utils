package wechat

import (
	"fmt"
	"testing"
)

func TestPrintSlice(t *testing.T) {
	s1 := make([]string, 0)
	s1 = append(s1, "a", "b", "c")
	s2 := make([]string, 0)
	fmt.Printf("%q\n", s1)
	fmt.Printf("%v\n", s2)

	fmt.Printf("%s\n", buildStringList(s1))
	fmt.Printf("%v\n", buildStringList(s2))
}
