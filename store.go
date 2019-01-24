package main

import (
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"

	"9fans.net/go/plumb"
)

type storeMsg struct {
	src string
	dst string
	wdir string
	msgtype string
	attr *plumb.Attribute
	ndata int
	data string
}

func (s storeMsg) send() error {
	fd, err := os.OpenFile("/mnt/plumb/send", os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	message := &plumb.Message{
		Src: s.src,
		Dst: s.dst,
		Dir: s.wdir,
		Type: s.msgtype,
		Attr: s.attr,
		Data: []byte(s.data),
	}
	return message.Send(fd)
}

func newStoreMsg(mediaType, wdir string, attr *plumb.Attribute) *storeMsg {
	return &storeMsg{
		src: os.Args[0],
		dst: "",
		wdir: wdir,
		msgtype: mediaType,
		attr: attr,
		ndata: len(os.Args[1]),
		data: os.Args[1],
	}
}

func paramsToAttr(params map[string]string) *plumb.Attribute {
	// Attribute is a linked list - we only get one from content-type, the encoding
	var attr *plumb.Attribute
	for key, value := range params {
		attr.Name = key
		attr.Value = value
	}
	return attr
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
	attr := paramsToAttr(params)
	storeMsg := newStoreMsg(mediaType, wdir, attr)
	err = storeMsg.send()
	if err != nil {
		log.Fatal(err)
	}
}
