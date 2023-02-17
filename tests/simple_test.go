package tests

import (
	"fmt"
	"strings"
	"testing"
)

func TestString(t *testing.T) {
	a := []string{
		"v_CricketShot_g11_c06",
		"1001_test.mp4_1676474809000671000.mp4",
	}
	for _, s := range a {
		if strings.HasPrefix(s, "v_") && strings.HasPrefix(s, ".mp4") {
			fmt.Println(s[strings.Index(s, "_")+1 : 2+strings.Index(s[2:], "_")])
		}

	}
}
