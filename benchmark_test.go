// Copyright 2019 Daniel Mundt. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: MIT
//

package query

import (
	"testing"
)

func BenchmarkQuery_Expand(b *testing.B) {
	a := shuffle(span(1, 100000))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		From(a).
			// Duplicate elements:
			Expand(func(e T) []T {
				return []T{e, e}
			}).
			// Pull the lazy iterator:
			ForEach(func(e T) {})
	}
}

func BenchmarkQuery_Join(b *testing.B) {
	inner := shuffle(span(1, 100000))
	outer := shuffle(span(1, 100000))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		From(outer).
			Join(From(inner),
				func(e T) interface{} {
					return e
				},
				func(e T) interface{} {
					return e
				},
				func(o, i interface{}) interface{} {
					return []T{o, i}
				}).
			// Pull the lazy iterator:
			ForEach(func(e T) {})
	}
}

func BenchmarkQuery_MapTo(b *testing.B) {
	a := shuffle(span(1, 100000))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		From(a).
			// Copy elements:
			MapTo(func(e T) T {
				return e
			}).
			// Pull the lazy iterator:
			ForEach(func(e T) {})
	}
}

func BenchmarkQuery_Sort(b *testing.B) {
	data := shuffle(span(1, 100000))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		From(data).
			// Sort random elements in ascending order:
			Sort(func(t1, t2 T) bool {
				return t1.(int) > t2.(int)
			}).
			// Pull the lazy iterator:
			ForEach(func(T) {})
	}
}
