# Web-Crawler
A web crawler written in Go


## usage

```
./web_crawler --domain www.example.com

```

The current functionality is that it will generate a map of resources, keyed by URL. Each resource will have a list of children URLs. These children urls are either valid HTML pages, or static resources such as png, zip, etc.

This data structure represents a graph, where loops are possible

I have not yet implemented a pretty printer that can print this out in a nice way. 


