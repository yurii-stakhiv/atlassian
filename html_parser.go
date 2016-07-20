package tokenizer

import (
	"errors"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"net/http"
)

var errorNoTitle = errors.New("no title found")

func readHTMLTitle(r io.Reader) (string, error) {
	t := html.NewTokenizer(r)
	for {
		tokenType := t.Next()
		if tokenType == html.ErrorToken {
			if t.Err() == io.EOF {
				break
			}
			return "", t.Err()
		}
		if tokenType != html.StartTagToken {
			continue
		}
		tokenName, _ := t.TagName()
		if string(tokenName) != "title" {
			continue
		}

		tokenType = t.Next()
		if tokenType != html.TextToken {
			break
		}

		return string(t.Text()), nil
	}
	return "", errorNoTitle
}

type HTTPFetcher struct {
}

var defaultFetcher = &HTTPFetcher{}

func (f *HTTPFetcher) FetchTitle(url string) (string, error) {
	resp, err := http.Get(url)
	defer drainResp(resp)
	if err != nil {
		return "", err
	}

	title, err := readHTMLTitle(resp.Body)
	if err != nil {
		return "", err
	}
	return title, nil
}

func drainResp(r *http.Response) error {
	if r != nil {
		io.Copy(ioutil.Discard, r.Body)
		r.Body.Close()
	}
	return nil
}
