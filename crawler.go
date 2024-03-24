package main

import (
	"log"
	"net/url"
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
					if isValidURL(attr.Val) {
						wg.Add(1)
						queueChannel <- attr.Val
					}
				}
			}
		}
	}
}

func isAnchor(token html.Token) bool {
	return token.DataAtom.String() == "a"
}

func isValidURL(str string) bool {
	u, err := url.ParseRequestURI(str)

	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	for _, target := range links {
		u2, _ := url.ParseRequestURI(target)
		if u2.Scheme == u.Scheme && u2.Host == u.Host {
			return true
		}
	}

	return false
}
