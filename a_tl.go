package main

import (
	"fmt"
	"net/http"
	"time"
)

type HttpResponse struct {
	url      string
	response *http.Response
	err      error
	req      string
}

func asyncHttpGets(url string) []*HttpResponse {
	ch := make(chan *HttpResponse)
	responses := []*HttpResponse{}
	client := http.Client{}

	letters := "abcdefghijklmnopqrstuvwxyz0123456789"
	for _, l := range letters {
		for _, m := range letters {
			for _, k := range letters {
				go func(l, m, k rune) {
					fmt.Printf("trying %s\n", string(l)+string(m)+string(k))
					// fmt.Printf("Fetching %s \n", url+string(l)+string(m))
					resp, err := client.Get(url + string(l) + string(m) + string(k))
					ch <- &HttpResponse{url, resp, err, string(l) + string(m) + string(k)}
					if err != nil && resp != nil && resp.StatusCode != http.StatusOK {
						fmt.Printf("%s: %d\n", string(l)+string(m), resp.StatusCode)
						defer func() { resp.Body.Close() }()
					} else {
						fmt.Printf("* ")
						resp.Body.Close()
					}
				}(l, m, k)
				time.Sleep(10 * time.Millisecond)
			}
		}
	}

	for {
		select {
		case r := <-ch:
			// fmt.Printf("%s was fetched for %s\n", r.url, r.req)
			if r.err != nil {
				fmt.Println("with an error", r.err)
			}
			if r.response.StatusCode == 404 {
				responses = append(responses, r)
			}
			if r.req == "999" {
				return responses
			}
		case <-time.After(50 * time.Millisecond):
			fmt.Printf(".")
		}
	}
	return responses
}

func main() {
	results := asyncHttpGets("http://github.com/")
	for _, result := range results {
		if result != nil && result.response.StatusCode != 200 {
			fmt.Printf("%s status: %s\n", result.req,
				result.response.Status)
		}
	}
}
