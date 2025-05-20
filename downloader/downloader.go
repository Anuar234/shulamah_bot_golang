package downloader

import (
	"io"
	"os"
	"path/filepath"
	"regexp"
	"shulamah_bot_golang/lib/e"
	"strings"

	"github.com/kkdai/youtube/v2"
)

func DownloadYoutubeVideo(videoURL, outputDir string) (string, error) {
	video, err := FetchVideoDetails(videoURL)
	if err != nil {
		return "", e.Wrap("failed to fetch video details: %w", err)
	}

	fileName := GenerateFileName(video.Title)
	outputPath := filepath.Join(outputDir, fileName)

	if err := DownloadAndSaveVideo(video, outputPath); err !=nil{
		return "", e.Wrap("failed to download video: %w", err)
	}

	return outputPath, nil
}

func FetchVideoDetails(videoURL string) (*youtube.Video, error) {
	client := youtube.Client{}
	video, err := client.GetVideo(videoURL)
	if err != nil{
		return nil, e.Wrap("error fethcing video details: %w", err)
	}
	return video, nil
}

func DownloadAndSaveVideo(video *youtube.Video, outputPath string) error {
	stream, err := GetBestStream(video)
	if err != nil {
		return err
	}

	return SaveStreamToFile(stream, outputPath)
}

func GetBestStream(video *youtube.Video) (io.Reader, error) {
	client := youtube.Client{}
	formats := video.Formats.WithAudioChannels()

	bestFormat := formats[0]

	for _, f := range formats{
		if f.Bitrate > bestFormat.Bitrate {
			bestFormat = f
		}
	}

	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		return nil, e.Wrap("error retrieving stream: %w", err)
	}

	return stream, nil
}

func SaveStreamToFile(stream io.Reader, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return e.Wrap("error creating file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		return e.Wrap("error writing stream to file: %w", err)
	}

	return nil
}


func RemoveFile(filePath string) error {
	return os.Remove(filePath)
}

func SanitizeFileName(name string) string {
	re := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F]`)
	name = re.ReplaceAllString(name, "_")
	return strings.TrimSpace(name)
}


func GenerateFileName(title string) string {
	return SanitizeFileName(title) + ".mp4"
}