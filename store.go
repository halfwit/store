package main

import (
	"flag"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"strings"

	"9fans.net/go/plumb"
)

var (
	plumbfile = flag.String("p", "/mnt/plumb/send", "write the message to plumbfile (default /mnt/plumb/send)")
	attributes = flag.String("a", "", "set the attr field of the message (default is empty), expects key=value")
	source = flag.String("s", "", "set the src field of the message (default is store)")
	destination = flag.String("d", "store", "set the dst filed of the message (default is store)")
	directory = flag.String("w", "", "set the wdir field of the message (default is current directory)")
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

func newStoreMsg(mediaType, wdir, arg string, attr *plumb.Attribute) *storeMsg {
	sf := &storeMsg{
		src: os.Args[0],
		dst: "store",
		wdir: wdir,
		msgtype: mediaType,
		attr: attr,
		ndata: len(arg),
		data: arg,
	}
	if *plumbfile != "" {
		sf.src = *plumbfile
	}
	if *source != "" {
		sf.src = *source
	}
	if *destination != "" {
		sf.dst = *destination
	}
	return sf
}

func paramsToAttr(params map[string]string) *plumb.Attribute {
	// Attribute is a linked list - we only get one from content-type, the encoding
	attr := &plumb.Attribute{Name: "", Value: ""}
	for key, value := range params {
		attr.Name = key
		attr.Value = value
	}
	if *attributes != "" {
		attr.Name = strings.TrimLeft(*attributes, "=")
		attr.Value = strings.TrimRight(*attributes, "=")
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
	flag.Parse()
	if flag.Lookup("h") != nil {
		flag.Usage()
		os.Exit(1)
	}
	wdir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	for _, arg := range os.Args[1:] {
		ct, err := content(arg)
		if err != nil {
			log.Fatalf("Error fetching content-type for URL: %s\n", err)
		}
		mediaType, params, err := mime.ParseMediaType(ct)
		if err != nil {
			log.Fatal(err)
		}
		attr := paramsToAttr(params)
		storeMsg := newStoreMsg(mediaType, wdir, arg, attr)
		err = storeMsg.send()
		if err != nil {
			log.Fatal(err)
		}
	}
}
