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

	nameFileInput := strJoin(strJoin(nameDir, "/"), "The Bad.html")
	f, err := os.Open(nameFileInput)
	if err != nil {
		fmt.Println(err)
		return
	}
	rDBuf := *bufio.NewReader(f)

	parse(&rDBuf)
}

func parse(rDB *bufio.Reader) {
	pS := &parseState{}
	pSp := parseState{}
	el := element{}
	els := make([]element, 0)
	tag := make([]byte, 0)
	cH, _, err := rDB.ReadRune()
	for err == nil {
		switch cH {
		case '<':
			if !pS.isPre && !pS.isScript && !pS.isComment && !pS.inTag && isValidTagStart(rDB, pS) {
				pS.inTag = true
			} /*
				if pS.isPre {
					d, err := rDB.Peek(4)
					if err == nil && string(d) == "/pre" {
						pS.isPre = false
						pS.inTag = true
						rDB.Discard(4)
					}

					d, err = rDB.Peek(6)
					if err == nil && string(d) == "script" {
						pS.inTag = true
						pS.isScript = true
						rDB.Discard(6)
					}

				}

				if pS.isScript {
					d, err := rDB.Peek(7)
					if err == nil && string(d) == "/script" {
						pS.inTag = true
						pS.isScript = false
						rDB.Discard(7)
					}
				}

				d, err := rDB.Peek(2)
				if err == nil && string(d) == "--" {
					pS.isComment = true
					pS.inTag = true
					rDB.Discard(2)
				}
			*/

		case '>':
			if pS.inTag && isValidTagEnd(tag, pS) {
				el = *el.parseElement(tag)
				if len(el.tag.name) > 0 {
					els = append(els, el)
				}
				updateState(el, pS)
				el = element{}
				tag = []byte{}
				pS.inTag = false
			} else if pS.inTag && !isValidTagEnd(tag, pS) {
				tag = append(tag, byte(cH))
			}
		default:
			if pS.inTag {
				tag = append(tag, byte(cH))
			}
			if !pS.isPre {
				switch cH {
				case ' ':
					if pS.isSpace {
						pS.isSkip = true
					}
					pS.isSpace = true
				case '\n', '\t', '\r', '\f':
					pS.isSkip = true
					pS.isSpace = false
				default:
					pS.isSkip = false
					pS.isSpace = false
				}
			}
			pSp.inTag = pS.inTag
			updateState(el, pS)
		}
		//validText := !pSp.inTag && !pS.inTag && !pS.isScript && (pS.isPre || !pS.isStyle && !pS.isSkip)
		//validText := !pS.inTag && !pS.isScript && (pS.isPre || !pS.isStyle && !pS.isSkip)
		if !pS.inTag {
			fmt.Print(string(cH))
		}
		cH, _, err = rDB.ReadRune()
	}
	for _, el := range els {
		fmt.Println(el)
	}

}

/*
func parse01(rDBuf *bufio.Reader) {
	//var dT []byte
	elements := make([]element, 0)
	pS := &parseState{false, false, false, false, false, false}
	cH, _, err := rDBuf.ReadRune()
	//fmt.Print(string(d[i]), " | ", tagName, " | ", inTag, isInPre(tagName, pS), isInScript(tagName, pS), "\n")
	for err == nil {
		elem := &element{}
		switch cH {
		case '<':
			if isValidTagStart(rDBuf) {
				cH, _, err = rDBuf.ReadRune()
				for err == nil && cH != '>' && cH != ' ' && !isValidTag(rDBuf) {
					elem.tagName = append(elem.tagName, byte(cH))
					cH, _, err = rDBuf.ReadRune()
				}
				for err == nil && cH != '>' && cH != ' ' {
					elem.tagName = append(elem.tagName, byte(cH))
					cH, _, err = rDBuf.ReadRune()
				}
				//fmt.Println(string(elem.tagName))
				if elem.tagName[0] == '/' {
					elem.tagName = elem.tagName[1:]
					elem.tInfo.close = true
				} else if elem.tagName[len(elem.tagName)-1] == '/' {
					elem.tInfo.selfclosing = true
				} else {
					elem.tInfo.open = true
				}
				elements = append(elements, *elem)
			}
		default:
			if !pS.isPre {
				switch cH {
				case ' ':
					if pS.isSpace {
						pS.isSkip = true
					}
					pS.isSpace = true
				case '\n', '\t', '\r', '\f':
					pS.isSkip = true
					pS.isSpace = false
				default:
					pS.isSkip = false
					pS.isSpace = false
				}
			}
		}
		if len(elem.tagName) > 0 {
			updateState(elem, pS)
		}
		cH, _, err = rDBuf.ReadRune()
		//fmt.Print("'", string(d[i]), "'", d[i], pS, "\n")
	}

	for _, elem := range elements {
		fmt.Println(string(elem.tagName))
	}

}
*/

type parseState struct {
	inTag, isPre, isScript, isStyle, isSpace, isSkip, isComment bool
}

type tagInfo struct {
	open        bool
	close       bool
	selfclosing bool
	comment     bool
}

type element struct {
	tag     tagType
	content string
}

type tagType struct {
	name string
	//attr map[string]string
	attr string
	info string
}

func (e element) parseElement(b []byte) *element {
	//fmt.Println(string(b))
	return &e
}

func (e element) parseElement01(b []byte) *element {
	//fmt.Println(string(b))
	l := len(b)
	if l > 1 {
		tagName := make([]byte, 0)
		comment := make([]byte, 0)
		tInfo := "opening"
		i := 0
		if b[i] == '/' {
			tInfo = "closing" //tagInfo{false, true, false, false}
			i++
		} else if b[l-1] == '/' {
			tInfo = "self closing" //tagInfo{false, false, true, false}
		}
		//fmt.Println(string(b[i : i+3]))
		if b[i] == '!' && i+3 < l {
			if string(b[i:i+3]) == "!--" {
				tagName = b[i : i+3]
				i = i + 3
				//fmt.Print(string(b), ",", l, ",", i, ",", string(b[i-2]), ",", string(b[i-1]), ",", string(b[i]))
				for ; i < l && !(b[i-2] == '-' && b[i-1] == '-' && b[i] == '>'); i++ {
					comment = append(comment, b[i])
				}
				e.tag.name = string("!--")
				e.content = string(b[3:])
				tInfo = "comment"
			}

		}
		if l-3 >= 0 {
			if b[l-3] == '-' && b[l-2] == '-' && b[l-1] == '>' {
				e.tag.name = string("!--")
				e.content = string(b[:l-3])
				tInfo = "comment"
			}
		}

		for ; i < l && b[i] != ' ' && b[i] != '>' && b[i] != '/'; i++ {
			tagName = append(tagName, b[i])
		}
		if len(tagName) > 0 {
			e.tag.name = string(tagName)
			e.tag.info = tInfo
		}
	} else {
		e.tag.name = string(b)
		e.tag.info = "opening"
	}

	//fmt.Println(string(e.tagName))
	return &e
}

func updateState(e element, pS *parseState) {
	ts := map[string]*bool{
		"script":  &pS.isScript,
		"pre":     &pS.isPre,
		"style":   &pS.isStyle,
		"comment": &pS.isComment,
	}
	if ts[e.tag.name] != nil {
		switch e.tag.info {
		case "opening":
			*(ts[e.tag.name]) = true
		case "closing":
			*(ts[e.tag.name]) = false
		case "comment":
			*(ts[e.tag.name]) = true
		}
	}
}

func strJoin(strA string, strB string) string {
	return strings.Join([]string{strA, strB}, "")
}

func isValidTagEnd(d []byte, pS *parseState) bool {
	//fmt.Println(string(d))
	l := len(d)
	if l > 1 {
		ecs := []bool{
			d[l-2] == '-' && d[l-1] == '-',
			d[l-2] == ';' && d[l-1] == '"',
			d[l-1] == '/',
			d[l-1] == ' ',
			isAlpha(d[l-1]),
			isNum(d[l-1]),
			d[l-1] == '"',
			d[l-1] == '-',
			d[l-1] == '?',
		}
		if ecs[0] {
			pS.isComment = false
		}
		for _, ec := range ecs {
			if ec {
				return true
			}
		}

	} else if l > 0 {
		return isAlpha(d[0])
	}
	return false
}
func isValidTagStart(rDBuf *bufio.Reader, pS *parseState) bool {
	d, err := rDBuf.Peek(3)
	if err == nil {
		scs := []bool{
			d[0] == '!' && d[1] == '-' && d[2] == '-',
			d[0] == '?',
			d[0] == '!',
			d[0] == '/' && isAlpha(d[1]),
			isAlpha(d[0]),
		}
		if scs[0] {
			pS.isComment = true
		}
		for _, sc := range scs {
			if sc {
				return true
			}
		}
	}
	return false
}

func isAlpha(d byte) bool {
	return d > 64 && d < 91 || d > 96 && d < 123
}

func isNum(d byte) bool {

	return d > 47 && d < 58
}
