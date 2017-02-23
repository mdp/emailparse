package main

import (
	"bytes"
	"flag"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/mail"
	"os"
	"regexp"
	"strings"
	"time"
)

func getFuncMap(email Email) template.FuncMap {
	funcMap := template.FuncMap{
		"datef": func(d string) (string, error) {
			t, err := time.Parse(time.RFC1123Z, email.Date)
			if err != nil {
				return "", err
			}
			return t.Format(d), nil
		},
		"underscore": func(s string) string {
			re := regexp.MustCompile(" ")
			s = re.ReplaceAllLiteralString(s, "_")
			re = regexp.MustCompile("[^a-zA-Z0-9_-]")
			s = re.ReplaceAllLiteralString(s, "")
			return s
		},
	}
	return funcMap
}

// Email reps our email in a template format
type Email struct {
	Subject string
	Date    string
	From    string
	To      string
	Text    string
}

func getPart(msg *mail.Message, ctype string) (string, error) {
	mediaType, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil {
		log.Fatal(err)
	}
	var body bytes.Buffer
	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(msg.Body, params["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			slurp, err := ioutil.ReadAll(p)
			if err != nil {
				log.Fatal(err)
			}
			if strings.HasPrefix(p.Header.Get("Content-Type"), ctype) {
				body.Write(slurp)
			}
		}
	} else {
		text, err := ioutil.ReadAll(msg.Body)
		if err != nil {
			return body.String(), err
		}
		body.Write(text)
	}
	return body.String(), nil
}

func main() {
	bytes, _ := ioutil.ReadAll(os.Stdin)
	r := strings.NewReader(string(bytes))
	m, err := mail.ReadMessage(r)
	if err != nil {
		log.Fatal(err)
	}
	flag.Parse()
	h := m.Header
	if err != nil {
		log.Fatal(err)
	}
	text, err := getPart(m, "text/plain")
	if err != nil {
		log.Fatal(err)
	}
	email := Email{
		Subject: h.Get("Subject"),
		Date:    h.Get("Date"),
		From:    h.Get("From"),
		Text:    text,
	}
	tmpl, err := template.New("email").Funcs(getFuncMap(email)).Parse(strings.Join(flag.Args(), ""))
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}
	err = tmpl.Execute(os.Stdout, email)
	if err != nil {
		log.Fatal(err)
	}
}
