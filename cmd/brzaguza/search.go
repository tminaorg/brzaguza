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
	"github.com/tminaorg/brzaguza/search/qwant"
	"github.com/tminaorg/brzaguza/search/startpage"
)

func searchAll(query string) ([]structures.Result) {
	// Make concurrency group and channel for results
	var worker conc.WaitGroup
	resultChannel := make(chan structures.Result)

	// Search Google
	worker.Go(func() {
		results, err := googlesearch.Search(nil, query)
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
				if err != nil || len(results) == 0 {
					if err == nil {
						err = fmt.Errorf("No results found")
					}
					log.Error().
						Err(err).
						Msg("Failed searching Startpage")
				} else {
					for _, r := range results {
						resultChannel <- r
					}
					log.Debug().
						Msg("Finished searching Startpage")
				}
			})

		} else {
			for _, r := range results {
				resultChannel <- r
			}
			log.Debug().
				Msg("Finished searching Google")
		}
	})

	// Search Duckduckgo
	worker.Go(func() {
		results, err := duckduckgosearch.Search(nil, query)
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
				if err != nil || len(results) == 0 {
					if err == nil {
						err = fmt.Errorf("No results found")
					}
					log.Error().
						Err(err).
						Msg("Failed searching Bing")
				} else {
					for _, r := range results {
						resultChannel <- r
					}
					log.Debug().
						Msg("Finished searching Bing")
				}
			})

		} else {
			for _, r := range results {
				resultChannel <- r
			}
			log.Debug().
				Msg("Finished searching Duckduckgo")
		}
	})

	// Search Brave
	worker.Go(func() {
		results, err := bravesearch.Search(nil, query)
		if err != nil || len(results) == 0 {
			if err == nil {
				err = fmt.Errorf("No results found")
			}
			log.Error().
				Err(err).
				Msg("Failed searching Brave")
		} else {
			for _, r := range results {
				resultChannel <- r
			}
			log.Debug().
				Msg("Finished searching Brave")
		}
	})

	// Search Qwant
	worker.Go(func() {
		results, err := qwantsearch.Search(nil, query)
		if err != nil || len(results) == 0 {
			if err == nil {
				err = fmt.Errorf("No results found")
			}
			log.Error().
				Err(err).
				Msg("Failed searching Qwant")
		} else {
			for _, r := range results {
				resultChannel <- r
			}
			log.Debug().
				Msg("Finished searching Qwant")
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