package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/mmcdole/gofeed"
	"github.com/urfave/cli/v2"
)

const (
	DEFAULT_AMOUNT_VIDS   = 2
	DEFAULT_URLS_PATH     = "urls.txt"
	DEFAULT_DOWNLOAD_PATH = "~/Videos"
	CONFIG_DIR            = "~/.config/intentional_youtube"
	CONFIG_FILE           = "config.toml"
)

type Config struct {
	AmountVids   int    `toml:"amount_vids"`
	UrlsPath     string `toml:"urls_path"`
	DownloadPath string `toml:"download_path"`
}

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
		return os.MkdirAll(expandedPath, 0755)
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
	expandedPath, err := expandPath(path)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(expandedPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
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

func loadConfig(configPath string) (Config, error) {
	var config Config
	expandedPath, err := expandPath(configPath)
	if err != nil {
		return config, err
	}
	_, err = toml.DecodeFile(expandedPath, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

func createDefaultConfig(configPath string) error {
	defaultConfig := Config{
		AmountVids:   DEFAULT_AMOUNT_VIDS,
		UrlsPath:     filepath.Join(CONFIG_DIR, DEFAULT_URLS_PATH),
		DownloadPath: DEFAULT_DOWNLOAD_PATH,
	}

	expandedPath, err := expandPath(configPath)
	if err != nil {
		return err
	}

	file, err := os.Create(expandedPath)
	if err != nil {
		return err
	}
	defer file.Close()

	println("No config found, created one at " + CONFIG_DIR)
	println("Please add URLs in " + CONFIG_DIR + "/" + DEFAULT_URLS_PATH)

	return toml.NewEncoder(file).Encode(defaultConfig)
}

func createDefaultURLsFile(filePath string) error {
	// Default content for the URLs file
	defaultContent := `# Default URLs file
# Add your URLs here, one per line. Comments with # and blanklines are allowed.
# Example:
# https://www.youtube.com/feeds/videos.xml?channel_id=UC-lHJZR3Gqxm24_Vd_AJ5Yw`

	// Expand the file path
	expandedPath, err := expandPath(filePath)
	if err != nil {
		return err
	}

	// Create the file with default content
	file, err := os.Create(expandedPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the default content to the file
	_, err = file.WriteString(defaultContent)
	if err != nil {
		return err
	}

	return nil
}

func mainAction(c *cli.Context) error {
	// Expand the config directory path
	configDir, err := expandPath(CONFIG_DIR)
	if err != nil {
		return fmt.Errorf("failed to expand config directory: %w", err)
	}

	// Ensure the config directory exists or create it
	err = ensureDir(configDir)
	if err != nil {
		return fmt.Errorf("failed to ensure config directory: %w", err)
	}

	// Check if the main config file exists
	configPath := filepath.Join(configDir, CONFIG_FILE)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// If the main config file doesn't exist, create it
		err := createDefaultConfig(configPath)
		if err != nil {
			return fmt.Errorf("failed to create default config file: %w", err)
		}

		// Create the default URLs file
		urlsFilePath := filepath.Join(configDir, DEFAULT_URLS_PATH)
		err = createDefaultURLsFile(urlsFilePath)
		if err != nil {
			return fmt.Errorf("failed to create default URLs file: %w", err)
		}
	}

	// Load the config
	config, err := loadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config file: %w", err)
	}

	// Retrieve command line arguments or use values from config
	urlsPath := c.String("urls")
	if urlsPath == "" {
		urlsPath = config.UrlsPath
	}

	numVideos := c.Int("num-videos")
	if numVideos == 0 {
		numVideos = config.AmountVids
	}

	downloadPath := c.String("download-path")
	if downloadPath == "" {
		downloadPath = config.DownloadPath
	}

	// Ensure the download directory exists
	err = ensureDir(downloadPath)
	if err != nil {
		return fmt.Errorf("failed to ensure download directory: %w", err)
	}

	// Parse the URLs file
	urls, err := parseURLFile(urlsPath)
	if err != nil {
		return fmt.Errorf("failed to parse URL file: %w", err)
	}

	// Parse and download feeds
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
