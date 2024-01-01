package repository

import (
	"reflect"
	"testing"
)

func TestChunksWithInt(t *testing.T) {
	cases := []struct {
		name string
		s    []int
		size int
		want [][]int
	}{
		{name: "normal1", s: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, size: 3, want: [][]int{{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, {9}}},
		{name: "normal2", s: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, size: 5, want: [][]int{{0, 1, 2, 3, 4}, {5, 6, 7, 8, 9}}},
		{name: "same", s: []int{0, 1, 2, 3, 4}, size: 5, want: [][]int{{0, 1, 2, 3, 4}}},
		{name: "length_is_zero", s: []int{}, size: 10, want: [][]int{}},
		{name: "size_is_zero", s: []int{0, 1, 2}, size: 0, want: nil},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := Chunks(tt.s, tt.size)
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("want = %+v, but got = %+v", tt.want, result)
			}
		})
	}
}

type data struct {
	ID string
}

func TestChunksWithStruct(t *testing.T) {
	cases := []struct {
		name string
		s    []data
		size int
		want [][]data
	}{
		{name: "normal1", s: []data{{ID: "001"}, {ID: "002"}, {ID: "003"}, {ID: "004"}, {ID: "005"}}, size: 3, want: [][]data{{{ID: "001"}, {ID: "002"}, {ID: "003"}}, {{ID: "004"}, {ID: "005"}}}},
		{name: "normal2", s: []data{{ID: "001"}, {ID: "002"}, {ID: "003"}, {ID: "004"}}, size: 2, want: [][]data{{{ID: "001"}, {ID: "002"}}, {{ID: "003"}, {ID: "004"}}}},
		{name: "same", s: []data{{ID: "001"}, {ID: "002"}, {ID: "003"}, {ID: "004"}, {ID: "005"}}, size: 5, want: [][]data{{{ID: "001"}, {ID: "002"}, {ID: "003"}, {ID: "004"}, {ID: "005"}}}},
		{name: "length_is_zero", s: []data{}, size: 10, want: [][]data{}},
		{name: "size_is_zero", s: []data{{"001"}, {"002"}}, size: 0, want: nil},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := Chunks(tt.s, tt.size)
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("want = %+v, but got = %+v", tt.want, result)
			}
		})
	}
}
