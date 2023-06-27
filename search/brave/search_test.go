// Copyright 2020-21 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package bravesearch_test

import (
	"testing"

	bravesearch "github.com/tminaorg/brzaguza/search/brave"
)

func TestSearch(t *testing.T) {

	q := "Hello World"

	opts := bravesearch.SearchOptions{
		Limit: 20,
	}

	//lint:ignore SA1012 ignore this bare essentials by passing nil for context and removing context package (despite not being idiomatic go).
	returnLinks, err := bravesearch.Search(nil, q, opts)
	if err != nil {
		t.Errorf("something went wrong: %v", err)
		return
	}

	if len(returnLinks) == 0 {
		t.Errorf("no results returned: %v", returnLinks)
	}

	noURL := 0
	noTitle := 0
	noDesc := 0

	for _, res := range returnLinks {
		if res.URL == "" {
			noURL++
		}

		if res.Title == "" {
			noTitle++
		}

		if res.Description == "" {
			noDesc++
		}
	}

	if noURL == len(returnLinks) || noTitle == len(returnLinks) || noDesc == len(returnLinks) {
		t.Errorf("google dom changed")
	}
}
