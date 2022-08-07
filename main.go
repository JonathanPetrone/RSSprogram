package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/mmcdole/gofeed"
)

type csvData struct {
	Name string
	Link string
}

type bloggerFeed struct {
	FeedTitle   string     `json:"FeedTitle"`
	ItemTitle   string     `json:"ItemTitle"`
	ItemLink    string     `json:"Link"`
	Description string     `json:"Description"`
	Published   *time.Time `json:"Published"`
}

type exportList struct {
	Entry []bloggerFeed
}

type Writer interface {
	Write(p []byte) (n int, err error)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func bloggerListURLs() []string {
	blogger := make(map[string]string)
	bloggerList := []string{}

	// Open the CSV-file
	csvFile, err := os.Open("C:/Users/Jonathan/Documents/GitHub/RSSprogram/feeds.csv")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened CSV file \n")
	defer csvFile.Close()

	// Reading the CSV-file
	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	// For the range of the CSV-file we make a map and then extract URLs from the names
	for _, line := range csvLines {
		csv := csvData{
			Name: line[0],
			Link: line[1],
		}
		//fmt.Println(csv.Name + " " + csv.Link + " ")
		blogger[csv.Name] = csv.Link
		// fmt.Println(blogger[csv.Name])
		bloggerList = append(bloggerList, blogger[csv.Name])
		// fmt.Println(bloggerlist)

		//fmt.Println(blogger)
	}
	return bloggerList
}

func getFeedData(s []string) []bloggerFeed {
	fp := gofeed.NewParser()
	urls := s
	feedlist := []bloggerFeed{}

	for i := range urls {

		fmt.Println(s[i])

		// get current url from range of urls
		feed, _ := fp.ParseURL(s[i])

		// clean tags from description string
		p := bluemonday.StripTagsPolicy()
		clean_desc := p.Sanitize(feed.Items[0].Description)

		// fmt.Println("\nFeed title:", feed.Title, "\nFirst item title:", feed.Items[0].Title, "\nItem link:", feed.Items[0].Link, "\nDescription:", clean_desc, "Published:", feed.Items[0].PublishedParsed, "\n")

		// append each entry per range of urls
		feedlist = append(feedlist, bloggerFeed{FeedTitle: feed.Title, ItemTitle: feed.Items[0].Title, ItemLink: feed.Items[0].Link, Description: clean_desc, Published: feed.Items[0].PublishedParsed})
		//fmt.Println(feedlist)
	}
	return feedlist
}

func sortFunc(b []bloggerFeed) []bloggerFeed {
	bloggers := b

	// prints a single entry of publishtime  we put into the slice of bloggerFeed's (struct)
	fmt.Println(bloggers[1].Published)

	//Sorting the entries based on the Published date
	sort.Slice(bloggers, func(i, j int) bool { return bloggers[i].Published.After(*bloggers[j].Published) })
	//fmt.Println(bloggers)

	return bloggers

	/*
		(for debugging use)
		for x := (len(b) - 1); x > -1; x-- {
			fmt.Println(x)
			fmt.Println(bloggers[x].Published)
		}*/
}

func reduceBloggers(c []bloggerFeed) []bloggerFeed {
	sortedBloggers := c
	numOfBloggers := 8
	reducedBloggers := sortedBloggers[:numOfBloggers]

	/*ar reducedBloggers []bloggerFeed

	for f := 0; f < numOfBloggers; f++ {
		reducedBloggers1 := append(reducedBloggers, sortedBloggers[f])
		fmt.Println(reducedBloggers1)
		fmt.Println(f + 1)
	}

	finalSelection := reducedBloggers
	return finalSelection */
	return reducedBloggers
}

func produceJSON(d []bloggerFeed) {

	reducedBloggers := d
	//encoder := json.NewEncoder(os.Stdout)
	//encodedJSON := encoder.MarshalJson(reducedBloggers)
	u, err := json.Marshal(reducedBloggers)
	fmt.Println(u)

	os.WriteFile("C:/Users/Jonathan/Documents/GitHub/RSSprogram/blogs.json", u, 0644)

	check(err)
}

/*d1 := []byte("hello\ngo\n")
  err := os.WriteFile("/tmp/dat1", d1, 0644)
  check(err)


  f, err := os.Create("/tmp/dat2")
  check(err)*/

//defer f.Close()

// Skriv reducedBloggers till JSON-fil

func main() {
	produceJSON(reduceBloggers(sortFunc(getFeedData(bloggerListURLs()))))
}
