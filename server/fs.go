package server

import (
	"github.com/hari-govind/liveserver-go/config"

	"bytes"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

const HEAD_TAG_LEN = len("<head>")

func isHTML(filename string) bool {
	if strings.HasSuffix(filename, ".html") || strings.HasSuffix(filename, ".htm") {
		return true
	}
	return false
}

// TODO Instead of using http.Fileserver and using embedded structs to inject script, use http.ServeContent and serve file using ResponseWriter?
// To inject script tag to html files.
// Size will return len(SCRIPT_TAG) bytes more than actual file size in case of .html files
// so that script tag of len(SCRIPT_TAG) bytes can be appended to html files while reading
type injectScriptFileInfo struct {
	os.FileInfo
	size int64
}

type injectScriptFile struct {
	http.File
	filename     string
	numBytesRead *int64 // number of bytes read
}

type injectScriptFileSystem struct {
	http.FileSystem
}

func (fileInfo injectScriptFileInfo) Size() int64 {
	if isHTML(fileInfo.Name()) {
		return fileInfo.size + int64(len(SCRIPT_TAG))
	}
	return fileInfo.size
}

func (file injectScriptFile) Stat() (os.FileInfo, error) {
	fileInfo, err := os.Stat(path.Join(config.GetConfig().Root, file.filename))
	return injectScriptFileInfo{fileInfo, fileInfo.Size()}, err
}

func (file injectScriptFile) Read(p []byte) (n int, err error) {
	fileIsHtml := isHTML(file.filename)
	contentLength := len(p)
	if fileIsHtml {
		contentLength -= len(SCRIPT_TAG) //length will be len(p) once the SCRIPT_TAG is appended, so substracting it now
	}
	contents := make([]byte, contentLength)
	osFile, err := os.Open(path.Join(config.GetConfig().Root, file.filename))
	defer osFile.Close()
	if err != nil {
		log.Println("Cannot open file", err)
	}

	osFile.ReadAt(contents, *file.numBytesRead)

	*file.numBytesRead += int64(len(p))

	if fileIsHtml {
		headIndex := bytes.Index(contents, []byte("<head>"))
		if headIndex != -1 {
			// Append the script tag after <head>
			headIndex += HEAD_TAG_LEN
			contents = append(contents[:headIndex], append(SCRIPT_TAG, contents[headIndex:]...)...)
		} else {
			log.Println("Did not find <head> in ", file.filename)
		}
	}

	copy(p, contents)
	if err != nil {
		log.Println("Cannot close file", err)
	}
	return len(p), err
}

func (fsys injectScriptFileSystem) Open(filename string) (http.File, error) {
	file, err := fsys.FileSystem.Open(filename)
	numBytes := int64(0)
	return injectScriptFile{file, filename, &numBytes}, err
}
