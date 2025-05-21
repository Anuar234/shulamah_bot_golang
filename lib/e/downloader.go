package e

import (
	"errors"
	"os/exec"
	"path/filepath"
)

func DownloadYouTubeVideo(url string) (string, error) {
    output := "downloads/%(title).40s.%(ext)s"
    cmd := exec.Command("yt-dlp", "-f", "best[ext=mp4]", "-o", output, url)

    err := cmd.Run()
    if err != nil {
        return "", err
    }

    files, err := filepath.Glob("downloads/*.mp4")
    if err != nil || len(files) == 0 {
        return "", errors.New("downloaded file not found")
    }

    return files[0], nil
}