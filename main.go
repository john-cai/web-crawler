package main

import (
	"flag"
	"fmt"
	"log"
	"log/syslog"
)

var domain string

func init() {
	logwriter, e := syslog.New(syslog.LOG_NOTICE, "myprog")
	if e == nil {
		log.SetOutput(logwriter)
	}
}

func main() {
	flag.StringVar(&domain, "domain", "", "website to crawl eg: www.digitalocean.com")
	flag.Parse()

	if domain == "" {
		fmt.Println("Must provide domain")
		return
	}
	parser := NewParser(domain)
	parser.Parse()
	parser.PrintChildren()
}
