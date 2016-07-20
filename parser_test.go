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
	{"@@asdf@user", "asdf@user", true},
	{"asdf@user", "", false},
	{"asdf", "", false},
	{"@", "", false},
	{"", "", false},
}

func TestMentionParser(t *testing.T) {
	p := MentionParser{}

	for _, testCase := range mentionCases {
		res, ok := p.Parse([]byte(testCase.in))
		if testCase.ok != ok {
			t.Errorf("Expected '%s' parse result to be '%t'", testCase.in, testCase.ok)
			continue
		}

		if ok && res != Mention(testCase.expected) {
			t.Fatalf("Expected res to equal '%s', got '%s'", testCase.expected, res)
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
	p := EmoticonParser{}

	for _, testCase := range emoticonCases {
		res, ok := p.Parse([]byte(testCase.in))
		if testCase.ok != ok {
			t.Errorf("Expected '%s' parse result to be '%t'", testCase.in, testCase.ok)
			continue
		}

		if ok && res != Emoticon(testCase.expected) {
			t.Fatalf("Expected res to equal '%s', got '%s'", testCase.expected, res)
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
	{"google.com", "test title", false},
}

func TestLinkParser(t *testing.T) {
	f := &testFetcher{"test title", nil}
	p := &LinkParser{f}

	for _, testCase := range linkCases {
		res, ok := p.Parse([]byte(testCase.in))
		if testCase.ok != ok {
			t.Errorf("Expected '%s' parse result to be '%t'", testCase.in, testCase.ok)
			continue
		}

		if ok {
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
