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
	if err != nil {
		log.Fatal(err)
	}

	/*
			for _, file := range files {
				name := file.Name()
				if name[len(name)-5:] == ".html" {
					nameFileInput := strJoin(strJoin(nameDir, "/"), name)
					dat, err := ioutil.ReadFile(nameFileInput)
					textOutput := []byte("")
					if err == nil {
						textOutput = parse(dat)
					}
					nameFileOutput := strJoin(nameFileInput, ".txt")
					err = ioutil.WriteFile(nameFileOutput, textOutput, 0644)
				}

			}


		nameFileInput := strJoin(strJoin(nameDir, "/"), "The Ugly.html")
		dat, err := ioutil.ReadFile(nameFileInput)
		textOutput := []byte("")
		if err == nil {
			textOutput = parse(dat)
		}
		nameFileOutput := strJoin(nameFileInput, ".txt")
		err = ioutil.WriteFile(nameFileOutput, textOutput, 0644)
	*/
	nameFileInput := strJoin(strJoin(nameDir, "/"), "The Ugly.html")
	dat, err := ioutil.ReadFile(nameFileInput)
	nameFileOutput := strJoin(nameFileInput, ".txt")
	textOutput := parse02(&dat)
	err = ioutil.WriteFile(nameFileOutput, *textOutput, 0644)
	fmt.Println(string(*textOutput))

	/*
		pS := &parseState{false, false}
		fmt.Println(isInPre("<pre>", pS))
		fmt.Println(isInPre("ereyeyeyry", pS))
		fmt.Println(isInPre("</pre>", pS))
		fmt.Println(isInPre("<hr/>", pS))
	*/

	/*
		nameFileInput := strJoin(strJoin(nameDir, "/"), "The Ugly.html")
		f, err := os.Open(nameFileInput)
		if err == nil {
			parse(bufio.NewReader(f))
		}
	*/
}

type myReadInterface interface {
	setFile(string) *myReader
	Read() rune
}

type myReader struct {
	myFile    *os.File
	bReader   *bufio.Reader
	countRead int
	err       error
}

func (r myReader) setFile(nameFile string) myReader {
	r = myReader{}
	r.myFile, r.err = os.Open(nameFile)
	r.bReader = bufio.NewReader(r.myFile)
	r.countRead = 0
	return r
}

//Read reads the file contents through a Reader and returns the
func (r myReader) Read() rune {
	if r.err == nil {
		rU, _, _ := r.bReader.ReadRune()
		return rU
	}
	return rune(0)
}

func parse(rP *bufio.Reader) {
	cH, _, err := rP.ReadRune()
	d := make([]byte, 0)
	for err == nil {
		switch cH {
		case '<':
			d, err = rP.Peek(3)
			fmt.Println(isValidTagStart(d))
		case '>':
		default:
		}
		cH, _, err = rP.ReadRune()
	}
}

type parseState struct {
	inTag     bool
	inComment bool
	isInPre   bool
	isScript  bool
	hitSpace  bool
	skip      bool
}

type tagInfo struct {
	open        bool
	close       bool
	selfclosing bool
}

type element struct {
	tagName []byte
	content []byte
	tInfo   tagInfo
}

func testRead(nameFile string) {
	f, err := os.Open(nameFile)
	if err == nil {
		readerData := bufio.NewReader(f)
		runeData, rSize, err := readerData.ReadRune()
		for err == nil {
			fmt.Println(string(runeData), rSize, err)
			runeData, rSize, err = readerData.ReadRune()
		}
	}
}

func parse02(dP *[]byte) *[]byte {
	d := *dP
	var dT []byte
	elem, content := &element{}, make([]byte, 0)
	e := len(d)
	i, sT := 0, 0
	validText := false
	pS := &parseState{false, false, false, false, false, false}

	for ; i < e; i++ {
		//fmt.Print(string(d[i]), " | ", tagName, " | ", inTag, isInPre(tagName, pS), isInScript(tagName, pS), "\n")
		switch d[i] {
		case '<':
			if isValidTagStart(d[i:]) {
				sT = i
				pS.inTag = true
				l := len(d[i:])
				if l > 3 {
					pS.inComment = d[i+1] == '!' && d[i+2] == '-' && d[i+3] == '-'
					//fmt.Println(string(d[i]), pS)
				}
			}
		case '>':
			if pS.inTag && isValidTag(d[sT:i+1]) {
				elem = elem.getTag(d[sT : i+1])
				//fmt.Print("'", string(d[sT:i+1]), ", ", string(element.tagName), "'\n")
				//fmt.Println("")
				pS.inTag = false
			}
			pS.inComment = false
		default:
			if !pS.isInPre {
				switch d[i] {
				case ' ':
					if pS.hitSpace {
						pS.skip = true
					}
					pS.hitSpace = true
				case '\n', '\t', '\r', '\f':
					pS.skip = true
					pS.hitSpace = false
				default:
					pS.skip = false
					pS.hitSpace = false
				}

				switch string(elem.tagName) {
				case "title":
					if elem.tInfo.close {
						elem.content = append(make([]byte, 0), content...)
						elem.tInfo.close = false
						//content = make([]byte, 0)
						//fmt.Println(string(element.content))
					}

				}
			}

			//validText = !pS.isScript && !pS.inTag && !pS.inComment
			validText = pS.inComment
			if validText {
				content = append(content, d[i])
				fmt.Print(string(d[i]))
			}

		}
		updateState(elem, pS)
		//fmt.Print("'", string(d[i]), "'", d[i], pS, "\n")
	}
	return &dT
}

func (e element) getTag(b []byte) *element {
	l := len(b)
	e.tagName = make([]byte, 0)
	e.content = make([]byte, 0)
	e.tInfo = tagInfo{true, false, false}
	i := 1
	if b[1] == '/' {
		e.tInfo = tagInfo{false, true, false}
		i = 2
	} else if b[l-1] == '/' {
		e.tInfo = tagInfo{false, false, true}
	}
	//fmt.Println(string(b[i : i+3]))
	if b[i] == '!' && i+3 < l {
		if string(b[i:i+3]) == "!--" {
			e.tagName = b[i : i+3]
			i = i + 3
			//fmt.Print(string(b), ",", l, ",", i, ",", string(b[i-2]), ",", string(b[i-1]), ",", string(b[i]))
			for ; i < l && !(b[i-2] == '-' && b[i-1] == '-' && b[i] == '>'); i++ {
				e.content = append(e.content, b[i])
			}
		}
	}
	for ; i < l && b[i] != ' ' && b[i] != '>' && b[i] != '/'; i++ {
		e.tagName = append(e.tagName, b[i])
	}
	//fmt.Println(string(e.tagName))
	return &e
}

func updateState(elem *element, pS *parseState) {
	pS.inTag = elem.tInfo.open && !elem.tInfo.close && !elem.tInfo.selfclosing
	//fmt.Println(string(elem.tagName))
	switch string(elem.tagName) {
	case "script":
		pS.isScript = elem.tInfo.open && !elem.tInfo.close

	case "pre":
		pS.isInPre = elem.tInfo.open && !elem.tInfo.close
	}
}

func parse01(dat []byte) []byte {
	e := len(dat)
	out := make([]byte, 0)
	text := make([]byte, 0)
	tS := []string{
		"<title>",
		"<p>",
		"<h2>",
		"<pre>",
		"<a",
		"<div>",
		"<div",
		"<br>",
		"<hr>",
		"<hr />",
		"<hr/>",
		"<i>",
		"<span",
	}

	for i := 0; i < e; i++ {
		if !isComment(dat, i, e) {
			for _, tag := range tS {
				l := len(tag)
				dTag := string(dat[i : i+l])
				if i+l < e && dTag == tag {
					switch dTag {
					case "<title>":
						text, i = formatTitle(dat, i+l, e)
						text = stripSpace(text)
						text = append([]byte("\n"), text...)
						text = append(text, []byte("\n\n")...)
						out = append(out, text...)
					case "<p>":
						text, i = formatPara(dat, i+l, e)
						text = stripSpace(text)
						text = append([]byte("\n"), text...)
						out = append(out, text...)
					case "<h2>":
						text, i = formatH2(dat, i+l, e)
						text = stripSpace(text)
						out = append(out, text...)
					case "<pre>":
						text, i = formatPre(dat, i+l, e)
						out = append(out, text...)
					case "<a":
						text, i = formatA(dat, i+l, e)
						text = stripSpace(text)
						text = append([]byte("\n"), text...)
						out = append(out, text...)
					case "<div>":
						text, i = formatDiv(dat, i+l, e)
						text = stripSpace(text)
						text = append([]byte("\n"), text...)
						out = append(out, text...)
					case "<div":
						text, i = formatDivA(dat, i+l, e)
						text = stripSpace(text)
						text = append([]byte("\n"), text...)
						out = append(out, text...)
					case "<i>":
						text, i = formatI(dat, i+l, e)
						text = stripSpace(text)
						out = append(out, text...)
					case "<span":
						text, i = formatSpan(dat, i+l, e)
						text = stripSpace(text)
						out = append(out, text...)
					case "<br>":
						out = append(out, []byte("\n")...)

					case "<hr>", "<hr />", "<hr/>":
						out = append(out, []byte("\n")...)
						for j := 0; j < 80; j++ {
							out = append(out, []byte("_")...)
						}
						out = append(out, []byte("\n")...)
					}
				}
			}
		} else {
			i = skipComment(dat, i, e)
		}
	}
	return formatText(out)
}

func strJoin(strA string, strB string) string {
	return strings.Join([]string{strA, strB}, "")
}

func stripSpace(str []byte) []byte {
	outStr := make([]byte, 0)
	spaceCount := 0
	for _, cH := range str {
		if cH == ' ' {
			if spaceCount < 1 && cH != '\n' && cH != '\t' && cH != 13 {
				outStr = append(outStr, byte(cH))
				spaceCount = 0
			}
			spaceCount++
		} else if cH != '\n' && cH != '\t' && cH != 13 {
			outStr = append(outStr, byte(cH))
			spaceCount = 0
		}
	}
	outStr = []byte(strings.TrimSpace(string(outStr)))
	return outStr
}

func formatText(dat []byte) []byte {
	out := make([]byte, 0)
	tS := map[string]string{
		"<q>":   "</q>",
		"&amp;": "&",
	}
	bR := byte(0)
	e := len(dat)
	for sT, eT := range tS {
		if sT[0] != '&' {
			switch sT {
			case "<q>":
				bR = '"'
			case "<h1>", "<h2>", "<h3>":
				bR = 0
			}
			for i := 0; i < e; i++ {
				lST := len(sT)
				lET := len(eT)
				dSTag := string(dat[i : i+lST])
				if i+lST < e && dSTag == sT {
					if bR != 0 {
						out = append(out, bR)
					}
					for j := i + lST; j < e; j++ {
						dETag := string(dat[j : j+lET])
						if j+lET < e && dETag == eT {
							i = j + lET
							if bR != 0 {
								out = append(out, bR)
							}
							break
						}
						out = append(out, dat[j])
						i = j
					}
				}
				out = append(out, dat[i])
			}
			dat = out
			e = len(dat)
			out = make([]byte, 0)
		} else {
			for i := 0; i < e; i++ {
				l := len(sT)
				if i+l < e && string(dat[i:i+l]) == sT {
					out = dat[:i]
					out = append(out, eT...)
					out = append(out, dat[i+l:]...)
					dat = out
					e = len(dat)
				}
			}
		}

	}

	return dat
}

func getContent(dat []byte) []byte {

	return dat
}

func formatTitle(dat []byte, i int, e int) ([]byte, int) {
	text := make([]byte, 0)
	et := "</title>"
	l := len(et)
	for j := i; j+l < e-1; j++ {
		if string(dat[j:j+l]) == et {
			i = j + l
			break
		}
		if isComment(dat, j, e) {
			j = skipComment(dat, j, e)
		}
		text = append(text, dat[j])
		i = j
	}
	return text, i
}

func formatDiv(dat []byte, i int, e int) ([]byte, int) {
	text := make([]byte, 0)
	et := "</div>"
	l := len(et)
	for j := i; j+l < e-1; j++ {
		if string(dat[j:j+l]) == et {
			i = j + l
			break
		}
		if isComment(dat, j, e) {
			j = skipComment(dat, j, e)
		}
		text = append(text, dat[j])
		i = j
	}
	return text, i
}

func formatI(dat []byte, i int, e int) ([]byte, int) {
	text := make([]byte, 0)
	et := "</i>"
	l := len(et)
	for j := i; j+l < e-1; j++ {
		if string(dat[j:j+l]) == et {
			i = j + l
			break
		}
		if isComment(dat, j, e) {
			j = skipComment(dat, j, e)
		}
		text = append(text, dat[j])
		i = j
	}
	return text, i
}

func isComment(dat []byte, i int, e int) bool {
	if i+3 < e {
		if dat[i] == '<' && dat[i+1] == '!' && dat[i+2] == '-' && dat[i+3] == '-' {
			return true
		}
	}
	return false
}

func skipComment(dat []byte, i int, e int) int {
	for j := i; j < e; j++ {
		if dat[j] == '-' && dat[j+1] == '-' && dat[j+2] == '>' {
			return j + 3
		}
		i = j
	}
	return i + 1
}

func formatPara(dat []byte, i int, e int) ([]byte, int) {
	text := make([]byte, 0)
	et := "</p>"
	l := len(et)
	k := 0
	newTag := false
	textN := make([]byte, 0)
	for j := i; j < e; j++ {
		if string(dat[j:j+l]) == et {
			i = j + l
			break
		} else if isComment(dat, j, e) {
			j = skipComment(dat, j, e)
		} else if dat[j] == '<' {
			k = j
			newTag = true
			textN = append(textN, text...)
			text = append(text, dat[j])
		} else {
			text = append(text, dat[j])
			i = j
		}
		newTag = false
	}
	if newTag {
		i = k
		return textN, i
	}
	return text, i
}

func formatH2(dat []byte, i int, e int) ([]byte, int) {
	text := make([]byte, 0)
	et := "</h2>"
	l := len(et)
	for j := i; j < e; j++ {
		if string(dat[j:j+l]) == et {
			i = j + l
			break
		} else if isComment(dat, j, e) {
			j = skipComment(dat, j, e)
		} else {
			text = append(text, dat[j])
			i = j
		}
	}
	return text, i
}

func formatPre(dat []byte, i int, e int) ([]byte, int) {
	text := make([]byte, 0)
	et := "</pre>"
	l := len(et)
	for j := i; j < e; j++ {
		if string(dat[j:j+l]) == et {
			i = j + l
			break
		} else if isComment(dat, j, e) {
			j = skipComment(dat, j, e)
		} else if dat[j] == '<' {
			j = skipTag(dat, j, e)
		} else {
			text = append(text, dat[j])
			i = j
		}
	}
	return text, i
}

func formatA(dat []byte, i int, e int) ([]byte, int) {
	text := make([]byte, 0)
	et := "</a>"
	l := len(et)

	if string(dat[i-2:i+1]) == "<a " {
		for j := i + 1; j < e && string(dat[j:j+2]) != "\">"; j++ {
			i = j
		}
		i += 2
	}

	for j := i + 1; j < e; j++ {
		if string(dat[j:j+l]) == et {
			i = j + l
			break
		} else if isComment(dat, j, e) {
			j = skipComment(dat, j, e)
		} else {
			text = append(text, dat[j])
			i = j
		}
	}
	return text, i
}

func formatSpan(dat []byte, i int, e int) ([]byte, int) {
	text := make([]byte, 0)
	et := "</span>"
	l := len(et)

	if string(dat[i-5:i+1]) == "<span " {
		for j := i + 1; j < e && string(dat[j:j+2]) != "\">"; j++ {
			i = j
		}
		i += 2
	}

	for j := i + 1; j < e; j++ {
		if string(dat[j:j+l]) == et {
			i = j + l
			break
		} else if isComment(dat, j, e) {
			j = skipComment(dat, j, e)
		} else {
			text = append(text, dat[j])
			i = j
		}
	}
	return text, i
}

func formatDivA(dat []byte, i int, e int) ([]byte, int) {
	text := make([]byte, 0)
	et := "</div>"
	l := len(et)

	if string(dat[i-4:i+1]) == "<div " {
		for j := i + 1; j < e && string(dat[j:j+2]) != "\">"; j++ {
			i = j
		}
		i += 2
	}

	for j := i + 1; j < e; j++ {
		if string(dat[j:j+l]) == et {
			i = j + l
			break
		} else if isComment(dat, j, e) {
			j = skipComment(dat, j, e)
		} else {
			text = append(text, dat[j])
			i = j
		}
	}
	return text, i
}

func skipTag(dat []byte, i int, e int) int {
	tagsToSkip := []string{
		"<script>",
	}
	tSkip := make(map[string]int)

	for _, t := range tagsToSkip {
		tSkip[t] = len(t)
	}

	if !isComment(dat, i, e) {
		for tag, l := range tSkip {
			sTag := string(dat[i : i+l])
			if i+l < e && sTag == tag {
				i = findClosingTag(dat, []byte(sTag), i, e)
			}
		}

	} else {
		i = skipComment(dat, i, e)
		return i
	}

	return i
}

func isValidTag(d []byte) bool {
	l := len(d)
	//fmt.Println(string(d), d)
	if l > 2 {
		s := isValidTagStart(d)
		e01 := d[l-3] == ';' && d[l-2] == '"' && d[l-1] == '>'
		e02 := d[l-2] == '/' && d[l-1] == '>'
		e03 := d[l-2] == ' ' && d[l-1] == '>'
		e04 := isAlpha(d[l-2]) && d[l-1] == '>'
		e05 := d[l-2] == '"' && d[l-1] == '>'
		e06 := d[l-2] == '-' && d[l-1] == '>'
		e07 := d[l-2] == '?' && d[l-1] == '>'
		return s && (e01 || e02 || e03 || e04 || e05 || e06 || e07)
	}
	return false
}

func isValidTagStart(d []byte) bool {
	l := len(d)
	if l > 2 {
		s01 := d[0] == '<' && isAlpha(d[1])
		s02 := d[0] == '<' && d[1] == '/' && isAlpha(d[2])
		s03 := d[0] == '<' && d[1] == '!'
		s04 := d[0] == '<' && d[1] == '?'
		s := s01 || s02 || s03 || s04
		return s
	}
	return false
}

func isInScript(elem element, pS *parseState) bool {
	if string(elem.tagName) == "script" && elem.tInfo.open == true {
		pS.isScript = true
	} else if string(elem.tagName) == "script" && elem.tInfo.close == true {
		pS.isScript = false
	}

	return pS.isScript
}

func isInPre(tag []byte, pS *parseState) bool {
	if string(tag) == "<pre>" {
		pS.isInPre = true
	} else if string(tag) == "</pre>" {
		pS.isInPre = false
	}
	return pS.isInPre
}

func findClosingTag(dat []byte, tagName []byte, i int, e int) int {

	for j := i; j+1+len(tagName)+1 < e; j++ {
		a := string(dat[j+2 : j+1+len(tagName)])
		b := string(tagName[1:])
		if dat[j] == '<' && dat[j+1] == '/' && a == b {
			return j + len(tagName) + 2
		}
	}
	return i
}

func isAlpha(d byte) bool {
	return d > 64 && d < 91 || d > 96 && d < 123
}
