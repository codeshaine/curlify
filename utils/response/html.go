package response

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func FormatHTMLResponse(body []byte) string {
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return fmt.Sprintf("Error parsing HTML: %v", err)
	}

	var buf bytes.Buffer
	formatNode(&buf, doc, 0)
	return buf.String()
}

func formatNode(buf *bytes.Buffer, n *html.Node, level int) {
	if n.Type == html.ElementNode {
		buf.WriteString(strings.Repeat("  ", level))
		buf.WriteString("<" + n.Data)
		for _, attr := range n.Attr {
			buf.WriteString(fmt.Sprintf(` %s="%s"`, attr.Key, attr.Val))
		}
		buf.WriteString(">\n")
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		formatNode(buf, c, level+1)
	}

	if n.Type == html.ElementNode {
		buf.WriteString(strings.Repeat("  ", level))
		buf.WriteString("</" + n.Data + ">\n")
	}
}
