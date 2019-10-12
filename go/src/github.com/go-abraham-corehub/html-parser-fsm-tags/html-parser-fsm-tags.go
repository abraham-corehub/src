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

	nameFileInput := strJoin(strJoin(nameDir, "/"), "The Ugly.html")
	f, err := os.Open(nameFileInput)
	if err != nil {
		fmt.Println(err)
		return
	}
	rDBuf := *bufio.NewReader(f)

	getTag(&rDBuf, "tag")
}

func getTag(rDB *bufio.Reader, tagName string) {
	var fsm abFSM
	fsm = fsm.createStateTable("stateTable.txt") // State Machine is created from the State Table text file
	cS := "initial"
	lS := "initial"
	cH, _, err := rDB.ReadRune()
	for err == nil {
		cS = getState(fsm, lS, byte(cH))
		text := tagName
		if cS == text {
			fmt.Print(string(cH))
		}
		if cS != lS && lS == text {
			fmt.Println("")
		}
		lS = cS
		cH, _, err = rDB.ReadRune()
	}
}

// abFSM holds the State Transition Table and
// State Assignment Mappings which defines the State Machine
type abFSM struct {
	mSTD  map[cSIn]byte   // Maps current states to next states for predefined inputs
	mSTX  map[byte]byte   // Maps current states to next states without considering inputs
	mS2ID map[string]byte // Maps State Names to State IDs
	mID2S map[byte]string // Maps State IDs to State Names
}

type cSIn struct {
	cS byte
	in byte
}

func (fsm abFSM) createStateTable(nameFileStateTable string) abFSM {
	//https://programming.guide/go/read-file-line-by-line.html
	file, err := os.Open(nameFileStateTable)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	table := make([][]string, 0)
	col := make([]string, 0)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		//https://golang.org/pkg/strings/#Split
		row := strings.Split(line, ",")
		col = append(col, []string{row[0], row[2]}...)
		table = append(table, row)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	//fmt.Print(table)
	states := unique(col)
	//fmt.Println(states)
	fsm.mS2ID = make(map[string]byte)
	fsm.mID2S = make(map[byte]string)
	for i, state := range states {
		fsm.mS2ID[state] = byte(i)
		fsm.mID2S[byte(i)] = state
	}

	//fmt.Println(fsm.mID2S)
	//fmt.Println(fsm.mS2ID)
	fsm.mSTD = make(map[cSIn]byte)
	fsm.mSTX = make(map[byte]byte)

	//fmt.Println(fsm.mS2ID)

	for _, row := range table {
		if row[1] != "0" {
			fsm.mSTD[cSIn{fsm.mS2ID[row[0]], strChar2Byte(row[1])}] = fsm.mS2ID[row[2]]
		} else {
			fsm.mSTX[fsm.mS2ID[row[0]]] = fsm.mS2ID[row[2]]
		}
	}
	//showMapTable(fsm, fsm.mSTD)
	//fmt.Println(fsm.mID2S, fsm.mSTD, fsm.mSTX)
	return fsm
}

func showMapTable(fsm abFSM, m map[cSIn]byte) {
	for a, v := range m {
		fmt.Println(a.cS, string(a.in), v)
	}
}

func strChar2Byte(strChar string) byte {
	var m map[string]byte
	if len(strChar) > 1 {
		m = map[string]byte{"\\n": 10, "\\t": 9, "\\b": 8, "\\f": 12, "\\r": 13}
	} else {
		m = map[string]byte{strChar: byte(strChar[0])}
	}
	return m[strChar]
}

func getState(fsm abFSM, stateName string, inCh byte) string {
	nS, prs := fsm.mSTD[cSIn{fsm.mS2ID[stateName], inCh}]
	if prs {
		stateName = fsm.mID2S[nS]
	} else {
		nS, prs = fsm.mSTX[fsm.mS2ID[stateName]]
		if prs {
			stateName = fsm.mID2S[nS]
		}
	}
	return stateName
}

func strJoin(strA string, strB string) string {
	return strings.Join([]string{strA, strB}, "")
}

//unique https://www.golangprograms.com/remove-duplicate-values-from-slice.html
func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
