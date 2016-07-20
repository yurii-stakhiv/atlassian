package tokenizer

import (
	"testing"
)

type testCase struct {
	in       string
	expected string
	ok       bool
}

var mentionCases = []testCase{
	{"@user", "user", true},
	{"@@user", "user", true},
	{"@@@user", "user", true},
	{"@@asdf@user", "user", true},
	{"asdf@user", "user", true},
	{"asdf", "", false},
	{"@", "", false},
	{"", "", false},
}

func TestMentionParser(t *testing.T) {
	p := SimpleParser{"mention", mentionRe}

	for _, testCase := range mentionCases {
		tokens, ok := p.Parse([]byte(testCase.in))
		if testCase.ok != ok {
			t.Errorf("Expected '%s' parse result to be '%t'", testCase.in, testCase.ok)
			continue
		}

		if ok {
			res := tokens[0]
			if res.ID() != testCase.expected {
				t.Fatalf("Expected res to equal '%s', got '%s'", testCase.expected, res.ID())
			}
		}
	}
}

var emoticonCases = []testCase{
	{"(smile)", "smile", true},
	{"(a)", "a", true},
	{"(fifteenchars123)", "fifteenchars123", true},
	{"(sixteenchars12345)", "", false},
	{"()", "", false},
	{"", "", false},
	{"asdf", "", false},
}

func TestEmoticonParser(t *testing.T) {
	p := SimpleParser{"emoticon", emoticonRe}

	for _, testCase := range emoticonCases {
		tokens, ok := p.Parse([]byte(testCase.in))
		if testCase.ok != ok {
			t.Errorf("Expected '%s' parse result to be '%t'", testCase.in, testCase.ok)
			continue
		}

		if ok {
			res := tokens[0]
			if res.ID() != testCase.expected {
				t.Fatalf("Expected res to equal '%s', got '%s'", testCase.expected, res)
			}
		}
	}

}

type testFetcher struct {
	title string
	err   error
}

func (f *testFetcher) FetchTitle(url string) (string, error) {
	return f.title, f.err
}

var linkCases = []testCase{
	{"http://link.com", "test title", true},
	{"https://link.com", "test title", true},
	{"asdf", "test title", false},
	{"google.com", "test title", true},
}

func TestLinkParser(t *testing.T) {
	f := &testFetcher{"test title", nil}
	p := &LinkParser{f}

	for _, testCase := range linkCases {
		tokens, ok := p.Parse([]byte(testCase.in))
		if testCase.ok != ok {
			t.Errorf("Expected '%s' parse result to be '%t'", testCase.in, testCase.ok)
			continue
		}

		if ok {
			res := tokens[0]
			c := make(chan *AsyncRes, 1)
			res.Process(c)
			<-c
			r := res.(*LinkToken)
			if r.Title != testCase.expected {
				t.Fatalf("Expected title to equal '%s', got '%s'", testCase.expected, r.Title)
			}
		}
	}

}
