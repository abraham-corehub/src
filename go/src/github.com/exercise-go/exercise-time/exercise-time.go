package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	dts := time.Now()
	//ts := dts.Format("Mon Jan 2 15:04:05 MST 2006")
	tD := dts.Format("20060102")
	tT := strings.Replace(dts.Format("15.04.05.000"), ".", "", 3)
	dTS := tD + tT
	fmt.Println(dTS)
}
