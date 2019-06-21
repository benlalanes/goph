package link

import (
	"container/list"
	"golang.org/x/net/html"
	"strings"
)

type Link struct {
	Href string
	Text string
}

func ParseLinks(root *html.Node) []*Link {

	var links []*Link

	// We use Go's built-in linked list type as a queue for doing
	// breadth-first search on the root node.
	q := list.New()
	q.PushBack(root)

	for curr := q.Front(); curr != nil; curr = curr.Next() {

		node := curr.Value.(*html.Node)

		// Create link struct and append if this node refers to
		// an anchor element.
		if node.Type == html.ElementNode && node.Data == "a" {

			lnk := &Link{}

			// Get the href attribute of the anchor element.
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					lnk.Href = attr.Val
					break
				}
			}

			// Try to get the text of the anchor element by
			// checking all its children.
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.TextNode {
					lnk.Text = strings.TrimSpace(c.Data)
					break
				}
			}

			links = append(links, lnk)

			continue

		}

		// Otherwise, add all children of this node to the queue.
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			q.PushBack(c)
		}

	}

	return links
}
