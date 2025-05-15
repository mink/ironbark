package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
)

type node struct {
	XMLName xml.Name
	Attr    []xml.Attr `xml:",any,attr"`
	Content []byte     `xml:",innerxml"`
	Nodes   []node     `xml:",any"`
}

var selfClosingPathSet = map[string]struct{}{
	"epp/greeting/dcp/access/all":              {},
	"epp/greeting/dcp/access/none":             {},
	"epp/greeting/dcp/access/null":             {},
	"epp/greeting/dcp/access/personal":         {},
	"epp/greeting/dcp/access/personalAndOther": {},
	"epp/greeting/dcp/access/other":            {},

	"epp/greeting/dcp/statement/purpose/admin":   {},
	"epp/greeting/dcp/statement/purpose/contact": {},
	"epp/greeting/dcp/statement/purpose/prov":    {},
	"epp/greeting/dcp/statement/purpose/other":   {},

	"epp/greeting/dcp/statement/recipient/ours":      {},
	"epp/greeting/dcp/statement/recipient/other":     {},
	"epp/greeting/dcp/statement/recipient/public":    {},
	"epp/greeting/dcp/statement/recipient/same":      {},
	"epp/greeting/dcp/statement/recipient/unrelated": {},

	"epp/greeting/dcp/statement/retention/business":   {},
	"epp/greeting/dcp/statement/retention/indefinite": {},
	"epp/greeting/dcp/statement/retention/legal":      {},
	"epp/greeting/dcp/statement/retention/none":       {},
	"epp/greeting/dcp/statement/retention/stated":     {},
}

func ConvertSelfClosingTags(input []byte) []byte {
	var root node
	if err := xml.Unmarshal(input, &root); err != nil {
		return input
	}
	var buf bytes.Buffer
	writeNode(&buf, root, nil, 0)
	return bytes.TrimRight(buf.Bytes(), "\n")
}

func writeNode(buf *bytes.Buffer, n node, path []string, level int) {
	tag := n.XMLName.Local
	fullPath := append(path, tag)
	fullKey := strings.Join(fullPath, "/")

	writeIndent(buf, level)
	buf.WriteByte('<')
	buf.WriteString(tag)
	writeAttributes(buf, n.Attr)

	isEmpty := len(bytes.TrimSpace(n.Content)) == 0 && len(n.Nodes) == 0
	if isEmpty && isSelfClosable(fullKey) {
		buf.WriteString("/>\n")
		return
	}

	buf.WriteByte('>')

	if len(n.Nodes) > 0 {
		buf.WriteByte('\n')
		for _, child := range n.Nodes {
			writeNode(buf, child, fullPath, level+1)
		}
		writeIndent(buf, level)
		buf.WriteString("</")
		buf.WriteString(tag)
		buf.WriteString(">\n")
	} else {
		buf.Write(n.Content)
		buf.WriteString("</")
		buf.WriteString(tag)
		buf.WriteString(">\n")
	}
}

func isSelfClosable(path string) bool {
	_, ok := selfClosingPathSet[path]
	return ok
}

func writeIndent(buf *bytes.Buffer, level int) {
	if level > 0 {
		buf.Write(bytes.Repeat([]byte("  "), level))
	}
}

func writeAttributes(buf *bytes.Buffer, attrs []xml.Attr) {
	for _, attr := range attrs {
		if attr.Name.Space != "" {
			_, _ = fmt.Fprintf(buf, ` %s:%s="%s"`, attr.Name.Space, attr.Name.Local, attr.Value)
		} else {
			_, _ = fmt.Fprintf(buf, ` %s="%s"`, attr.Name.Local, attr.Value)
		}
	}
}
