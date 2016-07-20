package tokenizer

import (
	"bytes"
	"net/url"
	"regexp"
)

type AsyncRes struct {
	Val Token
	Err error
}

type Token interface {
	ID() string
	Type() string
	Process(chan<- *AsyncRes) bool
}

type Parser interface {
	Parse([]byte) (Token, bool)
}

type Mention string

func (m Mention) ID() string {
	return string(m)
}

func (m Mention) Process(chan<- *AsyncRes) bool {
	return false
}

func (m Mention) Type() string {
	return "mentions"
}

type MentionParser struct {
}

func (p *MentionParser) Parse(b []byte) (Token, bool) {
	if len(b) < 2 {
		return nil, false
	}

	if b[0] != '@' {
		return nil, false
	}

	idx := -1
	for i, v := range b {
		if v == '@' {
			idx = i
		} else {
			break
		}
	}

	return Mention(b[idx+1:]), true
}

type Emoticon string

func (e Emoticon) ID() string {
	return string(e)
}

func (e Emoticon) Process(chan<- *AsyncRes) bool {
	return false
}

func (e Emoticon) Type() string {
	return "emoticons"
}

type EmoticonParser struct {
}

var emoticonRe = regexp.MustCompile("^\\([A-Za-z0-9]{1,15}\\)$")

func (p *EmoticonParser) Parse(b []byte) (Token, bool) {
	if !emoticonRe.Match(b) {
		return nil, false
	}
	return Emoticon(b[1 : len(b)-1]), true
}

type LinkToken struct {
	Url     string `json:"url"`
	Title   string `json:"title"`
	fetcher TitleFetcher
}

func (lt *LinkToken) ID() string {
	return lt.Url
}

func (lt *LinkToken) Process(c chan<- *AsyncRes) bool {
	go func() {
		title, err := lt.fetcher.FetchTitle(lt.Url)
		if err != nil {
			c <- &AsyncRes{nil, err}
			return
		}
		lt.Title = title
		c <- &AsyncRes{lt, nil}
	}()

	return true
}

func (lt *LinkToken) Type() string {
	return "links"
}

type TitleFetcher interface {
	FetchTitle(string) (string, error)
}

type LinkParser struct {
	fetcher TitleFetcher
}

var httpPrefix = []byte("http://")
var httpsPrefix = []byte("https://")

func (p *LinkParser) valid(b []byte) (string, bool) {
	if !bytes.HasPrefix(b, httpPrefix) && !bytes.HasPrefix(b, httpsPrefix) {
		return "", false
	}
	stringUrl := string(b)
	_, err := url.Parse(stringUrl)
	if err != nil {
		return "", false
	}

	return stringUrl, true
}

func (p *LinkParser) Parse(b []byte) (Token, bool) {
	url, ok := p.valid(b)
	if !ok {
		return nil, false
	}
	fetcher := p.fetcher
	if fetcher == nil {
		fetcher = defaultFetcher
	}

	return &LinkToken{url, "", fetcher}, true
}
