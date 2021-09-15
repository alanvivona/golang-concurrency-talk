package main

import (
	"scrapper/fetcher"
	"scrapper/parser"
	"scrapper/utils"

	log "github.com/sirupsen/logrus"
)

func main() {
	rootURL, maxDepth := utils.ParseEnv()
	fetcher := fetcher.NewFetcher()
	targets := []parser.Target{{*rootURL, 0}}
	for i := 0; i < len(targets); i++ {
		target := targets[i]
		if target.Depth > maxDepth {
			log.Info("Skipping ", target.String())
			continue
		}
		log.Info("Fetching ", target.String())
		content, err := fetcher.Fetch(target.URL)
		if err != nil {
			continue
		}
		newTargets := target.ParseLinks(content)
		targets = append(targets, newTargets...)
	}
}
