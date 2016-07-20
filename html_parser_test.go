package tokenizer

import (
	"bytes"
	"testing"
)

const testHTMLData = `<html>
<head>
<title>HTML title</title>
</head>
<body>
The content of the document......
</body>
</html>
`

const testBadHTMLData = `<html>
<head>
</head>
<body>
The content of the document......
</body>
</html>
`

func TestHTMLParser(t *testing.T) {
	r := bytes.NewBuffer([]byte(testHTMLData))
	title, err := readHTMLTitle(r)
	if err != nil {
		t.Fatalf("Unexpected 'readTitle' error: %v\n", err)
	}

	if title != "HTML title" {
		t.Fatalf("Expected title 'HTML title', got '%s'", title)
	}
}

func TestHTMLParserBad(t *testing.T) {
	r := bytes.NewBuffer([]byte(testBadHTMLData))
	_, err := readHTMLTitle(r)
	if err != errorNoTitle {
		t.Fatalf("Expected 'readTitle' error 'errorNoTitle', got %v\n", err)
	}
}
