package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUrlFromLink(t *testing.T) {
	baseUrl := "www.example.com"
	tests := []struct {
		Name     string
		Url      string
		Expected string
	}{
		{
			Name:     "Valid HTML link",
			Url:      fmt.Sprintf("http://%s?whatever", baseUrl),
			Expected: fmt.Sprintf("http://%s", baseUrl),
		},
		{
			Name:     "Valid internal link",
			Url:      "/some/internal/resource",
			Expected: fmt.Sprintf("http://%s/some/internal/resource", baseUrl),
		},
		{
			Name:     "external resource",
			Url:      "http://someotherdomain.com",
			Expected: "",
		},
		{
			Name:     "css resource",
			Url:      "//some/other/resource",
			Expected: "",
		},
	}
	parser := NewParser(baseUrl)

	for _, test := range tests {

		result, err := parser.getUrlFromLink(test.Url)
		assert.NoError(t, err, test.Name)
		assert.Equal(t, test.Expected, result, test.Name)
	}
}

func TestGetLinkType(t *testing.T) {

	tests := []struct {
		Name             string
		Link             string
		ExpectedLinkType string
	}{
		{
			Name:             "Valid HTML link",
			Link:             "http://thisis.some.valid.url",
			ExpectedLinkType: LinkTypeHTML,
		},
		{
			Name:             "Valid HTML link",
			Link:             "http://thisis.some.valid.url/",
			ExpectedLinkType: LinkTypeHTML,
		},
		{
			Name:             "js resource",
			Link:             "http://somejavascriptresource.com/whatever.js",
			ExpectedLinkType: "js",
		},
		{
			Name:             "css resource",
			Link:             "/something/style.css",
			ExpectedLinkType: "css",
		},
		{
			Name:             "some non html resource",
			Link:             "/something/style.zzz",
			ExpectedLinkType: "zzz",
		},
	}
	parser := NewParser("")

	for _, test := range tests {

		linkType := parser.getLinkType(test.Link)
		assert.Equal(t, linkType, test.ExpectedLinkType)
	}

}

// Mock getter

func newMockGetter() *mockGetter {
	pageIndex := map[string]string{
		"http://www.example.com": `<html>
<body>
<a href="/a/">
<a href="/b/">
<a href="/c/">
<a href="//googletagmanager">
<img src="/assets/images/b.gif" alt="">
<img src="http://www.google.com/assets/images/b.gif" alt="">
</body>
</html>	
`,
		"http://www.example.com/a/": `<html>
<body>
<a href="/d">
</body>
</html>	
`,
		"http://www.example.com/b/": `<html>
<body>
some thing
</body>
</html>	
`, "http://www.example.com/c/": `<html>
<body>
<img src="/assets/images/b.gif" alt="">
</body>
</html>	
`,
	}
	return &mockGetter{
		pageIndex: pageIndex,
	}
}

type mockGetter struct {
	pageIndex map[string]string
}

func (m *mockGetter) get(url string) (string, error) {
	if body, ok := m.pageIndex[url]; ok {
		return body, nil
	}
	return "", fmt.Errorf("page not found")
}

func TestGetResources(t *testing.T) {
	baseUrl := "http://www.example.com"
	p := &parser{
		Cache:   NewPageCache(),
		getter:  newMockGetter(),
		baseUrl: baseUrl,
	}

	resources := p.GetResources(baseUrl)

	assert.Equal(t, len(resources), 4)
}

func TestParser(t *testing.T) {
	baseUrl := "www.example.com"
	p := &parser{
		Cache:   NewPageCache(),
		getter:  newMockGetter(),
		baseUrl: baseUrl,
	}
	p.Parse(fmt.Sprintf("http://%s", baseUrl))

	assert.Equal(t, 6, len(p.Cache.Visited))
}
