package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type Config struct {
	URLs           []string `json:"urls"`
	TimeoutSeconds int      `json:"timeout_seconds"`
	MaxRetries     int      `json:"max_retries"`
}

func loadConfig() (*Config, error) {
	data, err := os.ReadFile("config.json")
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	return &config, nil
}

func checkWebsite(url string, timeout time.Duration) (time.Duration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	start := time.Now()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	duration := time.Since(start)
	return duration, nil
}

type result struct {
	url      string
	duration time.Duration
	err      error
}

func main() {
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	timeout := time.Duration(config.TimeoutSeconds) * time.Second

	// Print initial "connecting" messages
	for i, url := range config.URLs {
		fmt.Printf("%d. Connecting to %s...\n", i+1, url)
	}

	// Create a map to store results
	resultMap := make(map[string]result)
	var resultMutex sync.Mutex

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Launch a goroutine for each URL
	for _, url := range config.URLs {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			duration, err := checkWebsite(u, timeout)

			resultMutex.Lock()
			resultMap[u] = result{url: u, duration: duration, err: err}
			resultMutex.Unlock()
		}(url)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Move cursor back to the start of the output
	fmt.Print("\033[F\033[K")
	for i := 0; i < len(config.URLs); i++ {
		fmt.Print("\033[F\033[K")
	}

	// Display results in original order
	for i, url := range config.URLs {
		r := resultMap[url]
		if r.err != nil {
			fmt.Printf("%d. Error checking %s: %v\n", i+1, url, r.err)
		} else {
			fmt.Printf("%d. Latency for %s: %v\n", i+1, url, r.duration)
		}
	}

	// Wait for user input before exiting
	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
}
