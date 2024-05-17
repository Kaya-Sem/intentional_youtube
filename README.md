# Intentional Youtube

This Golang program allows you to download videos from YouTube channels using their RSS feeds. By providing a list of channel URLs, the script fetches the latest videos from each channel and downloads them using the yt-dlp command-line tool.

### Requirements

- [yt-dlp](https://github.com/yt-dlp/yt-dlp)

## Setup

1. Clone or download the repository to your local machine

`git clone git@github.com:Kaya-Sem/intentional_youtube.git`

3. Install **yt-dlp**:

`pip install yt-dlp`

## Usage

1. The binary release or self-build file, should be made executable before use. You can do this in the directory it was downloaded. To use it globally, you also have to put it in your path

`sudo chmod +x intentional_youtube; cp intentional_youtube /usr/local/bin`

1. When first running the program, a config file and template URLs file will be generated at `~/.config/intentional_youtube/`. From then on, running the program will parse any URLS in the file, and download them.

The URLs file may contain blank lines and comments, formatted with `#`.

```markdown
# comment

https://www.youtube.com/feeds/videos.xml?channel_id=UCaiVt4r6YLPzJVgr7pOmD6w

# another comment

https://www.youtube.com/feeds/videos.xml?channel_id=UCXuqSBlHAE6Xw-yeJA0Tunw
```

The default values in `config.toml` can be overriden with flags, which you can find more about with the help flag:

`intentional_youtube --help`

```

```
