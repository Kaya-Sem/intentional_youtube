package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mmcdole/gofeed"
	"github.com/urfave/cli/v2"
)

const (
	DEFAULT_AMOUNT_VIDS   = 3
	DEFAULT_URLS_PATH     = "urls.txt"
	DEFAULT_DOWNLOAD_PATH = "~/Videos"
)

// Expand ~ to the home directory if used
func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, path[1:]), nil
	}
	return path, nil
}

// Ensure the directory exists, create it if it doesn't
func ensureDir(path string) error {
	expandedPath, err := expandPath(path)
	if err != nil {
		return err
	}

	if _, err := os.Stat(expandedPath); os.IsNotExist(err) {
		return err
	}
	return nil
}

func downloadVideo(url string, downloadPath string) {
	expandedPath, err := expandPath(downloadPath)
	if err != nil {
		fmt.Printf("Error expanding path %s: %v\n", downloadPath, err)
		return
	}

	// Change the current working directory to the download path
	err = os.Chdir(expandedPath)
	if err != nil {
		fmt.Printf("Error changing directory to %s: %v\n", expandedPath, err)
		return
	}

	cmd := exec.Command("yt-dlp", url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
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

func downloadFeed(feed *gofeed.Feed, numVideos int, downloadPath string) {
	numDownloaded := 0
	for _, entry := range feed.Items {
		url := entry.Link
		if strings.Contains(url, "youtube.com") || strings.Contains(url, "youtu.be") {
			downloadVideo(url, downloadPath)
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

	downloadPath := c.String("download-path")
	if downloadPath == "" {
		downloadPath = DEFAULT_DOWNLOAD_PATH
	}

	err := ensureDir(downloadPath)
	if err != nil {
		return fmt.Errorf("failed to ensure download directory: %w", err)
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
		downloadFeed(feed, numVideos, downloadPath)
	}

	return nil
}

func main() {
	app := &cli.App{
		Name:  "Intentional Youtube",
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
			&cli.StringFlag{
				Name:  "download-path",
				Usage: "Path to the download directory.",
			},
		},
		Action: mainAction,
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("Error running app: %v\n", err)
	}
}
