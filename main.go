package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mmcdole/gofeed"
	"github.com/urfave/cli/v2"
)

const (
	DEFAULT_AMOUNT_VIDS = 6
	DEFAULT_URLS_PATH   = "urls.txt"
)

func downloadVideo(url string) {
	cmd := exec.Command("yt-dlp", url)

	// Connect the command's sdout and stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error downloading video %s: %v\n", url, err)
	}
}

func parseURLFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "//") {
			lines = append(lines, line)
		}
	}

	return lines, scanner.Err()
}

func downloadFeed(feed *gofeed.Feed, numVideos int) {
	numDownloaded := 0
	for _, entry := range feed.Items {
		url := entry.Link
		if strings.Contains(url, "youtube.com") || strings.Contains(url, "youtu.be") {
			downloadVideo(url)
			numDownloaded++
			if numDownloaded >= numVideos {
				return
			}
		}
	}
}

func mainAction(c *cli.Context) error {
	urlsPath := c.String("urls")
	if urlsPath == "" {
		urlsPath = DEFAULT_URLS_PATH
	}

	numVideos := c.Int("num-videos")
	if numVideos == 0 {
		numVideos = DEFAULT_AMOUNT_VIDS
	}

	urls, err := parseURLFile(urlsPath)
	if err != nil {
		return fmt.Errorf("failed to parse URL file: %w", err)
	}

	fp := gofeed.NewParser()
	for index, item := range urls {
		fmt.Printf("Parsing URL %d: %s\n", index+1, item)
		feed, err := fp.ParseURL(item)
		if err != nil {
			fmt.Printf("Error parsing feed %s: %v\n", item, err)
			continue
		}
		downloadFeed(feed, numVideos)
	}

	return nil
}

func main() {
	app := &cli.App{
		Name:  "yt-downloader",
		Usage: "Download videos from YouTube RSS feeds.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "urls",
				Usage: "Path to the file containing URLs.",
			},
			&cli.IntFlag{
				Name:  "num-videos",
				Usage: "Number of latest videos to download.",
			},
		},
		Action: mainAction,
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("Error running app: %v\n", err)
	}
}
