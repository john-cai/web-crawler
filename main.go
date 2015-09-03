package main

import (
	"fmt"
)

func main() {
	baseDomain := "www.digitalocean.com"
	parser := NewParser(baseDomain)
	parser.Parse(fmt.Sprintf("http://%s", baseDomain))
	parser.PrintChildren()
}
