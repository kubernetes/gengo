package main

import (
	"fmt"

	v1 "k8s.io/gengo/v2/parser/testdata/type-alias/v1"
	v2 "k8s.io/gengo/v2/parser/testdata/type-alias/v2"
)

func main() {
	b1 := v1.Blah{}
	b2 := v2.Blah{}

	fmt.Println(b1, b2)
}
