package main

import (
	"sort"
	"github.com/sourcegraph/conc"
	"github.com/hashicorp/go-set"
	"github.com/rs/zerolog/log"
)

import (
	"github.com/tminaorg/brzaguza/structures"
	"github.com/tminaorg/brzaguza/search/bing"
	"github.com/tminaorg/brzaguza/search/brave"
	"github.com/tminaorg/brzaguza/search/duckduckgo"
	"github.com/tminaorg/brzaguza/search/google"
	"github.com/tminaorg/brzaguza/search/startpage"
)

func searchAll(query string) ([]structures.Result) {
	var worker conc.WaitGroup
	// Make channels for results
	resultChannel := make(chan structures.Result)

	// Search Google
	worker.Go(func() {
		results, _ := googlesearch.Search(nil, query)
		for _, r := range results {
			resultChannel <- r
		}
		log.Debug().
			Msg("Finished searching Google")
	})

	// Search Startpage
	worker.Go(func() {
		results, _ := startpagesearch.Search(nil, query)
		for _, r := range results {
			resultChannel <- r
		}
		log.Debug().
			Msg("Finished searching Startpage")
	})

	// Search Bing
	worker.Go(func() {
		results, _ := bingsearch.Search(nil, query)
		for _, r := range results {
			resultChannel <- r
		}
		log.Debug().
			Msg("Finished searching Bing")
	})

	// Search Duckduckgo
	worker.Go(func() {
		results, _ := duckduckgosearch.Search(nil, query)
		for _, r := range results {
			resultChannel <- r
		}
		log.Debug().
			Msg("Finished searching Duckduckgo")
	})

	// Search Brave
	worker.Go(func() {
		results, _ := bravesearch.Search(nil, query)
		for _, r := range results {
			resultChannel <- r
		}
		log.Debug().
			Msg("Finished searching Brave")
	})

	// Insert results from all searches into HashSet
	resultsHashSet := set.NewHashSet[structures.Result, string](10)
	var helper conc.WaitGroup
	helper.Go(func() {
		for r := range resultChannel {
			resultsHashSet.Insert(r)
		}
	})

	// Wait for results to come back and return them
	worker.Wait()
	results := resultsHashSet.Slice()
	sort.Sort(structures.ByRank(results))
	return results
}