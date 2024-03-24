package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func Call(url string) (io.ReadCloser, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("fetching:", url)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}
