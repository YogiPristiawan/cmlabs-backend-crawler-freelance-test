package main

import (
	"log"
	"net/url"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

func Craw(wg *sync.WaitGroup, link string, queueChannel chan string) {
	body, err := Call(link)
	if err != nil {
		log.Println("failed to fetch", link, err)
		return
	}
	defer body.Close()

	tokenizer := html.NewTokenizer(body)

	for tokenizer.Next() != html.ErrorToken {

		token := tokenizer.Token()

		if isAnchor(token) {
			attributes := token.Attr
			for _, attr := range attributes {
				if attr.Key == "href" {
					sanitizedUrl, valid := sanitizeURL(attr.Val, link)
					if valid {
						wg.Add(1)
						queueChannel <- sanitizedUrl
					}
				}
			}
		}
	}
}

func isAnchor(token html.Token) bool {
	return token.DataAtom.String() == "a"
}

func sanitizeURL(nextTarget string, currentTarget string) (sanitized string, valid bool) {
	nt, err := url.ParseRequestURI(nextTarget)
	if err != nil {
		return "", false
	}

	ct, err := url.ParseRequestURI(currentTarget)
	if err != nil {
		return "", false
	}

	if nt.Scheme == "" && nt.Host == "" && strings.HasPrefix(nt.Path, "/") {
		nt.Scheme = ct.Scheme
		nt.Host = ct.Host
	}

	for _, link := range links {
		l, _ := url.ParseRequestURI(link)

		s := strings.HasSuffix(nt.Host, l.Host)

		if l.Scheme == nt.Scheme && l.Host == nt.Host && s {
			return nt.String(), true
		}
	}

	return "", false
}
