#!/bin/sh
app="html2text"
[ -d $WINDIR ] && app="$app.exe"

for i in ../tst/*.html; do
	$app "$i" -o "$i.html2text.txt"
done
