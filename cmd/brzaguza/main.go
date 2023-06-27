package main

import "fmt"
import "github.com/rocketlaunchr/google-search"
import "github.com/tminaorg/brzaguza/search/bing"
import "github.com/tminaorg/brzaguza/search/brave"
import "github.com/tminaorg/brzaguza/search/duckduckgo"
import "github.com/tminaorg/brzaguza/search/startpage"

func main() {
	fmt.Println("Google")
	fmt.Println(googlesearch.Search(nil, "cars for sale in Toronto, Canada"))
	fmt.Println("Startpage")
	fmt.Println(startpagesearch.Search(nil, "cars for sale in Toronto, Canada"))
	fmt.Println("Bing")
	fmt.Println(bingsearch.Search(nil, "cars for sale in Toronto, Canada"))
	fmt.Println("Duckduckgo")
	fmt.Println(duckduckgosearch.Search(nil, "cars for sale in Toronto, Canada"))
	fmt.Println("Brave")
	fmt.Println(bravesearch.Search(nil, "cars for sale in Toronto, Canada"))
}