package feed

import (
	"fmt"
	"strings"
	"testing"
)

func TestString(t *testing.T) {
	s := "v_CricketShot_g11_c06"
	fmt.Println(s[strings.Index(s, "_")+1 : 2+strings.Index(s[2:], "_")])
}
