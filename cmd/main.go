package main

import (
	"encoding/json"
	"fmt"
	tokenizer "github.com/yurii-stakhiv/atlassian"
	"os"
)

func showUsage() {
	fmt.Println("./cmd <string>")
	os.Exit(0)
}

func main() {
	if len(os.Args) < 2 {
		showUsage()
	}

	scanner := tokenizer.NewScanner(nil)
	res, err := scanner.ScanBytes([]byte(os.Args[1]))
	if err != nil {
		fmt.Println(err)
		return
	}
	b, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}
