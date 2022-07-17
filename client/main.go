package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

var netClient = &http.Client{}

func init() {
	tr := &http.Transport{
		MaxIdleConns:        50,
		MaxIdleConnsPerHost: 50,
	}
	netClient = &http.Client{Transport: tr}
}

func testFetch(val chan bool) {
	res, err := netClient.Get("http://localhost:9090")
	if err != nil {
		log.Println(err)
	}
	sitemap, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	if res.Status != "200 OK" {
		log.Println(res.Status)
		log.Println(string(sitemap))
		val <- false
	} else {
		val <- true
	}
	defer res.Body.Close()
}

func main() {
	testCount := 500
	val := make(chan bool)

	for i := 0; i < testCount; i++ {
		go testFetch(val)
	}

	result := []bool{}

	for {
		result = append(result, <-val)
		// fmt.Println(result)
		if len(result) == testCount {
			break
		}
	}

	successCount := 0
	failCount := 0

	for i := 0; i < testCount; i++ {
		if result[i] {
			successCount = successCount + 1
		} else {
			failCount = failCount + 1
		}
	}

	log.Println("success count: ", successCount)
	log.Println("fail count: ", failCount)
}
