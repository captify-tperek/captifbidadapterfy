package main

import (
	"context"
	"regexp"
	"strings"
)

type LiveClassification struct {
	UrlsWithSegments map[string][]Segment
}

func (l *LiveClassification) classify(ctx context.Context, request BidRequest) ([]Segment, error) {
	//find url in request
	url := removeProtocol(request.Site.Page)
	//find segments for url

	var matchingEntry string
	for key := range l.UrlsWithSegments {
		matched, err := regexp.MatchString(key, url)
		if err != nil {
			continue
		}
		if matched {
			matchingEntry = key
			break
		}
	}

	segments, ok := l.UrlsWithSegments[matchingEntry]
	if !ok {
		return nil, nil
	} else {
		return segments, nil
	}

}
func removeProtocol(url string) string {
	// Find the index of "://"
	index := strings.Index(url, "://")
	if index == -1 {
		return url // No protocol found, return as is
	}
	// Slice the string from the end of "://"
	return url[index+3:]
}
