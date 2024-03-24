package main

import (
	"sync"
)

var links = []string{
	"https://cmlabs.co",
	"https://sequence.day",
}

func main() {
	var crawedLinkChannel = make(chan string)
	var wg sync.WaitGroup

	queueChannel, queueCloseFunc := Queue(&wg, crawedLinkChannel)
	htmlExporterAddLinkChan, htmlExporterCloseFunc := HtmlExporter()

	for _, link := range links {
		wg.Add(1)
		go func(l string) {
			queueChannel <- l
		}(link)
	}

	// spawn 100 workers
	for i := 0; i < 100; i++ {
		go func() {
			for link := range crawedLinkChannel {
				Craw(&wg, link, queueChannel)
				htmlExporterAddLinkChan <- link
				wg.Done()
			}
		}()
	}

	wg.Wait()
	close(crawedLinkChannel)
	queueCloseFunc()
	htmlExporterCloseFunc()
}
