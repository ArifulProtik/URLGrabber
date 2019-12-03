package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)
// Test 
func main() {
	result, err := Grabber("103.253.145.35")
	if err != nil {
		fmt.Println(err)
	}
	for  i := range result{
		fmt.Println(result[i])
	}
}
// StringInArray do If string in list return true false otherwise.
func StringInArray(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}
// Grabber Do Search the bing and collect array of sitelist
func Grabber(ip string) (output []string, err error) {
	if ip == "" {
		return output, errors.New("Empty Field Given")
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
			return output, err
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; rv:57.0) Gecko/20100101 Firefox/57.0")
		res, err := client.Do(req)
		if err != nil {
			return output, err
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return output, err
		}
		re := regexp.MustCompile(`<h2><a href="(.*?)"`)
		links := re.FindAllString(string(body), -1)
		if links != nil {
			for l := range links {
				o := strings.Split(links[l],`"`)
				d := strings.Split(o[1],"/")
				s := d[0]+"//"+d[2]
				if !StringInArray(s, output) {
					output = append(output, s)
				}
			}
		}
		page = page + 50
	}
	return output, nil
}