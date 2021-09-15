package main

import (
	"fmt"
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
	targetCh chan *parser.Target
	errCh    chan error
	endCh    chan bool
}

func NewScrapper(maxDepth int, targetCh chan *parser.Target, errCh chan error, endCh chan bool) *Scrapper {
	return &Scrapper{
		fetcher:  fetcher.NewFetcher(),
		wg:       &sync.WaitGroup{},
		maxDepth: maxDepth,
		targetCh: targetCh,
		errCh:    errCh,
		endCh:    endCh,
	}
}

func (s *Scrapper) Run(rootTarget *parser.Target) {
	s.run(rootTarget)
	s.wg.Wait()
	s.endCh <- true
}

func (s *Scrapper) run(target *parser.Target) {
	s.targetCh <- target

	if target.Depth > s.maxDepth {
		log.Info("Skipping ", target.String())
		return
	}

	content, err := s.fetcher.Fetch(target.URL)
	if err != nil {
		s.errCh <- fmt.Errorf("Failed to fetch %s", target.String())
		return
	}

	s.wg.Add(1)
	go func(t *parser.Target) {
		defer s.wg.Done()
		for _, newTarget := range t.ParseLinks(content) {
			s.run(&newTarget)
		}
	}(target)
}

func main() {
	rootURL, maxDepth := utils.ParseEnv()
	rootTarget := &parser.Target{URL: *rootURL, Depth: 0}

	endCh := make(chan bool)
	targetCh := make(chan *parser.Target)
	errCh := make(chan error)

	go func() {
		for {
			select {
			case target := <-targetCh:
				log.Info("Found target: ", target.String())
			case err := <-errCh:
				log.Error(err)
			case <-endCh:
				log.Warn("No more targets found")
				break
			}
		}
	}()

	NewScrapper(maxDepth, targetCh, errCh, endCh).Run(rootTarget)
}
