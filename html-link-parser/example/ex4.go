package main

import (
	"fmt"
	"strings"

	"github.com/zhenyang0405/gophercises/html-link-parser/parser"
)

var exampleHtml4 = `
<html>
<body>
  <a href="/dog-cat">dog cat <!-- commented text SHOULD NOT be included! --></a>
</body>
</html>
`

func main() {
	r := strings.NewReader(exampleHtml4)
	links, err := link.Parse(r)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", links)
}