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
				decodedUrl, err := url.QueryUnescape(strings.Split(strings.Split(fullGoogleUrl, "?q=")[1], "&")[0])
				if err != nil { return []string{}, err }
				urls = append(urls, decodedUrl)
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


func LinkedinJob(url string) (models.Job, error) {
	job := models.Job{SourceUrl: url, SourceName: companyNameFromJobUrl(url)}
	
	resp, err := http.Get(url)
	if err != nil { return job, err }
	
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil { return job, err }
	resp.Body.Close()
	
	
	h1Selector, err := selector.Selector("h1.title")
	if err != nil { return job, err }
	
	descriptionSelector, err := selector.Selector(".description-module div.rich-text")
	if err != nil { return job, err }
	
	locationSelector, err := selector.Selector("h2.sub-header span[itemprop=jobLocation] span")
	if err != nil { return job, err }
	
	howtoapplySelector, err := selector.Selector(".module.highlighted p")
	if err != nil { return job, err }
	
	companyUrlSelector, err := selector.Selector("#top-card .logo-container a")
	if err != nil { return job, err }
	
	companyNameSelector, err := selector.Selector("h2.sub-header span")
	if err != nil { return job, err }
	
	h5content, err := h5.NewFromString(string(content))
	if err != nil { return job, err }
	
	h1Nodes := h1Selector.Find(h5content.Top())
	if (len(h1Nodes) > 0) {
		job.JobTitle = strings.TrimSpace(h1Nodes[0].FirstChild.Data)
	}
	
	companyNameNodes := companyNameSelector.Find(h5content.Top())
	if (len(companyNameNodes) > 0) {
		job.CompanyName = strings.TrimSpace(companyNameNodes[0].FirstChild.Data)
	}
	
	locationNodes := locationSelector.Find(h5content.Top())
	if (len(locationNodes) > 0) {
		job.JobLocation = strings.TrimSpace(locationNodes[0].LastChild.Data)
	}
	
	howtoapplyNodes := howtoapplySelector.Find(h5content.Top())
	if (len(howtoapplyNodes) > 0) {
		job.HowToApply = strings.Replace(innerHtml(howtoapplyNodes[len(howtoapplyNodes)-1]), "h2>", "span>", -1)
	} else {
		job.HowToApply = "Apply On LinkedIn"
	}
	
	companyUrlNodes := companyUrlSelector.Find(h5content.Top())
	if (len(companyUrlNodes) > 0) {
		for j := 0; j < len(companyUrlNodes[0].Attr); j++ {
			if companyUrlNodes[0].Attr[j].Key == "href" {
				job.CompanyUrl = companyUrlNodes[0].Attr[j].Val;
			}
		}
	}
	
	descriptionNodes := descriptionSelector.Find(h5content.Top())
	if (len(descriptionNodes) > 0) {
		job.JobDescription = innerHtml(descriptionNodes[0])
	}
	
	return job, nil
}


func GithubJob(url string) (models.Job, error) {
	job := models.Job{SourceUrl: url, SourceName: companyNameFromJobUrl(url)}
	
	resp, err := http.Get(url)
	if err != nil { return job, err }
	
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil { return job, err }
	resp.Body.Close()
	
	
	h1Selector, err := selector.Selector("h1")
	if err != nil { return job, err }
	
	descriptionSelector, err := selector.Selector("div.column.main")
	if err != nil { return job, err }
	
	locationSelector, err := selector.Selector("p.supertitle")
	if err != nil { return job, err }
	
	howtoapplySelector, err := selector.Selector(".module.highlighted p")
	if err != nil { return job, err }
	
	companyUrlSelector, err := selector.Selector(".column.sidebar .url a")
	if err != nil { return job, err }
	
	companyNameSelector, err := selector.Selector(".column.sidebar h2")
	if err != nil { return job, err }
	
	h5content, err := h5.NewFromString(string(content))
	if err != nil { return job, err }
	
	h1Nodes := h1Selector.Find(h5content.Top())
	if (len(h1Nodes) > 0) {
		job.JobTitle = strings.TrimSpace(h1Nodes[0].FirstChild.Data)
	}
	
	companyNameNodes := companyNameSelector.Find(h5content.Top())
	if (len(companyNameNodes) > 0) {
		job.CompanyName = strings.TrimSpace(companyNameNodes[0].FirstChild.Data)
	}
	
	locationNodes := locationSelector.Find(h5content.Top())
	if (len(locationNodes) > 0) {
		job.JobLocation = strings.TrimSpace(strings.Split(locationNodes[0].LastChild.Data, "/")[1])
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
	
	descriptionNodes := descriptionSelector.Find(h5content.Top())
	if (len(descriptionNodes) > 0) {
		job.JobDescription = innerHtml(descriptionNodes[0])
	}
	
	return job, nil
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
	
	possibleCompanyNameSelector, err := selector.Selector("h2.job_company")
	if err != nil { return job, err }
	
	h5content, err := h5.NewFromString(string(content))
	if err != nil { return job, err }
	
	h1Nodes := h1Selector.Find(h5content.Top())
	if (len(h1Nodes) > 0) {
		job.JobTitle = h1Nodes[0].FirstChild.Data
	}
	
	possibleCompanyNameNodes := possibleCompanyNameSelector.Find(h5content.Top())
	if (len(possibleCompanyNameNodes) > 0) {
		job.CompanyName = possibleCompanyNameNodes[0].FirstChild.Data
	} else {
		titleNodes := titleSelector.Find(h5content.Top())
		if (len(titleNodes) > 0) {
			titlesWithoutJobName := strings.Split(titleNodes[0].FirstChild.Data, job.JobTitle)
			titleWithoutJobName := titlesWithoutJobName[len(titlesWithoutJobName)-1]
			dasherizedTitle := strings.Split(titleWithoutJobName, " - ")
			job.CompanyName = strings.TrimSpace(dasherizedTitle[1])
		}
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
	
	descriptionNodes := descriptionSelector.Find(h5content.Top())
	if (len(descriptionNodes) > 0) {
		job.JobDescription = innerHtml(descriptionNodes[0])
	}
	
	return job, nil
}
