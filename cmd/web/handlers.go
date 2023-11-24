package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

	"github.com/lnfu/youtube-downloader/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	ts, err := template.ParseFiles("./ui/html/home.page.tmpl")
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.serverError(w, err)
	}
}

type File struct {
	Key      string `json:"key"`
	Filename string `json:"filename"`
}

func (app *application) getFile(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		app.notFound(w)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.serverError(w, err)
		return
	}

	var f File

	if err := json.Unmarshal(body, &f); err != nil {
		app.serverError(w, err)
		return
	}
	file, err := os.Open("./files/" + f.Filename)
	if err != nil {
		fmt.Println("檔案讀取錯誤")
		app.serverError(w, err)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "video/mp4")
	// w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Disposition", "attachment;")

	_, err = io.Copy(w, file)
	if err != nil {
		app.serverError(w, err)
		return
	}

}

func (app *application) getOriginInfo(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query().Get("v")

	if len(v) == 0 {
		app.notFound(w)
		return
	}

	mi, err := app.origin.Get(v)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {

			app.origin.Insert(v)
			output, err := exec.Command("yt-dlp", "--get-title", "--get-duration", "https://www.youtube.com/watch?v="+v).Output()
			if err != nil || len(output) == 0 {
				app.serverError(w, err)
				return
			}
			lines := strings.Split(string(output), "\n")
			if len(lines) < 2 {
				app.serverError(w, err)
				return
			}
			title := app.formatTitle(lines[0])
			duration := app.parseDuration(lines[1])
			if err != nil {
				app.serverError(w, err)
				return
			}
			app.origin.UpdateInfo(v, title, duration)

			mi, err = app.origin.Get(v)

		} else {
			app.serverError(w, err)
			return
		}
	}

	switch mi.InfoStatus {
	case "running":

		for mi.InfoStatus == "running" {
			time.Sleep(time.Millisecond)
			mi, err = app.origin.Get(v)
			if err != nil {
				app.serverError(w, err)
				return
			}
		}
		break
	case "done":
		break
	case "failure":
		app.serverError(w, err)
		return
	default:
		app.serverError(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"title": mi.Title, "duration": mi.Duration}, nil)

}

func (app *application) getMediaInfo(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query().Get("v")
	t := r.URL.Query().Get("t")

	if len(v) == 0 {
		app.notFound(w)
		return
	}

	if t != "v" && t != "a" {
		app.notFound(w)
		return
	}

	mm, err := app.media.Get(v, t, time.Now())

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {

			key, err := generateRandomString(16)

			if err != nil {
				fmt.Println("金鑰產生錯誤")
				app.serverError(w, err)
				return
			}

			app.media.Insert(v, t, key)

			go app.download(v, t)

			mm, err = app.media.Get(v, t, time.Now())

		} else {
			app.serverError(w, err)
			return
		}
	}

	switch mm.MediaStatus {
	case "running":
		err = app.writeJSON(w, http.StatusOK, envelope{"status": "running"}, nil)
		break
	case "done":
		var extension string
		if mm.Type == "a" {
			extension = ".mp3"
		} else {
			extension = ".mp4"
		}

		err = app.writeJSON(w, http.StatusOK, envelope{"status": "done", "filename": mm.OriginId + extension, "key": mm.AccessKey}, nil)
		break
	case "failure":
		fmt.Println("影片下載失敗")
		app.serverError(w, err)
		return
	default:
		app.serverError(w, err)
		return
	}

}

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}

func (app *application) download(v, t string) {

	var err error
	switch t {
	case "v":
		_, err = exec.Command("yt-dlp", "--format", "bv[ext=mp4]+ba[ext=m4a]/b[ext=mp4]", "--output", "./files/"+v, "https://www.youtube.com/watch?v="+v).Output()
		break
	case "a":
		_, err = exec.Command("yt-dlp", "--format", "ba", "-x", "--audio-format", "mp3", "--output", "./files/"+v, "https://www.youtube.com/watch?v="+v).Output()
		break
	default:
		fmt.Println("錯誤格式")
		return
	}

	if err != nil {
		app.media.DownloadFailure(v, t)
		return
	}
	app.media.DownloadComplete(v, t)
}
