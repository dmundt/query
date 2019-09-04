// Copyright 2019 Daniel Mundt. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: MIT
//

package query

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

// equal compares two Queries for equality.
func (q *Query) equal(r *Query) bool {
	return q.String() == r.String()
}

// copy copies value e into a new slice containing a one element copy of e.
func copy(e T) []T {
	return []T{e}
}

// duplicate duplicates value e into a slice containing two element copies of e.
func duplicate(e T) []T {
	return []T{e, e}
}

// getSlice creates a new slice object [start, start+len) with step size step.
func getSlice(start, len, step int) []T {
	a := make([]T, len)
	for i := range a {
		a[i] = start + step*i
	}
	return a
}

// less return the comparison of values e1 and e2 as boolean value.
func less(e1, e2 T) bool {
	return e1.(int) < e2.(int)
}

// null discards any value e.
func null(e T) []T {
	return []T{}
}

// print prints any value e.
func print(e T) {
	fmt.Printf("%v\n", e)
}

// shuffle shuffles slice t with seed time.Now().UnixNano().
func shuffle(t []T) []T {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(t), func(i, j int) {
		t[i], t[j] = t[j], t[i]
	})
	return t
}

// sum sums value v and e.
func sum(v T, e T) interface{} {
	return v.(int) + e.(int)
}

// Span creates a new [begin, end) span object with step size 1.
func span(begin, end int) []T {
	if end < begin {
		return getSlice(begin, begin-end+1, -1)
	}
	return getSlice(begin, end-begin+1, 1)
}

// truth returns a boolean function with value b.
func truth(b bool) func(T) bool {
	return func(T) bool {
		return b
	}
}

func TestQuery_equal(t *testing.T) {
	type args struct {
		r *Query
	}
	tests := []struct {
		name string
		q    *Query
		args args
		want bool
	}{
		{"equal#1", From([]T{}), args{From([]T{})}, true},
		{"equal#2", From([]T{1}), args{From([]T{})}, false},
		{"equal#3", From([]T{}), args{From([]T{1})}, false},
		{"equal#4", From([]T{1}), args{From([]T{1})}, true},
		{"equal#5", From([]T{1, 2, 3}), args{From([]T{1, 2, 3})}, true},
		{"equal#6", From([]T{1, 2, 3}), args{From([]T{3, 2, 1})}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.equal(tt.args.r); got != tt.want {
				t.Errorf("Query.equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_duplicate(t *testing.T) {
	type args struct {
		e T
	}
	tests := []struct {
		name string
		args args
		want []T
	}{
		{"duplicate#1", args{}, []T{nil, nil}},
		{"duplicate#2", args{1}, []T{1, 1}},
		{"duplicate#3", args{"hello"}, []T{"hello", "hello"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := duplicate(tt.args.e); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("duplicate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getSlice(t *testing.T) {
	type args struct {
		start int
		len   int
		step  int
	}
	tests := []struct {
		name string
		args args
		want []T
	}{
		{"getslice#1", args{start: 1, len: 0, step: 1}, []T{}},
		{"getslice#2", args{start: 1, len: 1, step: 1}, []T{1}},
		{"getslice#3", args{start: 1, len: 9, step: 0}, []T{1, 1, 1, 1, 1, 1, 1, 1, 1}},
		{"getslice#4", args{start: 1, len: 9, step: 1}, []T{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{"getslice#5", args{start: 9, len: 9, step: -1}, []T{9, 8, 7, 6, 5, 4, 3, 2, 1}},
		{"getslice#5", args{start: 4, len: 9, step: -1}, []T{4, 3, 2, 1, 0, -1, -2, -3, -4}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSlice(tt.args.start, tt.args.len, tt.args.step); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_less(t *testing.T) {
	type args struct {
		t1 T
		t2 T
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"less#1", args{0, 0}, false},
		{"less#2", args{0, 1}, true},
		{"less#3", args{1, 0}, false},
		{"less#4", args{-1, 1}, true},
		{"less#5", args{1, -1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := less(tt.args.t1, tt.args.t2); got != tt.want {
				t.Errorf("less() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_null(t *testing.T) {
	type args struct {
		e T
	}
	tests := []struct {
		name string
		args args
		want []T
	}{
		{"null#1", args{}, []T{}},
		{"null#2", args{1}, []T{}},
		{"null#3", args{"hello"}, []T{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := null(tt.args.e); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("null() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_print(t *testing.T) {
	type args struct {
		t T
	}
	tests := []struct {
		name string
		args args
	}{
		{"print#1", args{}},
		{"print#2", args{0}},
		{"print#3", args{span(1, 10)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			print(tt.args.t)
		})
	}
}

func Test_shuffle(t *testing.T) {
	type args struct {
		t []T
	}
	tests := []struct {
		name string
		args args
		want []T
	}{
		// Test for inequality with shuffled sequence.
		// Entropy should be high enough using 1024 elements.
		{"shuffle#1", args{span(1, 1024)}, span(1, 1024)},
		{"shuffle#2", args{span(-1024, 1024)}, span(-1024, 1024)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shuffle(tt.args.t); reflect.DeepEqual(got, tt.want) {
				t.Errorf("shuffle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_span(t *testing.T) {
	type args struct {
		begin int
		end   int
	}
	tests := []struct {
		name string
		args args
		want []T
	}{
		{"span#1", args{}, []T{0}},
		{"span#2", args{begin: 1, end: 1}, []T{1}},
		{"span#3", args{begin: 1, end: 9}, []T{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{"span#4", args{begin: 9, end: 1}, []T{9, 8, 7, 6, 5, 4, 3, 2, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := span(tt.args.begin, tt.args.end); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("span() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sum(t *testing.T) {
	type args struct {
		v T
		e T
	}
	tests := []struct {
		name string
		args args
		want T
	}{
		{"sum#1", args{0, 0}, 0},
		{"sum#2", args{0, 1}, 1},
		{"sum#3", args{1, 1}, 2},
		{"sum#4", args{-1, 1}, 0},
		{"sum#5", args{-1, -1}, -2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sum(tt.args.v, tt.args.e); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_truth(t *testing.T) {
	type args struct {
		b bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"truth#1", args{b: false}, false},
		{"truth#2", args{b: true}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := truth(tt.args.b); !reflect.DeepEqual(got(t), tt.want) {
				t.Errorf("truth() = %v, want %v", got(t), tt.want)
			}
		})
	}
}

func TestQuery_Any(t *testing.T) {
	type args struct {
		f []func(T) bool
	}
	tests := []struct {
		name string
		q    *Query
		args args
		want bool
	}{
		{"any#1", From([]T{}), args{}, false},
		{"any#2", From([]T{}), args{[]func(T) bool{truth(true), truth(true)}}, false},
		{"any#3", From([]T{}), args{[]func(T) bool{truth(true), truth(false)}}, false},
		{"any#4", From([]T{}), args{[]func(T) bool{truth(false), truth(false)}}, false},
		{"any#5", From([]T{}), args{[]func(T) bool{truth(false), truth(true)}}, false},
		{"any#6", From(span(1, 9)), args{[]func(T) bool{truth(false), truth(false)}}, false},
		{"any#7", From(span(1, 9)), args{[]func(T) bool{truth(false), truth(true)}}, false},
		{"any#8", From(span(1, 9)), args{[]func(T) bool{truth(true), truth(true)}}, true},
		{"any#9", From(span(1, 9)), args{[]func(T) bool{truth(true), truth(false)}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.Any(tt.args.f...); got != tt.want {
				t.Errorf("Query.Any() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuery_At(t *testing.T) {
	type args struct {
		i int
	}
	tests := []struct {
		name string
		q    *Query
		args args
		want T
	}{
		{"at#1", From([]T{}), args{0}, nil},
		{"at#2", From(span(1, 9)), args{5}, 6},
		{"at#3", From(span(1, 9)), args{15}, nil},
		{"at#4", From(span(1, 9)), args{-100}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.At(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Query.At() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuery_Contains(t *testing.T) {
	type args struct {
		t T
	}
	tests := []struct {
		name string
		q    *Query
		args args
		want bool
	}{
		{"contains#1", From([]T{}), args{}, false},
		{"contains#2", From([]T{}), args{42}, false},
		{"contains#3", From(span(1, 9)), args{}, false},
		{"contains#4", From(span(1, 9)), args{5}, true},
		{"contains#5", From(span(1, 9)), args{10}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.Contains(tt.args.t); got != tt.want {
				t.Errorf("Query.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuery_Empty(t *testing.T) {
	tests := []struct {
		name string
		q    *Query
		want bool
	}{
		{"empty#1", From([]T{}), true},
		{"empty#2", From([]T{1}), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.IsEmpty(); got != tt.want {
				t.Errorf("Query.Empty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuery_Every(t *testing.T) {
	type args struct {
		f []func(T) bool
	}
	tests := []struct {
		name string
		q    *Query
		args args
		want bool
	}{
		{"every#1", From([]T{}), args{}, true},
		{"every#2", From([]T{}), args{[]func(T) bool{truth(true), truth(true)}}, true},
		{"every#3", From([]T{}), args{[]func(T) bool{truth(false), truth(false)}}, true},
		{"every#4", From(span(1, 9)), args{[]func(T) bool{truth(false), truth(false)}}, false},
		{"every#5", From(span(1, 9)), args{[]func(T) bool{truth(true), truth(true)}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.Every(tt.args.f...); got != tt.want {
				t.Errorf("Query.Every() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuery_First(t *testing.T) {
	tests := []struct {
		name string
		q    *Query
		want T
	}{
		{"first#1", From([]T{}), nil},
		{"first#2", From([]T{1}), 1},
		{"first#3", From(span(1, 9)), 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.First(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Query.First() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuery_Fold(t *testing.T) {
	type args struct {
		v T
		f func(t1, t2 T) interface{}
	}
	tests := []struct {
		name string
		q    *Query
		args args
		want interface{}
	}{
		{"fold#1", From([]T{}), args{}, nil},
		{"fold#2", From([]T{}), args{0, sum}, 0},
		{"fold#3", From([]T{}), args{10, sum}, 10},
		{"fold#4", From([]T{}), args{-10, sum}, -10},
		{"fold#5", From(span(1, 9)), args{0, sum}, 45},
		{"fold#6", From(span(1, 9)), args{10, sum}, 55},
		{"fold#7", From(span(1, 9)), args{-10, sum}, 35},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.Fold(tt.args.v, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Query.Fold() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuery_ForEach(t *testing.T) {
	type args struct {
		f func(T)
	}
	tests := []struct {
		name string
		q    *Query
		args args
	}{
		{"foreach#1", From([]T{}), args{}},
		{"foreach#2", From(span(1, 9)), args{print}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.q.ForEach(tt.args.f)
		})
	}
}

func TestQuery_Expand(t *testing.T) {
	type args struct {
		f func(e T) []T
	}
	tests := []struct {
		name string
		q    *Query
		args args
		want *Query
	}{
		{"expand#1", From([]T{}), args{null}, From([]T{})},
		{"expand#2", From([]T{1, 2, 3}), args{null}, From([]T{})},
		{"expand#3", From([]T{}), args{copy}, From([]T{})},
		{"expand#4", From([]T{1, 2, 3}), args{copy}, From([]T{1, 2, 3})},
		{"expand#5", From([]T{}), args{duplicate}, From([]T{})},
		{"expand#6", From([]T{1, 2, 3}), args{duplicate}, From([]T{1, 1, 2, 2, 3, 3})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.Expand(tt.args.f); !got.equal(tt.want) {
				t.Errorf("Query.Expand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFrom(t *testing.T) {
	type args struct {
		t []T
	}
	tests := []struct {
		name string
		args args
		want *Query
	}{
		{"from#1", args{[]T{}}, From([]T{})},
		{"from#2", args{span(1, 9)}, From(span(1, 9))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := From(tt.args.t); !got.equal(tt.want) {
				t.Errorf("From() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuery_Join(t *testing.T) {
	keySel := func(e T) interface{} {
		return e
	}
	resultSel := func(o, i interface{}) interface{} {
		return o
	}

	type args struct {
		inner     *Query
		outKeySel func(T) interface{}
		innKeySel func(T) interface{}
		resultSel func(o, i interface{}) interface{}
	}
	tests := []struct {
		name string
		q    *Query
		args args
		want *Query
	}{
		{"join#1", From([]T{}), args{From([]T{}), nil, nil, nil}, From([]T{})},
		{"join#2", From([]T{}), args{From([]T{}), keySel, keySel, resultSel}, From([]T{})},
		{"join#3", From(span(1, 9)), args{From([]T{}), keySel, keySel, resultSel}, From([]T{})},
		{"join#4", From(span(4, 9)), args{From(span(1, 6)), keySel, keySel, resultSel}, From(span(4, 6))},
		{"join#5", From(span(1, 6)), args{From(span(6, 9)), keySel, keySel, resultSel}, From([]T{6})},
		{"join#6", From(span(1, 5)), args{From(span(6, 9)), keySel, keySel, resultSel}, From([]T{})},
		{"join#7", From([]T{}), args{From(span(6, 9)), keySel, keySel, resultSel}, From([]T{})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.Join(tt.args.inner, tt.args.outKeySel, tt.args.innKeySel, tt.args.resultSel); !got.equal(tt.want) {
				t.Errorf("Query.Join() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuery_Last(t *testing.T) {
	tests := []struct {
		name string
		q    *Query
		want T
	}{
		{"last#1", From([]T{}), nil},
		{"last#2", From([]T{1}), 1},
		{"last#3", From(span(1, 9)), 9},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.Last(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Query.Last() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuery_MapTo(t *testing.T) {
	type args struct {
		f func(e T) T
	}
	tests := []struct {
		name string
		q    *Query
		args args
		want *Query
	}{
		{"mapto#1", From([]T{}), args{}, From([]T{})},
		{"mapto#2", From([]T{}), args{func(e T) T { return e.(int) + 10 }}, From([]T{})},
		{"mapto#3", From([]T{1, 2, 3, 4, 5}), args{func(e T) T { return e.(int) + 10 }}, From([]T{11, 12, 13, 14, 15})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.MapTo(tt.args.f); !got.equal(tt.want) {
				t.Errorf("Query.MapTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuery_Reduce(t *testing.T) {
	type args struct {
		f func(v T, e T) interface{}
	}
	tests := []struct {
		name string
		q    *Query
		args args
		want interface{}
	}{
		{"reduce#1", From([]T{}), args{}, nil},
		{"reduce#2", From([]T{}), args{sum}, nil},
		{"reduce#3", From([]T{1}), args{sum}, 1},
		{"reduce#4", From(span(1, 9)), args{sum}, 45},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.Reduce(tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Query.Reduce() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuery_Skip(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		q    *Query
		args args
		want *Query
	}{
		{"skip#1", From([]T{}), args{0}, From([]T{})},
		{"skip#2", From([]T{}), args{5}, From([]T{})},
		{"skip#3", From(span(1, 9)), args{0}, From(span(1, 9))},
		{"skip#4", From(span(1, 9)), args{5}, From(span(6, 9))},
		{"skip#5", From(span(1, 9)), args{8}, From([]T{9})},
		{"skip#6", From(span(1, 9)), args{9}, From([]T{})},
		{"skip#7", From(span(1, 9)), args{100}, From([]T{})},
		{"skip#8", From(span(1, 9)), args{-100}, From([]T{})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.Skip(tt.args.n); !got.equal(tt.want) {
				t.Errorf("Query.Skip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuery_Sort(t *testing.T) {
	type args struct {
		f []func(t1, t2 T) bool
	}
	tests := []struct {
		name string
		q    *Query
		args args
		want *Query
	}{
		{"sortby#1", From([]T{}), args{}, From([]T{})},
		{"sortby#2", From([]T{1}), args{[]func(t1, t2 T) bool{less, less}}, From([]T{1})},
		{"sortby#3", From(shuffle(span(1, 9))), args{[]func(t1, t2 T) bool{less, less}}, From(span(1, 9))},
		{"sortby#4", From(span(9, 1)), args{[]func(t1, t2 T) bool{less, less}}, From(span(1, 9))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.Sort(tt.args.f...); !got.equal(tt.want) {
				t.Errorf("Query.SortBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBy_Sort(t *testing.T) {
	type args struct {
		t []interface{}
	}
	tests := []struct {
		name string
		by   by
		args args
	}{
		{"sort#1", []func(t1, t2 T) bool{less, less}, args{[]interface{}{}}},
		{"sort#2", []func(t1, t2 T) bool{less, less}, args{[]interface{}{1}}},
		{"sort#3", []func(t1, t2 T) bool{less, less}, args{[]interface{}{1, 6, 3, 7, 5, 8, 4, 9, 2}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.by.Sort(tt.args.t)
		})
	}
}

func Test_sorter_Len(t *testing.T) {
	tests := []struct {
		name string
		s    sorter
		want int
	}{
		{"len#1", sorter{}, 0},
		{"len#2", sorter{t: []interface{}{1}}, 1},
		{"len#3", sorter{t: []interface{}{1, 6, 3, 7, 5, 8, 4, 9, 2}}, 9},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Len(); got != tt.want {
				t.Errorf("sorter.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sorter_Swap(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		s    sorter
		args args
	}{
		{"swap#1", sorter{t: []interface{}{1}}, args{0, 0}},
		{"swap#2", sorter{t: []interface{}{1, 2}}, args{0, 1}},
		{"swap#3", sorter{t: []interface{}{1, 2}}, args{1, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Swap(tt.args.i, tt.args.j)
		})
	}
}

func Test_sorter_Less(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		s    sorter
		args args
		want bool
	}{
		{"less#1", sorter{t: []interface{}{1}, less: []func(t1, t2 T) bool{less, less}}, args{0, 0}, false},
		{"less#2", sorter{t: []interface{}{1, 1}, less: []func(t1, t2 T) bool{less, less}}, args{0, 1}, false},
		{"less#3", sorter{t: []interface{}{2, 1}, less: []func(t1, t2 T) bool{less, less}}, args{0, 1}, false},
		{"less#4", sorter{t: []interface{}{1, 2}, less: []func(t1, t2 T) bool{less, less}}, args{0, 1}, true},
		{"less#5", sorter{t: []interface{}{2, 1}, less: []func(t1, t2 T) bool{less, less}}, args{1, 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Less(tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("sorter.Less() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuery_Take(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		q    *Query
		args args
		want *Query
	}{
		{"take#1", From([]T{}), args{0}, From([]T{})},
		{"take#2", From(span(1, 9)), args{0}, From([]T{})},
		{"take#3", From([]T{}), args{5}, From([]T{})},
		{"take#4", From(span(1, 9)), args{5}, From(span(1, 5))},
		{"take#5", From(span(1, 9)), args{9}, From(span(1, 9))},
		{"take#6", From(span(1, 9)), args{100}, From(span(1, 9))},
		{"take#7", From(span(1, 9)), args{-100}, From([]T{})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.Take(tt.args.n); !got.equal(tt.want) {
				t.Errorf("Query.Take() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToSlice(t *testing.T) {
	type args struct {
		q *Query
	}
	tests := []struct {
		name string
		args args
		want []interface{}
	}{
		{"toslice#1", args{From([]T{})}, []interface{}{}},
		{"toslice#2", args{From([]T{1})}, []interface{}{1}},
		{"toslice#3", args{From(span(1, 9))}, []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToSlice(tt.args.q); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuery_Where(t *testing.T) {
	type args struct {
		f []func(T) bool
	}
	tests := []struct {
		name string
		q    *Query
		args args
		want *Query
	}{
		{"where#1", From([]T{}), args{}, From([]T{})},
		{"where#2", From([]T{}), args{[]func(T) bool{truth(false), truth(false)}}, From([]T{})},
		{"where#3", From([]T{}), args{[]func(T) bool{truth(true), truth(false)}}, From([]T{})},
		{"where#4", From([]T{}), args{[]func(T) bool{truth(false), truth(true)}}, From([]T{})},
		{"where#5", From([]T{}), args{[]func(T) bool{truth(true), truth(true)}}, From([]T{})},
		{"where#6", From([]T{1, 2, 3}), args{[]func(T) bool{truth(false), truth(false)}}, From([]T{})},
		{"where#7", From([]T{1, 2, 3}), args{[]func(T) bool{truth(true), truth(false)}}, From([]T{})},
		{"where#8", From([]T{1, 2, 3}), args{[]func(T) bool{truth(false), truth(true)}}, From([]T{})},
		{"where#9", From([]T{1, 2, 3}), args{[]func(T) bool{truth(true), truth(true)}}, From([]T{1, 2, 3})},
		{"where#10", From(span(1, 9)),
			args{[]func(T) bool{func(e T) bool {
				return e.(int) < 4
			}}}, From([]T{1, 2, 3})},
		{"where#11", From(span(1, 9)),
			args{[]func(T) bool{func(e T) bool {
				return e.(int) > 9
			}}}, From([]T{})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.Where(tt.args.f...); !got.equal(tt.want) {
				t.Errorf("Query.Where() = %v, want %v", got, tt.want)
			}
		})
	}
}
