package main

import (
	"fmt"
	"sync"
)

type QueueChannel chan string

func Queue(wg *sync.WaitGroup, crawedLinkChannel chan string) (queueCh QueueChannel, closeFn func()) {
	queueCh = make(chan string, 10)
	var mtx sync.RWMutex

	go func() {
		var alreadyCrawed = make(map[string]bool)

		for link := range queueCh {
			mtx.RLock()
			_, ok := alreadyCrawed[link]
			mtx.RUnlock()

			if !ok {
				mtx.Lock()
				alreadyCrawed[link] = true
				mtx.Unlock()

				crawedLinkChannel <- link
			} else {
				wg.Done()
			}
		}

		fmt.Printf("total link: %d\n", len(alreadyCrawed))
	}()

	return queueCh, func() {
		close(queueCh)
	}
}
