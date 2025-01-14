package searchengine

import (
	"GoMapEnum/src/utils"
	"fmt"
	"regexp"
	"strings"
)

// SEARCH_ENGINE contains url for search on Google and Bing
var SEARCH_ENGINE = map[string]string{"google": `https://www.google.com/search?q=site:linkedin.com/in+"%s"&num=100&start=%d`,
	"bing": `http://www.bing.com/search?q=site:linkedin.com/in+"%s"&first=%d`}

// REGEX_TITLE is the regex to extract title from search engine's results
var REGEX_TITLE = `<h[23](.*?")?>(.*?)<\/h[23]>`

// REGEX_LINKEDIN is the regex to extract field from the title
var REGEX_LINKEDIN = `<h[23](.*?")?>(?P<FirstName>.*?) (?P<LastName>.*?) [-–] (?P<Title>.*?) [-–] (?P<Company>.*?)(\| LinkedIn)(.*?)<\/h[23]>`

// Gather will search a company name and returned the list of people in specified format
func (options *Options) Gather() []string {
	var output []string
	log = options.Log
	// Always insensitive case compare
	options.Company = strings.ToLower(options.Company)
	for searchEngine, formatUrl := range SEARCH_ENGINE {
		log.Target = searchEngine
		log.Verbose("Searching on " + searchEngine + " about " + options.Company)
		url := fmt.Sprintf(formatUrl, options.Company, 0)
		// Get the results of the search
		body, statusCode, err := utils.GetBodyInWebsite(url, options.Proxy, nil)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		// Too many requests. No results returned
		if statusCode == 429 {
			log.Error("Too many requests")
			continue
		}
		reTitle := regexp.MustCompile(REGEX_TITLE)
		reData := regexp.MustCompile(REGEX_LINKEDIN)
		// Extract all links of the body
		links := reTitle.FindAllString(body, -1)

		for _, link := range links {
			// Extract all the title
			result := utils.ReSubMatchMap(reData, link)
			// Compare the company name with case insensitive
			companyName := strings.Trim(strings.ToLower(result["Company"]), " ")
			// Check for the company name (exact match or not)
			if (!options.ExactMatch && strings.Contains(companyName, options.Company)) || (options.ExactMatch && companyName == options.Company) {

				var email string
				email = options.Format
				log.Verbose(result["FirstName"] + " - " + result["LastName"] + " - " + result["Title"] + " - " + result["Company"])
				// output with the specified format
				email = strings.ReplaceAll(email, "{first}", result["FirstName"])
				email = strings.ReplaceAll(email, "{f}", result["FirstName"][0:1])
				email = strings.ReplaceAll(email, "{last}", result["LastName"])
				email = strings.ReplaceAll(email, "{l}", result["LastName"][0:1])
				email = strings.ToLower(email)
				log.Success(email)
				output = append(output, email)
			}
		}
	}
	return output
}
