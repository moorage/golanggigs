package scrapers

import (
	"code.google.com/p/go-html-transform/css/selector"
	"code.google.com/p/go-html-transform/h5"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func Google(query string, site string) ([]string, error) {
	resp, err := http.Get("http://www.google.com/search?um=1&ie=UTF-8&hl=en&tbm=isch&source=og&sa=N&tab=wi&q=" + url.QueryEscape(query) + "+site%3A" + url.QueryEscape(site))
	if err != nil { return []string{}, err }

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil { return []string{}, err }
	resp.Body.Close()
	
	chn, err := selector.Selector("a")
	if err != nil { return []string{}, err }
	
	h5content, err := h5.NewFromString(string(content))
	if err != nil { return []string{}, err }
	
	urls := []string{}
	nodes := chn.Find(h5content.Top())
	for i := 0; i < len(nodes); i++ {
		for j := 0; j < len(nodes[i].Attr); j++ {
			if nodes[i].Attr[j].Key == "href" && strings.Contains(nodes[i].Attr[j].Val, site+"/") {
				fullGoogleUrl := nodes[i].Attr[j].Val
				urls = append(urls, strings.Split(strings.Split(fullGoogleUrl, "?q=")[1], "&")[0])
			}
		}
	}
	return urls, nil
}

func Github(query string) ([]string, error) {
	resp, err := http.Get("https://jobs.github.com/positions?&location=&description=" + url.QueryEscape(query))
	if err != nil { return []string{}, err }

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil { return []string{}, err }
	resp.Body.Close()
	
	chn, err := selector.Selector("a")
	if err != nil { return []string{}, err }
	
	h5content, err := h5.NewFromString(string(content))
	if err != nil { return []string{}, err }
	
	urls := []string{}
	nodes := chn.Find(h5content.Top())
	for i := 0; i < len(nodes); i++ {
		for j := 0; j < len(nodes[i].Attr); j++ {
			if nodes[i].Attr[j].Key == "href" && strings.Contains(nodes[i].Attr[j].Val, "/positions/") {
				urls = append(urls, "https://jobs.github.com" + nodes[i].Attr[j].Val);
			}
		}
	}
	return urls, nil
}

func StackOverflow(query string) ([]string, error) {
	resp, err := http.Get("http://careers.stackoverflow.com/jobs?location=&searchTerm=" + url.QueryEscape(query))
	if err != nil { return []string{}, err }

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil { return []string{}, err }
	resp.Body.Close()
	
	chn, err := selector.Selector("h3 a.job-link")
	if err != nil { return []string{}, err }
	
	h5content, err := h5.NewFromString(string(content))
	if err != nil { return []string{}, err }
	
	urls := []string{}
	nodes := chn.Find(h5content.Top())
	for i := 0; i < len(nodes); i++ {
		for j := 0; j < len(nodes[i].Attr); j++ {
			if nodes[i].Attr[j].Key == "href" {
				urls = append(urls, "http://careers.stackoverflow.com" + nodes[i].Attr[j].Val);
			}
		}
	}
	return urls, nil
}

func Dice(query string) ([]string, error) {
	resp, err := http.Get("http://www.dice.com/job/results?caller=basic&x=all&p=&q=" + url.QueryEscape(query))
	if err != nil { return []string{}, err }

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil { return []string{}, err }
	resp.Body.Close()
	
	chn, err := selector.Selector("a")
	if err != nil { return []string{}, err }
	
	h5content, err := h5.NewFromString(string(content))
	if err != nil { return []string{}, err }
	
	urls := []string{}
	nodes := chn.Find(h5content.Top())
	for i := 0; i < len(nodes); i++ {
		for j := 0; j < len(nodes[i].Attr); j++ {
			if nodes[i].Attr[j].Key == "href" && strings.Contains(nodes[i].Attr[j].Val, "/job/result/") {
				urls = append(urls, "http://www.dice.com" + nodes[i].Attr[j].Val);
			}
		}
	}
	return urls, nil
}
