package scrapers

import (
	"code.google.com/p/go.net/html"
	"code.google.com/p/go-html-transform/css/selector"
	"code.google.com/p/go-html-transform/h5"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"../models"
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

func companyNameFromJobUrl(url string) string {
	return strings.Split(strings.Split(strings.Split(url, "://")[1], "/")[0], ".")[1]
}

func innerHtml(n *html.Node) string {
	if (n == nil) {
		return ""
	}
	
	result := ""
	if n.Type == html.TextNode {
		result = result + strings.TrimSpace(n.Data)
	}
	
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if n.Type == html.ElementNode {
			childHtml := innerHtml(c)
			if len(childHtml) > 0 {
				result = result + "<"+n.Data+">" + childHtml + "</"+n.Data+">"
			}
		}
	}
	return result
}

func ResumatorJob(url string) (models.Job, error) {
	job := models.Job{SourceUrl: url, SourceName: companyNameFromJobUrl(url)}
	
	resp, err := http.Get(url)
	if err != nil { return job, err }
	
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil { return job, err }
	resp.Body.Close()
	
	titleSelector, err := selector.Selector("title")
	if err != nil { return job, err }
	
	h1Selector, err := selector.Selector("h1")
	if err != nil { return job, err }
	
	descriptionSelector, err := selector.Selector("#resumator-job-description")
	if err != nil { return job, err }
	
	locationSelector, err := selector.Selector("#resumator-job-location")
	if err != nil { return job, err }
	
	howtoapplySelector, err := selector.Selector("#resumator-content-introduction h2.resumator-jobs-text")
	if err != nil { return job, err }
	
	companyUrlSelector, err := selector.Selector("#resumator-company-website a")
	if err != nil { return job, err }
	
	h5content, err := h5.NewFromString(string(content))
	if err != nil { return job, err }
	
	titleNodes := titleSelector.Find(h5content.Top())
	h1Nodes := h1Selector.Find(h5content.Top())
	if (len(h1Nodes) > 0) {
		job.JobTitle = h1Nodes[0].FirstChild.Data
	}
	if (len(titleNodes) > 0) {
		job.CompanyName = strings.TrimSpace(strings.Split(strings.Split(titleNodes[0].FirstChild.Data, job.JobTitle)[1], " - ")[1])
	}
	
	locationNodes := locationSelector.Find(h5content.Top())
	if (len(locationNodes) > 0) {
		job.JobLocation = strings.TrimSpace(locationNodes[0].LastChild.Data)
	}
	
	howtoapplyNodes := howtoapplySelector.Find(h5content.Top())
	if (len(howtoapplyNodes) > 0) {
		job.HowToApply = strings.Replace(innerHtml(howtoapplyNodes[len(howtoapplyNodes)-1]), "h2>", "span>", -1)
	}
	
	companyUrlNodes := companyUrlSelector.Find(h5content.Top())
	if (len(companyUrlNodes) > 0) {
		for j := 0; j < len(companyUrlNodes[0].Attr); j++ {
			if companyUrlNodes[0].Attr[j].Key == "href" {
				job.CompanyUrl = companyUrlNodes[0].Attr[j].Val;
			}
		}
	}
	
	job.JobDescription = innerHtml(descriptionSelector.Find(h5content.Top())[0])
	
	return job, nil
}
