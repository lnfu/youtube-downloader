package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) parseDuration(str string) int {
	parts := strings.Split(str, ":")

	totalSeconds := 0
	power := 1
	for index := len(parts) - 1; index >= 0; index-- {
		count, _ := strconv.Atoi(parts[index])
		totalSeconds += count * power
		power *= 60
	}
	return totalSeconds
}

func (app *application) formatTitle(str string) string {

	replacedStr := strings.ReplaceAll(str, "/", "_")
	replacedStr = strings.ReplaceAll(replacedStr, "\\", "_")
	replacedStr = strings.ReplaceAll(replacedStr, ".", "_")
	replacedStr = strings.ReplaceAll(replacedStr, ":", "_")
	replacedStr = strings.ReplaceAll(replacedStr, "*", "_")
	replacedStr = strings.ReplaceAll(replacedStr, "?", "_")
	replacedStr = strings.ReplaceAll(replacedStr, "\"", "_")
	replacedStr = strings.ReplaceAll(replacedStr, "'", "_")
	replacedStr = strings.ReplaceAll(replacedStr, ">", "_")
	replacedStr = strings.ReplaceAll(replacedStr, "<", "_")
	replacedStr = strings.ReplaceAll(replacedStr, "|", "_")
	replacedStr = strings.ReplaceAll(replacedStr, " ", "_")
	replacedStr = strings.ReplaceAll(replacedStr, "__", "_")

	return replacedStr

}

type envelope map[string]interface{}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	js = append(js, '\n')
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}
