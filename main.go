package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 1 define handler function
func sayHello(w http.ResponseWriter, r *http.Request) {
	// allow cross origin
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Write([]byte("<p>Hello World</p>"))
}

func main() {
	// a.1 create file handler
	fileHanlder := http.FileServer(http.Dir("./video/"))
	// a.2 register handler to servemux
	http.Handle("/video/", http.StripPrefix("/video/", fileHanlder))

	// b.2 register handler
	http.HandleFunc("/api/upload", uploadHandler)

	// c.2 register vlog list handler
	http.HandleFunc("/api/list", getFileListHandler)

	http.HandleFunc("/say-hello", sayHello)
	// 3 run the serve
	http.ListenAndServe(":8008", nil)

}

// b.1 define service logic
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// allow cross origin
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 1) -> corp 1st 10MB of r.Body and try to parse.
	// if no err then the file is less than 10MB
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024*10)
	err := r.ParseMultipartForm(1024 * 1024 * 10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 2) Acquire the file uploaded
	file, fileHeader, err := r.FormFile("uploadFile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 3) check file format, if not .mp4 then give error
	fFmt := strings.HasSuffix(fileHeader.Filename, ".mp4")
	if fFmt {
		http.Error(w, "Format is not .mp4", http.StatusInternalServerError)
		return
	}

	// 4) get random file name to prevent dup file name
	md5Byte := md5.Sum([]byte(fileHeader.Filename + time.Now().String()))
	md5Str := fmt.Sprintf("%x", md5Byte)
	newFileName := md5Str + ".mp4"

	// 5) save into ./video folder
	dst, err := os.Create("./video/" + newFileName)
	defer dst.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}

// 3.1 vlog get file list handler
func getFileListHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("getting the list of the videos")

	// allow cross origin
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	files, _ := filepath.Glob("video/*")

	resp := []string{}
	for _, file := range files {
		resp = append(resp, "http://"+r.Host+"/video/"+filepath.Base(file))
	}
	respJSON, _ := json.Marshal(resp)
	w.Write([]byte(respJSON))
}
