package main

import (
	"fmt"
	"os"
	"sync"
)

var head = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Document</title>
</head>
<body>
<ul>
`

var foot = `</ul></body></html>`

func HtmlExporter() (addLinkChannel chan string, closeFunc func()) {
	addLinkChannel = make(chan string, 10)

	go func() {
		var mtx sync.Mutex

		file, err := os.Create("index.html")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		file.WriteString(head)

		for l := range addLinkChannel {
			var a = fmt.Sprintf(`<li><a href="%s">%s</a></li>`, l, l)

			mtx.Lock()

			file.WriteString(a)

			mtx.Unlock()
		}

		_, err = file.WriteString(foot)
		if err != nil {
			panic(err)
		}
	}()

	return addLinkChannel, func() {
		close(addLinkChannel)
	}
}
