package parser

import (
	"io"
	"net/url"

	"fmt"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/html"
)

var attributeHREF = "href"
var tagNameLink = byte('a')

type Target struct {
	URL   url.URL
	Depth int
}

func (t *Target) String() string {
	return fmt.Sprintf("url:%s depth:%d", t.URL.String(), t.Depth)
}

func (t Target) ParseLinks(text io.ReadCloser) []Target {
	if text == nil {
		return []Target{}
	}
	defer text.Close()
	// using a map to ignore duplicates
	links := map[string]int{}

	tokenizer := html.NewTokenizer(text)
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			if tokenizer.Err() == io.EOF {
				// no more tags
				log.Debug("Reached EOF")
				break
			}
			// failed to process token. skip this one
			log.Warn("Failed to process token")
			continue
		}

		// current token is valid
		if tokenIsLink(tokenizer) && tokenType == html.StartTagToken {
			log.WithFields(log.Fields{"raw-tag": string(tokenizer.Raw())}).Debug("Found Link tag")
			var keyB, valB []byte
			for more := true; more; {
				keyB, valB, more = tokenizer.TagAttr()
				key, val := string(keyB), string(valB)
				log.WithFields(log.Fields{"key": key, "val": val}).Debug("Got attribute")
				if key == attributeHREF {
					links[val]++
					log.WithFields(log.Fields{attributeHREF: val}).Debug("Got reference")
				}
			}
		}

	}

	targets := make([]Target, 0, len(links))
	for link := range links {
		parsed, err := url.Parse(link)
		if err != nil {
			log.WithFields(log.Fields{"link": link}).Error(err)
		}
		parsed.IsAbs()
		parsed.Scheme = t.URL.Scheme
		parsed.Host = t.URL.Host
		targets = append(targets, Target{URL: *parsed, Depth: t.Depth + 1})
	}
	return targets
}

func tokenIsLink(t *html.Tokenizer) bool {
	tagName, _ := t.TagName()
	return len(tagName) == 1 && tagName[0] == tagNameLink
}
