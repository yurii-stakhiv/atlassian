package tokenizer

import (
	"bufio"
	"bytes"
	"io"
)

var defautlParsers = []Parser{&MentionParser{}, &EmoticonParser{}, &LinkParser{}}

type TokenMap map[string][]Token
type Set map[string]bool

type Scanner struct {
	parsers []Parser
}

func NewScanner(parsers []Parser) *Scanner {
	if parsers == nil {
		parsers = defautlParsers
	}
	return &Scanner{parsers}
}

// Scan io.Reader and return token map
func (s *Scanner) Scan(r io.Reader) (TokenMap, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)

	m := make(TokenMap)
	set := make(Set)

	asyncCount := 0
	asyncRes := make(chan *AsyncRes)

	for scanner.Scan() {
		word := scanner.Bytes()
		if len(word) < 2 {
			continue
		}

		for _, p := range s.parsers {
			token, ok := p.Parse(word)
			if !ok {
				continue
			}
			// Not sure if deduplication is needed
			key := token.Type() + ":" + token.ID()
			if set[key] {
				continue
			}
			set[key] = true

			if token.Process(asyncRes) {
				asyncCount++
			} else {
				m[token.Type()] = append(m[token.Type()], token)
			}
			break
		}
	}

	for i := 0; i < asyncCount; i++ {
		res := <-asyncRes
		if res.Err == nil {
			token := res.Val
			m[token.Type()] = append(m[token.Type()], token)
		}
	}
	return m, nil
}

func (s *Scanner) ScanBytes(b []byte) (TokenMap, error) {
	return s.Scan(bytes.NewBuffer(b))
}
