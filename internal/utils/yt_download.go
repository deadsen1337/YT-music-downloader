package utils

import (
	"errors"
	"fmt"
	"github.com/kkdai/youtube/v2"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

var (
	ytregex = regexp.MustCompile(`(http:|https:)?\/\/(www\.)?(youtube.com|youtu.be)\/(watch)?(\?v=)?(\S+)?`)
	Version = "0.0.9"
)

func searchffmpeg() {
	path, err := exec.LookPath("ffmpeg")
	if err != nil {
		fmt.Println("ffmpeg not found", path)
	}
}

func YouTubeDownload(link string) (string, error) {
	searchffmpeg()
	if !ytregex.MatchString(link) {
		return "", errors.New("bad link: " + link)
	}
	videoID := link
	client := youtube.Client{}

	video, err := client.GetVideo(videoID)
	if err != nil {
		return "", err
	}

	formats := video.Formats.WithAudioChannels()
	stream, _, err := client.GetStream(video, &formats[2])
	if err != nil {
		return "", err
	}

	fileVideo := strings.ReplaceAll(video.Title+".mpeg", "/", "|")
	//So.. the character / is and space, so i need to replace it.
	mp3file := strings.ReplaceAll(video.Title+".mp3", "/", "|")

	file, err := os.Create(fileVideo)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		return "", err
	}
	switch runtime.GOOS {
	case "linux", "darwin":
		cmd := exec.Command("ffmpeg", "-i", fileVideo, mp3file)
		if cmd.Run() != nil {
			return "", err
		}
	case "windows":
		cmd := exec.Command("ffmpeg.exe", "-i", fileVideo, mp3file)
		if cmd.Run() != nil {
			return "", err
		}
	default:
		fmt.Println("Unknown OS")
	}
	switch runtime.GOOS {
	case "linux", "darwin":
		del := exec.Command("sh", "-c", "rm *.mpeg").Run()
		if del != nil {
			fmt.Println(del)
		}
	case "windows":
		del := exec.Command("cmd", "/C", "del", "*.mpeg")
		err = del.Run()
		if err != nil {
			log.Error(err)
		}
	default:
		fmt.Println("Unknown OS")
	}
	return video.Title + ".mp3", nil
}
