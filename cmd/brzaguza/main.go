package main

import "fmt"
import "github.com/sourcegraph/conc"

import (
	"github.com/tminaorg/brzaguza/search/bing"
	"github.com/tminaorg/brzaguza/search/brave"
	"github.com/tminaorg/brzaguza/search/duckduckgo"
	"github.com/tminaorg/brzaguza/search/google"
	"github.com/tminaorg/brzaguza/search/startpage"
)

func main() {
	var worker conc.WaitGroup
	worker.Go(func() {
		results := googlesearch.Search(nil, "cars for sale in Toronto, Canada")
		fmt.Println("\n\nGoogle:\n")
		fmt.Println(results)
	})
	worker.Go(func() {
		results := startpagesearch.Search(nil, "cars for sale in Toronto, Canada")
		fmt.Println("\n\nStartpage:\n")
		fmt.Println(results)
	})
	worker.Go(func() {
		results := bingsearch.Search(nil, "cars for sale in Toronto, Canada")
		fmt.Println("\n\nBing:\n")
		fmt.Println(results)
	})
	worker.Go(func() {
		results := duckduckgosearch.Search(nil, "cars for sale in Toronto, Canada")
		fmt.Println("\n\nDuckduckgo:\n")
		fmt.Println(results)
	})
	worker.Go(func() {
		results := bravesearch.Search(nil, "cars for sale in Toronto, Canada")
		fmt.Println("\n\nBrave:\n")
		fmt.Println(results)
	})
	worker.Wait()
}
