// Copyright 2019 Daniel Mundt. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package query provides primitives for querying slices
// and user-defined collections.
package query

import (
	"fmt"
	"sort"
)

// T is an interface that has to be implemented by a custom collection in
// order to work with Query.
type T interface{}

// Iterator is an alias for function which iterates over slices.
type Iterator func() (elem T, ok bool)

// Query is the type returned from query functions. It can be iterated manually.
type Query struct {
	Iterate func() Iterator
}

// String converts the query to a string.
func (q *Query) String() string {
	return fmt.Sprintf("%v", ToSlice(q))
}

// Any checks whether any element of this collection satisfies all predicates.
//
// Checks every element in iteration order, and returns true
// if any of them make test return true, otherwise returns false.
func (q *Query) Any(f ...func(e T) bool) bool {
	next := q.Iterate()
	for elem, ok := next(); ok; elem, ok = next() {
		has := true
		for k := 0; k < len(f); k++ {
			has = has && f[k](elem)
		}
		if has {
			return true
		}
	}
	return false
}

// At returns the ith element.
//
// The index i must be non-negative and less than length.
// Index zero represents the first element (so Query.At(0) is equivalent to Query.First()).
//
// May iterate through the elements in iteration order,
// ignoring the first i elements and then returning the next.
func (q *Query) At(i int) (elem T) {
	if i < 0 {
		return nil
	}
	next := q.Iterate()
	for ; i >= 0; i-- {
		elem, _ = next()
	}
	return
}

// Contains returns true if the collection contains an element equal to element.
// This operation will check each element in order for being equal to element,
// unless it has a more efficient way to find an element equal to element.
func (q *Query) Contains(e T) bool {
	next := q.Iterate()
	for elem, ok := next(); ok; elem, ok = next() {
		if elem == e {
			return true
		}
	}
	return false
}

// IsEmpty returns true if there are no elements in this collection.
func (q *Query) IsEmpty() bool {
	next := q.Iterate()
	_, ok := next()
	return !ok
}

// Every checks whether every element of this collection satisfies all tests.
// Checks every element in iteration order, and returns false
// if any of them make test return false, otherwise returns true.
func (q *Query) Every(f ...func(e T) bool) bool {
	next := q.Iterate()
	has := true
	for elem, ok := next(); ok; elem, ok = next() {
		for k := 0; k < len(f); k++ {
			has = has && f[k](elem)
		}
	}
	return has
}

// Expand expands each element of this Query into zero or more elements.
//
// The resulting Query runs through the elements returned by f
// for each element of this, in iteration order.
//
// The returned Query is lazy, and calls f
// for each element of this every time it's iterated.
func (q *Query) Expand(f func(e T) []T) *Query {
	iterate := func() Iterator {
		return expand(q, f)
	}
	return &Query{iterate}
}

type expState struct {
	outer T
	inner []T
	i     int
	len   int
}

func expand(q *Query, f func(e T) []T) Iterator {
	next := q.Iterate()
	s := expState{}

	return func() (elem T, ok bool) {
		for {
			if s.i >= s.len {
				s.outer, ok = next()
				if !ok {
					return
				}
				s.inner = f(s.outer)
				s.len = len(s.inner)
				s.i = 0
			}

			if s.i < s.len {
				elem = s.inner[s.i]
				s.i++
				return elem, true
			}
		}
	}
}

// First returns the first element.
func (q *Query) First() (first T) {
	next := q.Iterate()
	first, _ = next()
	return
}

// Fold reduces a collection to a single value by iteratively combining
// each element of the collection with an existing value.
//
// Uses v as the initial value, then iterates through the elements
// and updates the value with each element using the combine function.
func (q *Query) Fold(v T, f func(v, e T) interface{}) interface{} {
	next := q.Iterate()
	for elem, ok := next(); ok; elem, ok = next() {
		v = f(v, elem)
	}
	return v
}

// ForEach applies the function f to each element of this collection in iteration order.
func (q *Query) ForEach(f func(e T)) {
	next := q.Iterate()
	for elem, ok := next(); ok; elem, ok = next() {
		f(elem)
	}
}

// From initializes a query with passed slice as the source.
func From(a []T) *Query {
	iterate := func() Iterator {
		return from(a)
	}
	return &Query{iterate}
}

func from(a []T) Iterator {
	i := 0
	return func() (elem T, ok bool) {
		ok = i < len(a)
		if ok {
			elem = a[i]
			i++
		}
		return
	}
}

// Join correlates the elements of two collection based on matching keys.
//
// A join refers to the operation of correlating the elements of two sources of
// information based on a common key. Join brings the two information sources
// and the keys by which they are matched together in one method call.
//
// Join preserves the order of the elements of outer collection, and for each of
// these elements, the order of the matching elements of inner.
func (q *Query) Join(inner *Query,
	outKeySel func(e T) interface{},
	innKeySel func(e T) interface{},
	resultSel func(o, i interface{}) interface{}) *Query {
	iterate := func() Iterator {
		return join(q, inner, outKeySel, innKeySel, resultSel)
	}
	return &Query{iterate}
}

type lut map[T][]T

func makeLut(it Iterator, f func(e T) interface{}) (result lut) {
	next := it
	result = make(lut)

	for elem, ok := next(); ok; elem, ok = next() {
		key := f(elem)
		result[key] = append(result[key], elem)
	}
	return
}

type joinState struct {
	outer T
	inner []T
	i     int
	len   int
}

func join(q *Query, inner *Query,
	outKeySel func(e T) interface{},
	innKeySel func(e T) interface{},
	resultSel func(o, i interface{}) interface{}) Iterator {
	next := q.Iterate()
	lut := makeLut(inner.Iterate(), innKeySel)
	s := joinState{}

	return func() (elem T, ok bool) {
		if s.i >= s.len {
			has := false
			for !has {
				s.outer, ok = next()
				if !ok {
					return
				}
				s.inner, has = lut[outKeySel(s.outer)]
				s.len = len(s.inner)
				s.i = 0
			}
		}
		elem = resultSel(s.outer, s.inner[s.i])
		s.i++
		return elem, true
	}
}

// Last returns the last element.
func (q *Query) Last() (last T) {
	next := q.Iterate()
	for elem, ok := next(); ok; elem, ok = next() {
		last = elem
	}
	return
}

// MapTo returns a new lazy Query with elements that are created by
// calling f on each element of this Query in iteration order.
//
// This method returns a view of the mapped elements.
// As long as the returned Iterable is not iterated over,
// the supplied function f will not be invoked.
// The transformed elements will not be cached.
// Iterating multiple times over the returned Iterable will invoke
// the supplied function f multiple times on the same element.
//
// Methods on the returned query are allowed to omit
// calling f on any element where the result isn't needed.
func (q *Query) MapTo(f func(e T) T) *Query {
	iterate := func() Iterator {
		return mapTo(q, f)
	}
	return &Query{iterate}
}

func mapTo(q *Query, f func(e T) T) Iterator {
	next := q.Iterate()
	return func() (elem T, ok bool) {
		elem, ok = next()
		if ok {
			return f(elem.(T)), ok
		}
		return
	}
}

// Reduce reduces a collection to a single value by iteratively combining
// elements of the collection using the provided function.
//
// The iterable must have at least one element.
// If it has only one element, that element is returned.
//
// Otherwise this method starts with the first element from the iterator,
// and then combines it with the remaining elements in iteration order, as if by:
//  v := q.First()
//  q.Skip(1).ForEach(func(e T) {
//  	v = comb(v, e)
//  })
//  return v
func (q *Query) Reduce(f func(v, e T) interface{}) interface{} {
	next := q.Iterate()
	if v, ok := next(); ok {
		for elem, ok := next(); ok; elem, ok = next() {
			v = f(v, elem)
		}
		return v
	}
	return nil
}

// Skip returns an Query that provides all but the first n elements.
//
// When the returned query is iterated, it starts iterating over this,
// first skipping past the initial n elements. If this has fewer than n elements,
// then the resulting Query is empty. After that, the remaining elements are iterated
// in the same order as in this query.
//
// The n must not be negative.
func (q *Query) Skip(n int) *Query {
	iterate := func() Iterator {
		return skip(q, n)
	}
	return &Query{iterate}
}

func skip(q *Query, n int) Iterator {
	next := q.Iterate()
	return func() (elem T, ok bool) {
		if n < 0 {
			return
		}
		for ; n > 0; n-- {
			_, ok = next()
			if !ok {
				return
			}
		}
		return next()
	}
}

// Sort sorts the elements of a collection in predicate order.
// Elements are sorted according to a key while keeping
// the original order of equal elements.
func (q *Query) Sort(f ...func(e, f T) bool) *Query {
	iterate := func() Iterator {
		return sortBy(q, f)
	}
	return &Query{iterate}
}

func sortBy(q *Query, f []func(e, f T) bool) Iterator {
	a := ToSlice(q)
	by(f).Sort(a)

	i := 0
	return func() (elem T, ok bool) {
		ok = i < len(a)
		if ok {
			elem = a[i]
			i++
		}
		return
	}
}

// by is the type of a "less" function array that defines the ordering of its arguments.
type by []func(e, j T) bool

// Sort is a method on the function type, by, that sorts the collection according to the function array.
func (f by) Sort(t []interface{}) {
	s := &sorter{
		t:    t,
		less: f, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Stable(s)
}

// sorter joins a by function and a slice of its elements to be sorted.
type sorter struct {
	t    []interface{}
	less by // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *sorter) Len() int {
	return len(s.t)
}

// Swap is part of sort.Interface.
func (s *sorter) Swap(i, j int) {
	s.t[i], s.t[j] = s.t[j], s.t[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *sorter) Less(i, j int) bool {
	// Try all but the last comparison.
	var k int
	for k = 0; k < len(s.less)-1; k++ {
		less := s.less[k]
		switch {
		case less(s.t[i], s.t[j]):
			// s.t[i] < s.t[j], so we have a decision.
			return true
		case less(s.t[j], s.t[i]):
			// s.t[i] > s.t[j], so we have a decision.
			return false
		}
		// s.t[i] == s.t[j]; try the next comparison.
	}
	// All comparisons to here said "equal", so just return whatever
	// the final comparison reports.
	return s.less[k](s.t[i], s.t[j])
}

// Take returns a lazy query of the n first elements of this query.
//
// The returned Query may contain fewer than n elements,
// if this contains fewer than n elements.
//
// The elements can be computed by stepping through iterator
// until n elements have been seen.
//
// The n must not be negative.
func (q *Query) Take(n int) *Query {
	iterate := func() Iterator {
		return take(q, n)
	}
	return &Query{iterate}
}

func take(q *Query, n int) Iterator {
	next := q.Iterate()
	return func() (elem T, ok bool) {
		if n <= 0 {
			return
		}
		n--
		return next()
	}
}

// ToSlice iterates over a collection and saves the results in the slice pointed
// by v. It overwrites the existing slice, starting from index 0.
func ToSlice(q *Query) []interface{} {
	a := []interface{}{}
	next := q.Iterate()
	for elem, ok := next(); ok; elem, ok = next() {
		a = append(a, elem)
	}
	return a
}

// Where returns a new lazy Query with all elements that satisfy all predicate tests.
//
// The matching elements have the same order in the returned iterable as they have in iterator.
//
// This function returns a view of the mapped elements. As long as the returned Iterable
// is not iterated over, the supplied function test will not be invoked.
// Iterating will not cache results, and thus iterating multiple times over the returned
// Query may invoke the supplied function test multiple times on the same element.
func (q *Query) Where(f ...func(e T) bool) *Query {
	iterate := func() Iterator {
		return where(q, f)
	}
	return &Query{iterate}
}

// where returns a new lazy iterator with all elements that satisfy all predicate tests.
func where(q *Query, f []func(e T) bool) Iterator {
	next := q.Iterate()
	return func() (elem T, ok bool) {
		for elem, ok = next(); ok; elem, ok = next() {
			has := true
			for k := 0; k < len(f); k++ {
				has = has && f[k](elem)
			}
			if has {
				return
			}
		}
		return
	}
}
