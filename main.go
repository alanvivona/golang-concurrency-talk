package main

import (
	"scrapper/fetcher"
	"scrapper/parser"
	"scrapper/utils"
	"sync"

	log "github.com/sirupsen/logrus"
)

type Scrapper struct {
	fetcher  *fetcher.Fetcher
	wg       *sync.WaitGroup
	maxDepth int
}

func NewScrapper(maxDepth int) *Scrapper {
	return &Scrapper{
		fetcher:  fetcher.NewFetcher(),
		wg:       &sync.WaitGroup{},
		maxDepth: maxDepth,
	}
}

func (s *Scrapper) Run(rootTarget *parser.Target) {
	s.run(rootTarget)
	s.wg.Wait()
}

func (s *Scrapper) run(target *parser.Target) error {
	log.Info("Fetching ", target.String())

	if target.Depth > s.maxDepth {
		log.Info("Skipping ", target.String())
		return nil
	}

	content, err := s.fetcher.Fetch(target.URL)
	if err != nil {
		return err
	}

	s.wg.Add(1)
	go func(t *parser.Target) {
		defer s.wg.Done()
		for _, newTarget := range t.ParseLinks(content) {
			s.run(&newTarget)
		}
	}(target)

	return nil
}

func main() {
	rootURL, maxDepth := utils.ParseEnv()
	rootTarget := &parser.Target{URL: *rootURL, Depth: 0}
	NewScrapper(maxDepth).Run(rootTarget)
}
