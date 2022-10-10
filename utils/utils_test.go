package utils

import (
	"fmt"
	"reflect"
	"testing"
)

func TestEmpty(t *testing.T) {
	testMap := make(map[int]int)
	testSlice := make([]int, 0, 1)
	var testPtr *int
	var TestInterface interface{}
	var testArray [0]int

	var complex64TestTrue complex64
	complex64TestFalse := 2
	var complex128TestTrue complex128
	complex128TestFalse := 2

	tests := []struct {
		name string
		in   interface{}
		want bool
	}{
		{name: "string want true", in: "", want: true},
		{name: "string want false", in: "11", want: false},
		{name: "array want true", in: testArray, want: true},
		{name: "array want false", in: [2]int{1, 2}, want: false},
		{name: "Map want true", in: map[int]int{}, want: true},
		{name: "Map want true", in: testMap, want: true},
		{name: "Map want false", in: map[int]int{1: 1}, want: false},
		{name: "slice want true", in: testSlice, want: true},
		{name: "slice want true", in: []int{}, want: true},
		{name: "slice want false", in: []int{1, 2}, want: false},
		{name: "bool want true", in: false, want: true},
		{name: "bool want false", in: true, want: false},
		{name: "Int want true", in: 0, want: true},
		{name: "Int want false", in: -1, want: false},
		{name: "Int8 want true", in: 0, want: true},
		{name: "Int8 want false", in: 1, want: false},
		{name: "Int16 want true", in: 0, want: true},
		{name: "Int16 want false", in: -1, want: false},
		{name: "Int32 want true", in: 0, want: true},
		{name: "Int32 want false", in: 1, want: false},
		{name: "Int64 want true", in: 0, want: true},
		{name: "Int64 want false", in: -1, want: false},
		{name: "Uint want true", in: uint(0), want: true},
		{name: "Uint want false", in: uint(1), want: false},
		{name: "Complex64 want true", in: complex64TestTrue, want: true},
		{name: "Complex64 want false", in: complex64TestFalse, want: false},
		{name: "Complex128 want true", in: complex128TestTrue, want: true},
		{name: "Complex128 want false", in: complex128TestFalse, want: false},
		{name: "Float32 want true", in: 0.0, want: true},
		{name: "Float32 want false", in: 1.1, want: false},
		{name: "Interface want false", in: &TestInterface, want: false},
		{name: "Interface want false", in: 1.1, want: false},
		{name: "Ptr want true", in: testPtr, want: true},
		{name: "Ptr want false", in: &testMap, want: false},
	}
	for _, v := range tests {
		got := Empty(v.in)
		if got != v.want {
			t.Errorf("%s test failed, expected : %v, got : %v", v.name, v.want, got)
		}
	}
}

func TestInterfaceToString(t *testing.T) {
	tests := map[string]struct {
		input interface{}
		want  string
	}{
		"string":       {input: "abc", want: "abc"},
		"empty string": {input: "", want: ""},
		"int":          {input: 123, want: "123"},
		"bool":         {input: false, want: "false"},
		"others":       {input: []string{"a", "b"}, want: "[\"a\",\"b\"]"},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := InterfaceToString(tc.input)
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func TestZip(t *testing.T) {
	var zipPath = "../static/sdk/src"
	// 目标文件，压缩后的文件
	var dst = fmt.Sprintf("../static/sdk/dst/sdk.zip")
	if err := ZipFiles(dst, GetDicAllChildFilePath(zipPath)); err != nil {
		t.Log(err)
	}
}
