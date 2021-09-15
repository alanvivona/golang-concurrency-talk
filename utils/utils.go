package utils

import (
	"net/url"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func ParseEnv() (rootURL *url.URL, maxDepth int) {
	rootURL, err := url.Parse(os.Getenv("ROOT"))
	if err != nil {
		log.Fatal(err)
	}
	maxDepth, err = strconv.Atoi(os.Getenv("DEPTH"))
	if err != nil {
		log.Fatal(err)
	}
	if maxDepth < 1 {
		log.Fatal("DEPTH should be bigger than 1")
	}
	log.Infof("Env root:%s max-depth:%d\n", rootURL, maxDepth)
	return
}
