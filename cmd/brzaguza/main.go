package main

import "fmt"
import "github.com/tminaorg/brzaguza/search/bing"
import "github.com/tminaorg/brzaguza/search/brave"
import "github.com/tminaorg/brzaguza/search/duckduckgo"
import "github.com/tminaorg/brzaguza/search/google"
import "github.com/tminaorg/brzaguza/search/startpage"

func main() {
	fmt.Println("\n\nGoogle:\n")
	fmt.Println(googlesearch.Search(nil, "cars for sale in Toronto, Canada"))
	fmt.Println("\n\nStartpage:\n")
	fmt.Println(startpagesearch.Search(nil, "cars for sale in Toronto, Canada"))
	fmt.Println("\n\nBing:\n")
	fmt.Println(bingsearch.Search(nil, "cars for sale in Toronto, Canada"))
	fmt.Println("\n\nDuckduckgo:\n")
	fmt.Println(duckduckgosearch.Search(nil, "cars for sale in Toronto, Canada"))
	fmt.Println("\n\nBrave:\n")
	fmt.Println(bravesearch.Search(nil, "cars for sale in Toronto, Canada"))
}
