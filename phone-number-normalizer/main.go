package main

import (
	"bytes"
	"strings"
)

//func main() {
//	scanner := bufio.NewScanner(os.Stdin)
//	scanner.Scan()
//	input := scanner.Text()
//	normalizePhone := normalizer(input)
//	fmt.Println(normalizePhone)
//}

func normalizer(phone string) string {
	var buf bytes.Buffer
	if strings.HasPrefix(phone, "+6") {
		phone = strings.TrimPrefix(phone, "+6")
	} else if strings.HasPrefix(phone, "6") {
		phone = strings.TrimPrefix(phone, "6")
	}
	for _, char := range phone {
		if char >= '0' && char <= '9' {
			buf.WriteRune(char)
		}
	}
	if len(buf.String()) > 11 {
		return "false"
	}
	return buf.String()
}