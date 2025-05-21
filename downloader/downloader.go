package downloader

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/kkdai/youtube/v2"
)

func DownloadVideo(videoURL string, downloadDir string) (string, error) {
	client := youtube.Client{}

	video, err := client.GetVideo(videoURL)
	if err != nil {
		return "", fmt.Errorf("failed to get video: %w", err)
	}

	stream, _, err := client.GetStream(video, &video.Formats.WithAudioChannels()[0])
	if err != nil {
		return "", fmt.Errorf("failed to get stream: %w", err)
	}

	filename := sanitizeFilename(video.Title) + ".mp4"
	filePath := filepath.Join(downloadDir, filename)

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		return "", fmt.Errorf("failed to download video: %w", err)
	}

	return filePath, nil
}

func sanitizeFilename(name string) string {
	// Убираем все символы, которые могут вызвать проблемы в путях файлов
	return strings.Map(func(r rune) rune {
		if strings.ContainsRune("\\/:*?\"<>|", r) {
			return -1
		}
		return r
	}, name)
}
