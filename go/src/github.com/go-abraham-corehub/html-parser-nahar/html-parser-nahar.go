package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {

	nameDir := "testCases"
	_, err := ioutil.ReadDir(nameDir)
	if err == nil {
		nameFileInput := strJoin(strJoin(nameDir, "/"), "The Ugly.html")
		iF, err := os.Open(nameFileInput)
		if err == nil {
			html := *bufio.NewReader(iF)

			nameFileOutput := strJoin(nameFileInput, ".txt")
			oF, err := os.Create(nameFileOutput)
			if err == nil {
				txt := *bufio.NewWriter(oF)
				html2txt(&html, &txt)
				txt.Flush()
			} else {
				catchError(err)
			}
		}
	}
}

func html2txt(html *bufio.Reader, txt *bufio.Writer) {
	blockTags := &[]string{
		"<p",
		"h1",
		"h2",
		"h3",
		"h4",
		"h5",
		"h6",
		"ol",
		"/li", // newline for end tags only
		"address",
		"blockquote",
		"dl",
		"dd",
		"dt",
		"div",
		"fieldset",
		"form",
		"noscript",
		"table",
		"/tr",
		"tbody",
		"tfoot",
		"thead",
	}

	specialTags := &[]string{
		"map",
		"script",
		"style",
		"object",
		"applet",
	}

	pS := &parseState{
		false,
		false,
		false,
		true,
		false,
		false,
		0,
		txt,
		blockTags,
		specialTags,
	}

	p, err := html.Peek(1)
	for err == nil {
		switch p[0] {
		case '<':
			if isTag(html) {
				html.Discard(1)
				p, err := html.Peek(1)
				if err == nil {
					if isComment(p[0]) {
						skipComment(html)
					} else if p[0] == '/' {
						//html.Discard(1)
						//fmt.Println(string(p), "isTag, isClosing")
						//dealETags(html, pS)
					} else {
						//dealTags(html, pS)
					}
				}
			}
		}
		html.Discard(1)
		if !pS.done {
			fmt.Print("")
		}
		p, err = html.Peek(1)
	}
}

func printNextNChars(html *bufio.Reader, count int) {
	p, err := html.Peek(count)
	if err == nil {
		fmt.Println(string(p))
	}
}

func fR(html *bufio.Reader, pS *parseState) {
	bts, _ := html.Peek(html.Buffered())
	pS.txt.Write(bts)
}

type parseState struct {
	inTitle, inBody, inPre, hitSpace, newLine, done bool
	titleLen                                        int
	txt                                             *bufio.Writer
	blockTags                                       *[]string
	specialTags                                     *[]string
}

func isTag(html *bufio.Reader) bool {
	p, err := html.Peek(3)
	if err == nil {
		return isTagName(p[1:])
	}
	return false
}

func isTagName(p []byte) bool {
	if p[0] == '/' {
		return isAlpha(p[1]) || isComment(p[1])
	}
	return isAlpha(p[0]) || isComment(p[0])
}

func isAlpha(c byte) bool {
	return (c > 65 && c < 95) || (c > 96 && c < 123)
}

func isAlphaNum(c byte) bool {

	return (c > 65 && c < 95) || (c > 96 && c < 123) || (c > 47 && c < 58)
}

func isComment(c byte) bool {

	return c == '!' || c == '?' || c == '%'
}

func skipComment(html *bufio.Reader) {
	p, err := html.Peek(3)
	if err == nil {
		switch p[0] {
		case '!':
			if string(p[1:]) == "--" {
				html.Discard(3)
				skipPast("-->", html)
			} else {
				html.Discard(1)
				//skipPast(">", html)
			}
		case '?':
			//html.Discard(1)
			//skipPast("?>", html)
		case '%':
			//html.Discard(1)
			//skipPast("%>", html)
		}
	}
}

func skipPast(str string, html *bufio.Reader) {
	//comment := make([]byte, 0)
	l := len(str)
	p, err := html.Peek(l)
	for err == nil {
		if string(p) == str {
			//fmt.Println(string(comment))
			html.Discard(l)
			return
		}
		//comment = append(comment, p[0])
		html.Discard(1)
		p, err = html.Peek(l)
	}
}

func skipTag(html *bufio.Reader) {
	var p []byte
	var err error
	for ; err == nil; _, _, err = html.ReadRune() {
		p, err = html.Peek(1)
		catchError(err)
		if p[0] == '"' {
			html.Discard(1)
			skipPast("\"", html)
		} else if p[0] == '\'' {
			html.Discard(1)
			skipPast("'", html)
		}

		if p[0] == '>' {
			html.ReadRune()
			return
		}
	}
}

func catchError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func dealETags(html *bufio.Reader, pS *parseState) {
	tagName := make([]byte, 0)
	p, err := html.Peek(1)
	if err == nil && isAlpha(p[0]) {
		cH, _, err := html.ReadRune()
		for ; isAlphaNum(byte(cH)) && err == nil; cH, _, err = html.ReadRune() {
			tagName = append(tagName, byte(cH))
		}
	}

	switch string(tagName) {
	case "body":
		doETagBody(pS)
	case "title":
		doETagTitle(html, pS)
	case "pre":
		doETagPre(pS)
	default:
		if isBlockTag(html, pS) {
			doBlockETag(pS)
		}
	}
	skipTag(html)

}

func dealTags(html *bufio.Reader, pS *parseState) {
	tagName := make([]byte, 0)
	p, err := html.Peek(1)
	if err == nil && isAlpha(p[0]) {
		cH, _, err := html.ReadRune()
		for ; isAlphaNum(byte(cH)) && err == nil; cH, _, err = html.ReadRune() {
			tagName = append(tagName, byte(cH))
		}
	}

	switch string(tagName) {
	case "br":
		doTagBr(pS)
	case "hr":
		doTagHr(pS)
	case "body":
		doTagBody(pS)
	case "title":
		doTagTitle(pS)
	case "pre":
		doTagPre(html, pS)
	default:
		if isBlockTag(html, pS) {
			doBlockTag(pS)
		} else if !isSpecialTag(html, pS) {
			skipPast("/", html)
		}
	}
	skipTag(html)
}

func doTagBody(pS *parseState) {
	pS.inBody = true
	pS.hitSpace = true
}

func doETagBody(pS *parseState) {

	//fmt.Println("in body")
	pS.inBody = false
	pS.done = true
	pS.txt.WriteRune('\n')
}

func doETagTitle(html *bufio.Reader, pS *parseState) {
	//fmt.Println("in title")
	if pS.inBody {
		return
	}
	pS.inTitle = false
	pS.titleLen = 0
	pS.txt.WriteRune('\n')

	b, err := html.Peek(1)
	catchError(err)

	if isSpace(b[0]) {
		pS.titleLen--
	}
	if pS.titleLen > 79 {
		pS.titleLen = 80
	}
	for ; pS.titleLen > 0; pS.titleLen-- {
		pS.txt.WriteRune('=')
	}

	pS.txt.WriteString("\n\n\n")

	pS.hitSpace = true

	return
}

func doTagPre(html *bufio.Reader, pS *parseState) {
	pS.inPre = true
	pS.hitSpace = false
	c, err := html.ReadByte()
	for isSpace(c) && err == nil {
		c, err = html.ReadByte()
	}
}

func doETagPre(pS *parseState) {
	pS.inPre = false
	pS.hitSpace = true
	pS.txt.WriteString("\n")
}

func isBlockTag(html *bufio.Reader, pS *parseState) bool {
	for _, bT := range *pS.blockTags {
		if bT[0] == '/' || bT[0] == '<' {
			p, err := html.Peek(len(bT))
			if err == nil {
				if strings.Compare(string(p), bT[1:]) == 0 {
					return true
				} else if strings.Compare(string(p), bT) == 0 {
					return true
				}
			}
		}
	}
	return false
}

func isSpecialTag(html *bufio.Reader, pS *parseState) bool {
	for _, sT := range *pS.specialTags {
		p, err := html.Peek(len(sT))
		if err == nil {
			if strings.Compare(string(p), sT) == 0 {
				return true
			}
		}
	}
	return false
}

func doNewLine(pS *parseState) {
	if !pS.newLine {
		pS.txt.WriteString("\n")
	}
}

func doBlockETag(pS *parseState) {
	doNewLine(pS)
	pS.hitSpace = true
}

func doBlockTag(pS *parseState) {
	doNewLine(pS)
	pS.hitSpace = true
}

func doTagBr(pS *parseState) {
	doNewLine(pS)
	pS.newLine = false
}

func doTagHr(pS *parseState) {
	doNewLine(pS)
	for i := 0; i < 80; i++ {
		pS.txt.WriteString("_")
	}
	pS.txt.WriteString("\n")
	pS.hitSpace = true
}

func doTagTitle(pS *parseState) {
	if pS.inBody {
		return
	}
	pS.inTitle = true
	pS.titleLen = 0
}

func isSpace(c byte) bool {

	return c == ' ' || c == '\n' || c == '\r' || c == '\t' || c == '\v' || c == '\f'
}

func strJoin(strA string, strB string) string {
	return strings.Join([]string{strA, strB}, "")
}
