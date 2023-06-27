// Copyright 2020-22 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package duckduckgosearch

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/proxy"
	"github.com/gocolly/colly/v2/queue"
)

// Result represents a single result from Google Search.
type Result struct {

	// Rank is the order number of the search result.
	Rank int `json:"rank"`

	// URL of result.
	URL string `json:"url"`

	// Title of result.
	Title string `json:"title"`

	// Description of the result.
	Description string `json:"description"`
}

const stdBase = "https://duckduckgo.com/?q="
const defaultAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"

// SearchOptions modifies how the Search function behaves.
type SearchOptions struct {

	// CountryCode sets the ISO 3166-1 alpha-2 code of the localized Google Search homepage to use.
	// The default is "us", which will return results from https://www.google.com.
	CountryCode string

	// LanguageCode sets the language code.
	// Default: en
	LanguageCode string

	// Limit sets how many results to fetch (at maximum).
	Limit int

	// Start sets from what rank the new result set should return.
	Start int

	// UserAgent sets the UserAgent of the http request.
	// Default: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36"
	UserAgent string

	// OverLimit searches for more results than that specified by Limit.
	// It then reduces the returned results to match Limit.
	OverLimit bool

	// ProxyAddr sets a proxy address to avoid IP blocking.
	ProxyAddr string

	// FollowNextPage, when set, scrapes subsequent result pages.
	FollowNextPage bool
}

// Search returns a list of search results from Google.
func Search(ctx context.Context, searchTerm string, opts ...SearchOptions) ([]Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if err := RateLimit.Wait(ctx); err != nil {
		return nil, err
	}

	c := colly.NewCollector(colly.MaxDepth(1))
	if len(opts) == 0 {
		opts = append(opts, SearchOptions{})
	}

	if opts[0].UserAgent == "" {
		c.UserAgent = defaultAgent
	} else {
		c.UserAgent = opts[0].UserAgent
	}

	var lc string
	if opts[0].LanguageCode == "" {
		lc = "en"
	} else {
		lc = opts[0].LanguageCode
	}

	q, _ := queue.New(1, &queue.InMemoryQueueStorage{MaxSize: 10000})

	limit := opts[0].Limit
	if opts[0].OverLimit {
		limit = int(float64(opts[0].Limit) * 1.5)
	}

	results := []Result{}
	nextPageLink := ""
	var rErr error
	filteredRank := 1
	rank := 1

	c.OnRequest(func(r *colly.Request) {
		if err := ctx.Err(); err != nil {
			r.Abort()
			rErr = err
			return
		}
		if opts[0].FollowNextPage && nextPageLink != "" {
			req, err := r.New("GET", nextPageLink, nil)
			if err == nil {
				q.AddRequest(req)
			}
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		rErr = err
	})

	// https://www.w3schools.com/cssref/css_selectors.asp
	c.OnHTML("ol.react-results--main > li > article", func(e *colly.HTMLElement) {

		sel := e.DOM

		linkHref, _ := sel.Find("div > div > a").Attr("href")
		linkText := strings.TrimSpace(linkHref)
		titleText := strings.TrimSpace(sel.Find("div > h2 > a > span").Text())
		descText := strings.TrimSpace(sel.Find("div > div > span").Text())

		rank += 1
		if linkText != "" && linkText != "#" && titleText != "" {
			result := Result{
				Rank:        filteredRank,
				URL:         linkText,
				Title:       titleText,
				Description: descText,
			}
			results = append(results, result)
			filteredRank += 1
		}

		// check if there is a next button at the end.
		// Added this selector as the Id is the same for every language checked on google.com .pt and .es the text changes but the id remains the same
		nextPageHref, _ := sel.Find("a #pnnext").Attr("href")
		nextPageLink = strings.TrimSpace(nextPageHref)

	})

	c.OnHTML("ol.react-results--main > li > article", func(e *colly.HTMLElement) {

		sel := e.DOM

		// check if there is a next button at the end.
		// Added this selector as the Id is the same for every language checked on google.com .pt and .es the text changes but the id remains the same
		if nextPageHref, exists := sel.Attr("href"); exists {
			start := getStart(strings.TrimSpace(nextPageHref))
			nextPageLink = buildUrl(searchTerm, opts[0].CountryCode, lc, limit, start)
			q.AddURL(nextPageLink)
		} else {
			nextPageLink = ""
		}
	})

	url := buildUrl(searchTerm, opts[0].CountryCode, lc, limit, opts[0].Start)

	if opts[0].ProxyAddr != "" {
		rp, err := proxy.RoundRobinProxySwitcher(opts[0].ProxyAddr)
		if err != nil {
			return nil, err
		}
		c.SetProxyFunc(rp)
	}

	q.AddURL(url)
	q.Run(c)

	if rErr != nil {
		if strings.Contains(rErr.Error(), "Too Many Requests") {
			return nil, ErrBlocked
		}
		return nil, rErr
	}

	// Reduce results to max limit
	if opts[0].Limit != 0 && len(results) > opts[0].Limit {
		return results[:opts[0].Limit], nil
	}

	return results, nil
}

func getStart(uri string) int {
	u, err := url.Parse(uri)
	if err != nil {
		fmt.Println(err)
	}
	q := u.Query()
	ss := q.Get("start")
	si, _ := strconv.Atoi(ss)
	return si

}

func base(url string) string {
	if strings.HasPrefix(url, "http") {
		return url
	} else {
		return stdBase + url
	}
}

func buildUrl(searchTerm string, countryCode string, languageCode string, limit int, start int) string {
	searchTerm = strings.Trim(searchTerm, " ")
	searchTerm = strings.Replace(searchTerm, " ", "+", -1)
	countryCode = strings.ToLower(countryCode)

	var url string

	if start == 0 {
		url = fmt.Sprintf("%s%s&hl=%s", stdBase, searchTerm, languageCode)
	} else {
		url = fmt.Sprintf("%s%s&hl=%s&start=%d", stdBase, searchTerm, languageCode, start)
	}

	if limit != 0 {
		url = fmt.Sprintf("%s&num=%d", url, limit)
	}

	return url
}
