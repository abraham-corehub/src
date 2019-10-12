package main

import (
	"reflect"
	"testing"
)

func Test_strChar2Byte(t *testing.T) {
	type args struct {
		strChar string
	}
	tests := []struct {
		name string
		args args
		want byte
	}{
		// TODO: Add test cases.
		{"case1", args{"hello"}, byte('h')},
		{"case2", args{"\n"}, byte('\n')},
		{"case3", args{"\r"}, byte('\r')},
		{"case4", args{"\f"}, byte('\f')},
		{"case5", args{" "}, byte(' ')},
		{"case6", args{"\t"}, byte('\t')},
		{"case7", args{"\""}, byte('"')},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := strChar2Byte(tt.args.strChar); got != tt.want {
				t.Errorf("strChar2Byte() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stripTagsHTML(t *testing.T) {
	type args struct {
		dat []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stripTagsHTML(tt.args.dat); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stripTagsHTML() = %v, want %v", got, tt.want)
			}
		})
	}
}
