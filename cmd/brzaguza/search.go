package main

import (
	"fmt"
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
		results, err := googlesearch.Search(nil, query)
		for _, r := range results {
			resultChannel <- r
		}
		if err != nil || len(results) == 0 {
			if err == nil {
				err = fmt.Errorf("No results found")
			}
			log.Error().
				Err(err).
				Msg("Failed searching Google, falling back to Startpage")
			
			// Search Startpage because Google failed
			worker.Go(func() {
				results, err := startpagesearch.Search(nil, query)
				for _, r := range results {
					resultChannel <- r
				}
				if err != nil {
					log.Error().
						Err(err).
						Msg("Failed searching Startpage")
				} else {
					log.Debug().
						Msg("Finished searching Startpage")
				}
			})

		} else {
			log.Debug().
				Msg("Finished searching Google")
		}
	})

	// Search Duckduckgo
	worker.Go(func() {
		results, err := duckduckgosearch.Search(nil, query)
		for _, r := range results {
			resultChannel <- r
		}
		if err != nil || len(results) == 0 {
			if err == nil {
				err = fmt.Errorf("No results found")
			}
			log.Error().
				Err(err).
				Msg("Failed searching Duckduckgo, falling back to Bing")
			
			// Search Bing because Duckduckgo failed
			worker.Go(func() {
				results, err := bingsearch.Search(nil, query)
				for _, r := range results {
					resultChannel <- r
				}
				if err != nil {
					log.Error().
						Err(err).
						Msg("Failed searching Bing")
				} else {
					log.Debug().
						Msg("Finished searching Bing")
				}
			})

		} else {
			log.Debug().
				Msg("Finished searching Duckduckgo")
		}
	})

	// Search Brave
	worker.Go(func() {
		results, err := bravesearch.Search(nil, query)
		for _, r := range results {
			resultChannel <- r
		}
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed searching Brave")
		} else {
			log.Debug().
				Msg("Finished searching Brave")
		}
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