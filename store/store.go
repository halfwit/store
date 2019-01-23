package main

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type storeMsg struct {
	src string
	dst string
	wdir string
	msgtype string
	attr map[string]string
	ndata int
	data string
}

func (s storeMsg) send() error {
	f, err := os.OpenFile("/mnt/storage/send", os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	attr := make([]string, 0, len(s.attr))
	for v := range s.attr {
		attr = append(attr, v)
	}
	fmt.Fprintf(f, "%s\n%s\n%s\n%s\n%d\n%s", 
		s.src,
		s.dst,
		s.wdir,
		s.msgtype,
		strings.Join(attr, " "),
		s.ndata,
		s.data,
	)
	return nil
}

func newStoreMsg(mediaType, wdir string, params map[string]string) *storeMsg {
	return &storeMsg{
		src: os.Args[0],
		dst: "",
		wdir: wdir,
		msgtype: mediaType,
		attr: params,
		ndata: len(os.Args[1]),
		data: os.Args[1],
	}
}

func content(testUrl string) (string, error) {

	// We read in 512 bytes 
	buf := make([]byte, 512)

	u, err := url.ParseRequestURI(testUrl)
	if err != nil {
		return "", err
	}
	r, err := http.Get(u.String())
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	n, err := r.Body.Read(buf)
	if err != nil {
		return "", err
	}
	return http.DetectContentType(buf[:n]), nil
}

func main() {
	if (len(os.Args) != 2) {
		log.SetFlags(0)
		log.Fatalf("Usage: %s <URL>\n", os.Args[0])
	}
	wdir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	ct, err := content(os.Args[1])
	if err != nil {
		log.Fatalf("Error fetching content-type for URL: %s\n", err)
	}
	mediaType, params, err := mime.ParseMediaType(ct)
	if err != nil {
		log.Fatal(err)
	}
	storeMsg := newStoreMsg(mediaType, wdir, params)
	err = storeMsg.send()
	if err != nil {
		log.Fatal(err)
	}
}
