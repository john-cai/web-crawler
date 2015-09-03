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
	LinkTypeHTML           = "html"
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
	Cache   *PageCache
	getter  Getter
	baseUrl string
}

type PageCache struct {
	mutex   *sync.Mutex
	Visited map[string]*resource // a map containing body of pages already visited
}

func NewParser(url string) *parser {
	return &parser{
		Cache:   NewPageCache(),
		getter:  &getter{},
		baseUrl: url,
	}
}
func (p *PageCache) Set(key string, val *resource) {

}

func NewPageCache() *PageCache {
	return &PageCache{
		mutex:   &sync.Mutex{},
		Visited: make(map[string]*resource),
	}
}

// Get all resources in an HTML document. This includes all hyperlinks, js, css resoures
func (p *parser) Parse(url string) {
	fmt.Printf("parsing %s\n", url)
	if _, ok := p.Cache.Visited[url]; !ok {
		p.Cache.mutex.Lock()

		if p.getLinkType(url) != LinkTypeHTML {
			p.Cache.Visited[url] = &resource{Url: url}
			p.Cache.mutex.Unlock()
			return
		}
		resources := p.GetResources(url)

		resource := &resource{
			Url:      url,
			Children: resources,
		}
		p.Cache.Visited[url] = resource
		p.Cache.mutex.Unlock()

		if len(resources) == 0 {
			return
		}

		for _, resource := range resources {
			p.Parse(resource)
		}
	}

}

type getter struct{}

type Getter interface {
	get(url string) (string, error)
}

func (g *getter) get(url string) (string, error) {
	resp, err := http.Get(url)
	//log
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		//handle error
		return "", nil
	}

	page, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return "", err
	}

	return string(page), nil
}

func (p *parser) GetResources(url string) []string {
	links := make([]string, 0)

	page, err := p.getter.get(url)

	if err != nil {
		return links
	}

	linksRegex := regexp.MustCompile(ResourcesRegex)

	allLinks := linksRegex.FindAllStringSubmatch(string(page), -1)

	for _, l := range allLinks {
		newUrl, err := p.getUrlFromLink(l[2])
		if newUrl == "" || newUrl == url {
			continue
		}
		if err != nil {

		}
		links = append(links, newUrl)
	}

	return links
}

func (p *parser) getUrlFromLink(hyperlink string) (string, error) {
	if strings.Contains(hyperlink, "?") {
		hyperlink = strings.Split(hyperlink, "?")[0]
	}

	var buffer bytes.Buffer

	if strings.HasPrefix(hyperlink, "//") {
		return "", nil
	}
	if strings.Index(hyperlink, "/") == 0 {
		buffer.WriteString("http://")
		buffer.WriteString(p.baseUrl)
		buffer.WriteString(hyperlink)
		return buffer.String(), nil
	}

	if strings.Index(hyperlink, p.baseUrl) == -1 {
		return "", nil
	}

	return hyperlink, nil
}

func (p *parser) getLinkType(link string) string {

	extensionSplit := strings.Split(link, ".")

	extension := extensionSplit[len(extensionSplit)-1]
	extRegex := regexp.MustCompile(ResourceExtensionRegex)

	if extRegex.MatchString(link) {

		return extension
	}

	return LinkTypeHTML
}

func (p *parser) PrintChildren() {
	for url, resource := range p.Cache.Visited {
		fmt.Printf("%s has %d children\n", url, len(resource.Children))
	}
}
