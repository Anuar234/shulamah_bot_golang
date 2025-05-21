package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"shulamah_bot_golang/lib/e"
	"strconv"
	"strings"
	"time"

	"github.com/kkdai/youtube/v2"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

const(
	getUpdatesMethod="getUpdates"
	sendMessageMethod="sendMessage"
)


func New(host string, token string) *Client {
	return &Client{
		host: host,
		basePath: newBasePath(token),
		client: http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}


func (c *Client) Updates(offset int, limit int) (updates []Update, err error){
	defer func() {err = e.WrapIfErr("cant get Updates", err)}()


	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	// do request <- getUpdates

	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil{
		return nil, err
	}

	var res UpdatesResponse

	if err:=json.Unmarshal(data, &res); err!=nil{
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) SendMessage(chatID int, text string) error{
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

_, err := c.doRequest(sendMessageMethod, q)
if err != nil {
	return e.Wrap("cant send message", err)
}

return nil
}

func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {

	defer func() {err = e.WrapIfErr("cant do request", err)}()

	u := url.URL{
		Scheme: "https",
		Host:	c.host,
		Path:	path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err !=nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {_=resp.Body.Close()} ()

	body, err :=io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// Add this constant
const (
    // Your existing constants
    CmdYouTubeDownload = "/ytdl"
)

// Add this to your DoCmd method
func (p *Processor) DoCmd(text string, chatID int, username string) error {
    text = strings.TrimSpace(text)
    
    // Handle commands
    if strings.HasPrefix(text, CmdYouTubeDownload) {
        return p.handleYouTubeDownload(chatID, text)
    }
    
    // Your existing command handling
    return nil
}

// This function gets the Telegram bot token through the bot API
func (p *Processor) getBotToken() (string, error) {
    // Use reflection or another method to get the token
    // As a fallback, you can read it from a config file or environment variable
    // For now, this is a placeholder
    return "", fmt.Errorf("implement this method to get the bot token")
}

// Function to validate YouTube URLs
func isValidYouTubeURL(url string) bool {
    patterns := []string{
        `^(https?\:\/\/)?(www\.)?(youtube\.com|youtu\.?be)\/.+$`,
        `^(https?\:\/\/)?(www\.)?youtube\.com\/watch\?v=([a-zA-Z0-9_-]{11})`,
        `^(https?\:\/\/)?(www\.)?youtu\.be\/([a-zA-Z0-9_-]{11})`,
    }
    
    for _, pattern := range patterns {
        match, _ := regexp.MatchString(pattern, url)
        if match {
            return true
        }
    }
    return false
}

// Function to download a YouTube video
func downloadYouTubeVideo(videoURL string) (string, error) {
    // Create YouTube client
    client := youtube.Client{}
    
    // Parse video ID from URL
    video, err := client.GetVideo(videoURL)
    if err != nil {
        return "", fmt.Errorf("failed to get video info: %w", err)
    }
    
    // Get formats with audio
    formats := video.Formats.WithAudioChannels()
    if len(formats) == 0 {
        return "", fmt.Errorf("no suitable video format found")
    }
    
    // Choose a format - preferably mp4 with audio
    var selectedFormat *youtube.Format
    for _, format := range formats {
        if format.MimeType == "video/mp4" {
            selectedFormat = &format
            break
        }
    }
    
    if selectedFormat == nil {
        // Fall back to the first format with audio
        selectedFormat = &formats[0]
    }
    
    // Create download directory
    downloadDir := "downloads"
    if err := os.MkdirAll(downloadDir, 0755); err != nil {
        return "", fmt.Errorf("failed to create download directory: %w", err)
    }
    
    // Create output file
    fileName := fmt.Sprintf("%s-%d.mp4", video.ID, time.Now().Unix())
    outputPath := filepath.Join(downloadDir, fileName)
    outputFile, err := os.Create(outputPath)
    if err != nil {
        return "", fmt.Errorf("failed to create output file: %w", err)
    }
    defer outputFile.Close()
    
    // Download video content
    stream, size, err := client.GetStream(video, selectedFormat)
    if err != nil {
        return "", fmt.Errorf("failed to get video stream: %w", err)
    }
    defer stream.Close()
    
    // Check if the file is too large (Telegram has a 50MB limit for bots)
    if size > 50*1024*1024 {
        return "", fmt.Errorf("video is too large (%.2f MB). Telegram bots can only send files up to 50MB", float64(size)/(1024*1024))
    }
    
    // Copy content to output file
    _, err = io.Copy(outputFile, stream)
    if err != nil {
        return "", fmt.Errorf("failed to save video: %w", err)
    }
    
    // Return the path to the downloaded file
    return outputPath, nil
}

// Method to handle YouTube download commands
func (p *Processor) handleYouTubeDownload(chatID int, text string) error {
    // Extract URL from command
    args := strings.TrimSpace(strings.TrimPrefix(text, CmdYouTubeDownload))
    if args == "" {
        return p.tg.SendMessage(chatID, "Please provide a YouTube video URL. Usage: /ytdl [youtube-url]")
    }
    
    // Validate URL
    if !isValidYouTubeURL(args) {
        return p.tg.SendMessage(chatID, "Invalid YouTube URL. Please provide a valid YouTube video link.")
    }
    
    // Inform user download is starting
    if err := p.tg.SendMessage(chatID, "Processing your YouTube video. This may take a moment..."); err != nil {
        return err
    }
    
    // Download the video
    filePath, err := downloadYouTubeVideo(args)
    if err != nil {
        errMsg := fmt.Sprintf("Failed to download video: %s", err.Error())
        return p.tg.SendMessage(chatID, errMsg)
    }
    defer os.Remove(filePath)  // Clean up the file when done
    
    // Inform user download is complete
    if err := p.tg.SendMessage(chatID, "Video download complete! Sending to you now..."); err != nil {
        return err
    }
}