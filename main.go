package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

const (
	ResourcesRegex         = `(href|src)="([^"]+)"`
	ResourceExtensionRegex = `^.*\.(jpg|JPG|gif|GIF|doc|DOC|pdf|PDF|zip|png|svg|css|js|eps)$`
	DigitalOceanBaseUrl    = "www.digitalocean.com"
)

type Resource interface {
	ParseLinks() error
}

type resource struct {
	Url      string
	Type     string
	Children []string
	Parent   *resource
	Body     string
	Parser   *parser
	Cache    *PageCache
}

type parser struct {
	RootUrl string
}

type PageCache struct {
	mutex   *sync.Mutex
	Visited map[string]*resource // a map containing body of pages already visited
}

func NewPageCache() *PageCache {
	return &PageCache{
		mutex:   &sync.Mutex{},
		Visited: make(map[string]*resource),
	}
}

// Get all resources in an HTML document. This includes all hyperlinks, js, css resoures
func (r *resource) Parse(url string) {

	if _, ok := r.Cache.Visited[url]; !ok {
		//cache hit
		r.Cache.mutex.Lock()

		if r.getLinkType(url) != "html" {
			r.Cache.Visited[url] = &resource{Url: url}
			r.Cache.mutex.Unlock()
			return
		}
		resources := r.GetResources(url)

		resource := &resource{
			Url:      url,
			Children: resources,
		}
		r.Cache.Visited[url] = resource
		r.Cache.mutex.Unlock()

		if len(resources) == 0 {
			return
		}

		for _, resource := range resources {
			if _, ok := r.Cache.Visited[resource]; !ok {
				/*if strings.HasPrefix(resource, "http://www.digitalocean.com/company") || strings.HasPrefix(resource, "http://www.digitalocean.com/community") || strings.HasPrefix(resource, "https://www.digitalocean.com/community") {
					continue
				}*/
				r.Parse(resource)
			}
		}
	} else {
		fmt.Println("cache hit")
	}

}

func (r *resource) GetResources(url string) []string {
	fmt.Printf("getting resources for %s\n", url)
	links := make([]string, 0)

	resp, err := http.Get(url)
	//log
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		fmt.Printf("resource: %+v\n", r)
		//handle error
	}

	p, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		fmt.Printf("resource: %+v\n", r)
	}

	//log

	if err != nil {
		//handle error
	}
	page := string(p)

	linksRegex := regexp.MustCompile(ResourcesRegex)

	allLinks := linksRegex.FindAllStringSubmatch(page, -1)

	for _, l := range allLinks {
		newUrl, err := r.getUrlFromLink(l[2])
		if newUrl == "" || newUrl == url {
			continue
		}
		if err != nil {

		}
		links = append(links, newUrl)
	}

	return links
}

func (r *resource) getUrlFromLink(hyperlink string) (string, error) {
	if strings.Contains(hyperlink, "?") {
		hyperlink = strings.Split(hyperlink, "?")[0]
	}

	var buffer bytes.Buffer

	if strings.Contains(hyperlink, "googletagmanager.com") {
		return "", nil
	}
	if strings.Index(hyperlink, "/") == 0 {
		buffer.WriteString("http://")
		buffer.WriteString(DigitalOceanBaseUrl)
		buffer.WriteString(hyperlink)
		return buffer.String(), nil
	}

	if strings.Index(hyperlink, DigitalOceanBaseUrl) == -1 || strings.Contains(hyperlink, "googletagmanager.com") {
		return "", nil
	}

	return hyperlink, nil
}

func (r *resource) getLinkType(link string) string {

	extensionSplit := strings.Split(link, ".")

	extension := extensionSplit[len(extensionSplit)-1]
	extRegex := regexp.MustCompile(ResourceExtensionRegex)

	if extRegex.MatchString(link) {

		return extension
	}

	return "html"
}

/*
func (r *resource) PrettyPrint(indentation int) {
	var b bytes.Buffer

	for i := 0; i < indentation; i++ {
		b.WriteString(" ")
	}

	for _, resource := range r.Children {
		fmt.Printf("%s%s\n", b.String(), resource)
		resource.PrettyPrint(indentation + 1)
	}

}
*/

func (r *resource) PrintChildren() {
	for url, resource := range r.Cache.Visited {
		fmt.Printf("%s has %d children\n", url, len(resource.Children))
	}
}
func main() {
	rootResource := &resource{
		Url:   "https://www.digitalocean.com",
		Type:  "html",
		Cache: NewPageCache(),
	}
	rootResource.Parse(rootResource.Url)
	rootResource.PrintChildren()
}
