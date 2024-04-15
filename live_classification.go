package main

import (
	"context"
	"strings"
)

type LiveClassification struct {
	UrlsWithSegments map[string][]Segment
}

func (l *LiveClassification) classify(ctx context.Context, request BidRequest) ([]Segment, error) {
	//find url in request
	url := removeProtocol(request.Site.Page)
	//find segments for url
	segments, ok := l.UrlsWithSegments[url]
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
