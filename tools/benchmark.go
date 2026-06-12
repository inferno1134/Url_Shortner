package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	ShortURL string `json:"short_url"`
}

var concurrencyLevels = []int{100, 500, 1000, 5000, 10000}

func main() {
	baseURL := flag.String("base", "http://localhost:8080", "app base URL")
	longURL := flag.String("url", "https://example.com", "long url to shorten")
	total := flag.Int("total", 10000, "total requests per concurrency level")
	flag.Parse()

	client := &http.Client{Timeout: 20 * time.Second}

	shortCode, err := createShortURL(client, *baseURL, *longURL)
	if err != nil {
		panic(err)
	}

	targetURL := fmt.Sprintf("%s/%s", *baseURL, shortCode)
	fmt.Println("Benchmark Target:", targetURL)
	fmt.Println()

	for _, conc := range concurrencyLevels {
		if conc > *total {
			fmt.Printf("skipping concurrency %d because total requests is %d\n", conc, *total)
			continue
		}

		fmt.Printf("=== concurrency %d, total requests %d ===\n", conc, *total)
		success, failed, avg, p50, p95, p99, rps := benchmarkGet(client, targetURL, *total, conc)
		fmt.Printf("success=%d failed=%d rps=%.1f avg=%.2fms p50=%.2fms p95=%.2fms p99=%.2fms\n\n",
			success, failed,
			rps,
			float64(avg.Microseconds())/1000.0,
			float64(p50.Microseconds())/1000.0,
			float64(p95.Microseconds())/1000.0,
			float64(p99.Microseconds())/1000.0,
		)
	}
}

func createShortURL(client *http.Client, baseURL, longURL string) (string, error) {
	body, err := json.Marshal(shortenRequest{URL: longURL})
	if err != nil {
		return "", err
	}

	resp, err := client.Post(baseURL+"/shorten", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("create short URL failed: status=%d body=%s", resp.StatusCode, string(data))
	}

	var result shortenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	parts := strings.Split(result.ShortURL, "/")
	return parts[len(parts)-1], nil
}

func benchmarkGet(client *http.Client, url string, total, concurrency int) (int, int, time.Duration, time.Duration, time.Duration, time.Duration, float64) {
	var success uint64
	var failed uint64
	var mu sync.Mutex
	durations := make([]time.Duration, 0, total)

	jobs := make(chan int, total)
	var wg sync.WaitGroup

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range jobs {
				req, err := http.NewRequest("GET", url, nil)
				if err != nil {
					atomic.AddUint64(&failed, 1)
					continue
				}

				start := time.Now()
				resp, err := client.Do(req)
				dur := time.Since(start)

				mu.Lock()
				durations = append(durations, dur)
				mu.Unlock()

				if err != nil {
					atomic.AddUint64(&failed, 1)
					continue
				}

				if resp.StatusCode >= 200 && resp.StatusCode < 400 {
					atomic.AddUint64(&success, 1)
				} else {
					atomic.AddUint64(&failed, 1)
				}
				io.Copy(io.Discard, resp.Body)
				resp.Body.Close()
			}
		}()
	}

	for i := 0; i < total; i++ {
		jobs <- i
	}
	close(jobs)

	start := time.Now()
	wg.Wait()
	elapsed := time.Since(start)

	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
	sum := time.Duration(0)
	for _, d := range durations {
		sum += d
	}

	avg := time.Duration(0)
	if len(durations) > 0 {
		avg = sum / time.Duration(len(durations))
	}

	return int(success), int(failed), avg, percentile(durations, 50), percentile(durations, 95), percentile(durations, 99), float64(total) / elapsed.Seconds()
}

func percentile(durations []time.Duration, p int) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	index := (len(durations)*p+99)/100 - 1
	if index < 0 {
		index = 0
	}
	if index >= len(durations) {
		index = len(durations) - 1
	}
	return durations[index]
}
