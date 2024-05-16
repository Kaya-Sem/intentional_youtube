# ⏯️ YouTube Channel Feed Downloader

This Golang program allows you to download videos from YouTube channels using their RSS feeds. By providing a list of channel URLs, the script fetches the latest videos from each channel and downloads them using the yt-dlp command-line tool.

### Requirements

- [yt-dlp](https://github.com/yt-dlp/yt-dlp)

## Setup

1. Clone or download the repository to your local machine

`git clone git@github.com:Kaya-Sem/intentional_youtube.git`

3. Install **yt-dlp**:

`pip install yt-dlp`

## Usage

1. Create a text file containing a list of URLs, each representing a YouTube channel feed. Each URL should be on a separate line and should include the channel ID appended to the base URL. The file may contain blank lines and comments, formatted with `//`

```markdown
// comment
https://www.youtube.com/feeds/videos.xml?channel_id=UCaiVt4r6YLPzJVgr7pOmD6w

// another comment
https://www.youtube.com/feeds/videos.xml?channel_id=UCXuqSBlHAE6Xw-yeJA0Tunw
```

2. Run the script and provide the path to the text file as a command-line argument:

```
TODO
```

This will fetch the latest videos from each channel in the URL file and download them using `yt-dlp`
