package link

import (
	"golang.org/x/net/html"
	"strings"
	"testing"
)

func TestParseLinks(t *testing.T) {

	tt := []struct {
		input  string
		expect []Link
	}{
		{
			input: `<a href='https://google.com'>Google</a>`,
			expect: []Link{
				Link{Href: "https://google.com", Text: "Google"},
			},
		},
	}

	for i, test := range tt {
		root, err := html.Parse(strings.NewReader(test.input))
		if err != nil {
			t.Fatal(err)
		}

		links := ParseLinks(root)

		if len(links) != len(test.expect) {
			t.Errorf("test %d - wanted %d links, got %d", i, len(links), len(test.expect))
		}

		for idx := 0; idx < len(links); idx++ {
			if *links[idx] != test.expect[idx] {
				t.Errorf("link %d - link parsed incorrectly; wanted %#v, got %#v",
					idx, *links[idx], test.expect[idx])
			}
		}

	}
}
