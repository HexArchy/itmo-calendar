package caldav

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const davNamespace = "DAV:"

// HandlePropPatch handles PROPPATCH requests for the read-only CalDAV server.
// It acknowledges all property changes with 207 Multi-Status without persisting them.
func HandlePropPatch(w http.ResponseWriter, r *http.Request) {
	props, err := parsePropPatchProps(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusMultiStatus)

	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	sb.WriteString(`<D:multistatus xmlns:D="DAV:">`)
	sb.WriteString(`<D:response>`)
	fmt.Fprintf(&sb, `<D:href>%s</D:href>`, xmlEscapeString(r.URL.Path))
	sb.WriteString(`<D:propstat><D:prop>`)

	for _, p := range props {
		if isDavProp(p) {
			fmt.Fprintf(&sb, `<D:%s/>`, p.Local)
		} else {
			fmt.Fprintf(&sb, `<%s xmlns="%s"/>`, p.Local, xmlEscapeString(p.Space))
		}
	}

	sb.WriteString(`</D:prop>`)
	sb.WriteString(`<D:status>HTTP/1.1 200 OK</D:status>`)
	sb.WriteString(`</D:propstat></D:response>`)
	sb.WriteString(`</D:multistatus>`)

	_, _ = io.WriteString(w, sb.String())
}

// parsePropPatchProps extracts property names from a PROPPATCH request body.
func parsePropPatchProps(body io.Reader) ([]xml.Name, error) {
	dec := xml.NewDecoder(body)
	var props []xml.Name
	inProp := false
	depth := 0

	for {
		tok, err := dec.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}

		start, isStart := tok.(xml.StartElement)
		if isStart {
			props, inProp, depth = handleStartElement(start, props, inProp, depth)
			continue
		}

		end, isEnd := tok.(xml.EndElement)
		if isEnd {
			inProp, depth = handleEndElement(end, inProp, depth)
		}
	}

	return props, nil
}

func handleStartElement(t xml.StartElement, props []xml.Name, inProp bool, depth int) ([]xml.Name, bool, int) {
	if isDavProp(t.Name) && t.Name.Local == "prop" {
		return props, true, 0
	}
	if inProp && depth == 0 {
		props = append(props, t.Name)
	}
	if inProp {
		depth++
	}
	return props, inProp, depth
}

func handleEndElement(t xml.EndElement, inProp bool, depth int) (bool, int) {
	if isDavProp(t.Name) && t.Name.Local == "prop" {
		return false, depth
	}
	if inProp {
		depth--
	}
	return inProp, depth
}

func isDavProp(n xml.Name) bool {
	return n.Space == davNamespace || n.Space == ""
}

// xmlEscapeString escapes special XML characters in a string.
func xmlEscapeString(s string) string {
	var b strings.Builder
	_ = xml.EscapeText(&b, []byte(s))
	return b.String()
}
