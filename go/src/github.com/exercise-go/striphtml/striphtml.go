package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
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
	parse(dat)

	//fmt.Println(formatText("hello\n\n\n     \n\n\n\n        \n\n\n\n         \n\n\nWorld\n\n\n\nhi"))
}

func parse(dat []byte) {
	//iSp := 0 // index Start previous
	iSc := 0 // index Start current
	iEp := 0 // index End previous
	iEc := 0 // index End current
	tags := make([][]byte, 0)
	//tagLsS := make([]int, 0) // tag Locations Start
	//tagLsE := make([]int, 0) // tag Locations End
	text := make([]byte, 0)

	for index, ch := range dat {
		switch ch {
		case '<':
			//iSp = iSc
			iSc = index
		case '>':
			iEp = iEc
			iEc = index
			tag := dat[iSc : iEc+1]
			//fmt.Println("tag : ", string(tag))
			tags = append(tags, tag)
			/*
				if bytes.Equal(tag, []byte("<hr>")) {
					text = bytesJoin(text, []byte("\n"))
					for i := 0; i < 80; i++ {
						text = bytesJoin(text, []byte("_"))
					}
					text = bytesJoin(text, []byte("\n"))
				}
			*/
			if iSc > iEp+1 {
				text = dat[iEp+1 : iSc]
				fmt.Println(formatText(string(text)))
			}

			/*
				tagLsS = append(tagLsS, iS)
				tagLsE = append(tagLsE, iE)
				lenTLS := len(tagLsS)
				if lenTLS > 1 && tagLsS[lenTLS-1]-tagLsE[lenTLS-2]+1 >= 0 {

					text = bytesJoin(text, dat[tagLsE[lenTLS-2]+1:tagLsS[lenTLS-1]])
				}
			*/
		}
	}
}

func formatText(text string) string {

	outText := text
	for _, sB := range []byte{' ', '\n', '\t', '\r'} {
		outText = stripByte(outText, sB)
	}

	return outText
}

func stripByte(text string, sB byte) string {

	outText := make([]byte, 0)
	count := 0
	for _, cH := range text {
		if sB == byte(cH) {
			if count < 1 {
				outText = append(outText, byte(cH))
			}
			count++
		} else {
			outText = append(outText, byte(cH))
			count = 0
		}
	}
	return string(outText)
}

func strJoin(strA string, strB string) string {
	return strings.Join([]string{strA, strB}, "")
}

func bytesJoin(byteA []byte, byteB []byte) []byte {
	return bytes.Join([][]byte{byteA, byteB}, []byte(""))
}

func categorizeTags(tag []byte) []byte {

	return tag

}

func extractTagName(tag []byte) []byte {
	iS := 0
	for i, ch := range tag {
		switch ch {
		case '<':
			iS = i
		case ' ':
			return tag[iS+1 : i]
		case '>':
			return tag[iS+1 : i]
		case '/':
			if tag[i+1] == '>' {
				return tag[iS+1 : i]
			}
			return tag[i+1 : len(tag)-1]
		}
	}

	return tag
}
