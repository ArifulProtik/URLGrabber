package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
)

var wg sync.WaitGroup
var mt sync.Mutex

// Final Literation
func main() {
	file, err := os.Open("ip.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		go Grabber(scanner.Text())
		wg.Add(1)

	}
	wg.Wait()

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// stringInArray do If string in list return true false otherwise.
func stringInArray(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Grabber Do Search the bing and collect array of sitelist
func Grabber(ip string) {
	defer wg.Done()
	var output []string
	outfile, err := os.Create("urls.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer outfile.Close()
	if ip == "" {

	}
	page := 1
	for page < 251 {
		client := &http.Client{}
		req, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf(
				"http://www.bing.com/search?q=ip:%s+&count=50&first=1",
				ip,
			),
			nil,
		)
		if err != nil {

		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; rv:57.0) Gecko/20100101 Firefox/57.0")
		res, err := client.Do(req)
		if err != nil {
			fmt.Println("Invalid Request")
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("Couldn't Read")
		}
		re := regexp.MustCompile(`<h2><a href="(.*?)"`)
		links := re.FindAllString(string(body), -1)
		if links != nil {
			for l := range links {
				o := strings.Split(links[l], `"`)
				d := strings.Split(o[1], "/")
				s := d[0] + "//" + d[2]
				if !stringInArray(s, output) {
					output = append(output, s)
				}
			}
		}
		page = page + 50
	}
	for _, links := range output {
		fmt.Println(links)
		fmt.Fprintf(outfile, links+"\n")
	}
}
