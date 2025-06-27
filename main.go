package main

import (
	"bufio"
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

	// Print initial "testing" messages
	for i, url := range config.URLs {
		fmt.Printf("%d. Testing %s... ", i+1, url)
		fmt.Println() // 换行，为结果预留位置
	}

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup
	var outputMutex sync.Mutex

	// Launch a goroutine for each URL
	for i, url := range config.URLs {
		wg.Add(1)
		go func(index int, u string) {
			defer wg.Done()
			duration, err := checkWebsite(u, timeout)

			// Lock output to prevent interference
			outputMutex.Lock()
			defer outputMutex.Unlock()

			// Move cursor to the specific line and update it
			// ANSI escape sequence to move cursor up and to beginning of line
			linesToMoveUp := len(config.URLs) - index
			fmt.Printf("\033[%dA", linesToMoveUp) // Move cursor up
			fmt.Printf("\033[2K")                 // Clear entire line
			fmt.Printf("\r")                      // Move cursor to beginning of line

			if err != nil {
				fmt.Printf("%d. Testing %s... ERROR: %v", index+1, u, err)
			} else {
				fmt.Printf("%d. Testing %s... %v", index+1, u, duration)
			}

			// Move cursor back down to the bottom
			fmt.Printf("\033[%dB", linesToMoveUp)
		}(i, url)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Move cursor to the end and ensure clean output
	fmt.Println()

	// Wait for user input before exiting
	fmt.Print("Press Enter to exit...")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadLine()
}
