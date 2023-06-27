package main

import "fmt"
import "github.com/rocketlaunchr/google-search"

func main() {
	fmt.Println(googlesearch.Search(nil, "cars for sale in Toronto, Canada"))
}