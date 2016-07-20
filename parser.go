package tokenizer

import (
	"regexp"
)

var (
	mentionRe  = regexp.MustCompile("@([A-Za-z0-9]+)($|[!,\\.;])")
	emoticonRe = regexp.MustCompile("\\(([A-Za-z0-9]{1,15})\\)")
	linkRe     = regexp.MustCompile("(https?://)?([\\da-z\\.-]+)\\.([a-z\\.]{2,6})([/\\w \\.-]*)*/?")
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
	Parse([]byte) ([]Token, bool)
}

type SimpleToken struct {
	val  string
	kind string
}

func (t *SimpleToken) ID() string {
	return t.val
}

func (t *SimpleToken) Process(chan<- *AsyncRes) bool {
	return false
}

func (t *SimpleToken) Type() string {
	return t.kind
}

func (t *SimpleToken) MarshalJSON() ([]byte, error) {
	return []byte("\"" + t.val + "\""), nil
}

type SimpleParser struct {
	Type string
	Re   *regexp.Regexp
}

func (p *SimpleParser) Parse(b []byte) ([]Token, bool) {
	f := p.Re.FindAllStringSubmatch(string(b), -1)
	if len(f) == 0 {
		return nil, false
	}

	res := make([]Token, 0, len(f))
	for _, i := range f {
		res = append(res, &SimpleToken{i[1], p.Type})
	}

	return res, true

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

func (p *LinkParser) Parse(b []byte) ([]Token, bool) {
	f := linkRe.FindAll(b, -1)
	if len(f) == 0 {
		return nil, false
	}

	fetcher := p.fetcher
	if fetcher == nil {
		fetcher = defaultFetcher
	}

	res := make([]Token, 0, len(f))
	for _, i := range f {
		res = append(res, &LinkToken{Url: string(i), fetcher: fetcher})
	}

	return res, true
}
