// Copyright 2019 Daniel Mundt. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package query

import (
	"fmt"
)

func ExampleFrom() {
	q := From([]T{1, 3, 5, 7, 9})
	fmt.Printf("Got from query: %v\n", q)

	// Output:
	// Got from query: [1 3 5 7 9]
}

func ExampleQuery_Any_odd() {
	isEven := func(e T) bool {
		return e.(int)&1 == 0
	}
	greaterThan0 := func(e T) bool {
		return e.(int) > 0
	}
	q := From([]T{1, 3, 5, 7, 9})
	v := q.Any(isEven, greaterThan0)
	fmt.Printf("An even number > 3 was found: %v\n", v)

	// Output:
	// An even number > 3 was found: false
}

func ExampleQuery_Any_greatherThan() {
	isEven := func(e T) bool {
		return e.(int)&1 == 0
	}
	greaterThan0 := func(e T) bool {
		return e.(int) > 0
	}
	q := From([]T{1, 3, 6, 7, 9})
	v := q.Any(isEven, greaterThan0)
	fmt.Printf("An even number > 3 was found: %v\n", v)

	// Output:
	// An even number > 3 was found: true
}

func ExampleQuery_At_found() {
	v := From([]T{1, 2, 3, 4, 5}).At(3)
	fmt.Printf("Element at index 5: %v\n", v)

	// Output:
	// Element at index 5: 4
}

func ExampleQuery_At_notFound() {
	v := From([]T{1, 2, 3, 4, 5}).At(15)
	fmt.Printf("Element at index 15: %v\n", v)

	// Output:
	// Element at index 15: <nil>
}

func ExampleQuery_Contains_notFound() {
	v := From([]T{1, 2, 3, 4, 5}).Contains(3)
	fmt.Printf("Contains 6: %v\n", v)

	// Output:
	// Contains 6: true
}

func ExampleQuery_Contains_found() {
	v := From([]T{1, 2, 3, 4, 5}).Contains(42)
	fmt.Printf("Contains 12: %v\n", v)

	// Output:
	// Contains 12: false
}

func ExampleQuery_Every_allOdd() {
	q := From([]T{1, 3, 5, 7, 9})
	v := q.Every(func(e T) bool {
		return e.(int)%2 != 0
	})
	fmt.Printf("All numbers are odd: %v\n", v)

	// Output:
	// All numbers are odd: true
}

func ExampleQuery_Every_allEven() {
	q := From([]T{1, 3, 5, 7, 9})
	v := q.Every(func(e T) bool {
		return e.(int)%2 != 1
	})
	fmt.Printf("All numbers are even: %v\n", v)

	// Output:
	// All numbers are even: false
}

func ExampleQuery_Expand_null() {
	q := From([]T{1, 2, 3, 4, 5})
	v := q.Expand(func(e T) []T {
		return []T{}
	})
	fmt.Printf("All elements are nulled: %v\n", v)

	// Output:
	// All elements are nulled: []
}

func ExampleQuery_Expand_identity() {
	q := From([]T{1, 2, 3, 4, 5})
	v := q.Expand(func(e T) []T {
		return []T{e}
	})
	fmt.Printf("All elements are identical: %v\n", v)

	// Output:
	// All elements are identical: [1 2 3 4 5]
}

func ExampleQuery_Expand_duplicate() {
	q := From([]T{1, 2, 3, 4, 5})
	v := q.Expand(func(e T) []T {
		return []T{e, e}
	})
	fmt.Printf("All elements are duplicated: %v\n", v)

	// Output:
	// All elements are duplicated: [1 1 2 2 3 3 4 4 5 5]
}

func ExampleQuery_Expand_odd() {
	// Remove even numbers:
	q := From([]T{1, 2, 3, 4, 5})
	v := q.Expand(func(e T) []T {
		if e.(int)%2 == 0 {
			return []T{}
		}
		return []T{e}
	})
	fmt.Printf("All elements are odd: %v\n", v)

	// Output:
	// All elements are odd: [1 3 5]
}

func ExampleQuery_First_found() {
	v := From([]T{1, 2, 3, 4, 5}).First()
	fmt.Printf("First element: %v", v)

	// Output:
	// First element: 1
}

func ExampleQuery_First_empty() {
	v := From([]T{}).First()
	fmt.Printf("First element: %v", v)

	// Output:
	// First element: <nil>
}

func ExampleQuery_Fold_sum() {
	// Calculating the sum of an query:
	sum := func(v, e T) interface{} {
		return v.(int) + e.(int)
	}
	v := From([]T{1, 2, 3}).Fold(0, sum)
	fmt.Printf("Folded elements to sum: %v", v)

	// Output:
	// Folded elements to sum: 6
}

func ExampleQuery_ForEach_append() {
	v := []T{}
	From([]T{1, 3, 5, 7, 9}).
		ForEach(func(e T) {
			v = append(v, e)
		})
	fmt.Printf("For each: %v", v)

	// Output:
	// For each: [1 3 5 7 9]
}

func ExampleQuery_ForEach_count() {
	v := 0
	From([]T{1, 3, 5, 7, 9}).
		ForEach(func(e T) {
			v++
		})
	fmt.Printf("For each: %v", v)

	// Output:
	// For each: 5
}

func ExampleQuery_ForEach_sum() {
	v := 0
	From([]T{1, 2, 3}).
		ForEach(func(e T) {
			v = v + e.(int)
		})
	fmt.Printf("For each: %v", v)

	// Output:
	// For each: 6
}

func ExampleQuery_IsEmpty_empty() {
	v := From([]T{}).IsEmpty()
	fmt.Printf("Empty query: %v\n", v)

	// Output:
	// Empty query: true
}

func ExampleQuery_IsEmpty_notEmpty() {
	v := From([]T{1}).IsEmpty()
	fmt.Printf("Empty query: %v\n", v)

	// Output:
	// Empty query: false
}

func ExampleQuery_Join_inner() {
	v := From([]T{1, 2, 3, 4, 5}).
		Join(From([]T{3, 4, 5, 6, 7}),
			// Outer key selector:
			func(e T) interface{} {
				return e
			},
			// Inner key selector:
			func(e T) interface{} {
				return e
			},
			// Result selector:
			func(o, i interface{}) interface{} {
				return []T{o, i}
			})
	fmt.Printf("Inner join: %v\n", v)

	// Output:
	// Inner join: [[3 3] [4 4] [5 5]]
}

func ExampleQuery_Last_found() {
	v := From([]T{1, 2, 3, 4, 5}).Last()
	fmt.Printf("Last element: %v", v)

	// Output:
	// Last element: 5
}

func ExampleQuery_Last_empty() {
	v := From([]T{}).Last()
	fmt.Printf("Last element: %v", v)

	// Output:
	// Last element: <nil>
}

func ExampleQuery_MapTo_add() {
	// Add a number to every slice collection element:
	add := func(e T) T {
		return e.(int) + 10
	}
	q := From([]T{1, 2, 3, 4, 5})
	v := q.MapTo(add)
	fmt.Printf("Map q to v: %v\n", v)

	// Output:
	// Map q to v: [11 12 13 14 15]
}

func ExampleQuery_MapTo_addToEven() {
	// Add a number to every even slice collection element:
	add := func(e T) T {
		if e.(int)%2 == 0 {
			return e.(int) + 10
		}
		return e.(int)
	}
	q := From([]T{1, 2, 3, 4, 5})
	v := q.MapTo(add)
	fmt.Printf("Map q to v: %v\n", v)

	// Output:
	// Map q to v: [1 12 3 14 5]
}

func ExampleQuery_Reduce_sum() {
	// Calculating the sum of an query:
	sum := func(v, e T) interface{} {
		return v.(int) + e.(int)
	}
	v := From([]T{1, 2, 3}).Reduce(sum)
	fmt.Printf("Reduced elements to sum: %v", v)

	// Output:
	// Reduced elements to sum: 6
}

func ExampleQuery_Skip_found() {
	v := From([]T{1, 2, 3, 4, 5}).Skip(2)
	fmt.Printf("Skipped 5 elements: %v", v)

	// Output:
	// Skipped 5 elements: [3 4 5]
}

func ExampleQuery_Skip_empty() {
	v := From([]T{1, 2, 3, 4, 5}).Skip(42)
	fmt.Printf("Skipped 5 elements: %v", v)

	// Output:
	// Skipped 5 elements: []
}

func ExampleQuery_Sort_decreasing() {
	decreasing := func(e, f T) bool {
		return e.(int) > f.(int)
	}
	q := From([]T{4, 3, 7, 2, 5, 9, 1, 6, 8})
	fmt.Println(q.Sort(decreasing))

	// Output:
	// [9 8 7 6 5 4 3 2 1]
}

func ExampleQuery_Sort_increasing() {
	increasing := func(e, f T) bool {
		return e.(int) < f.(int)
	}
	q := From([]T{4, 3, 7, 2, 5, 9, 1, 6, 8})
	fmt.Println(q.Sort(increasing))

	// Output:
	// [1 2 3 4 5 6 7 8 9]
}

func ExampleQuery_Take_some() {
	v := From([]T{1, 2, 3, 4, 5}).Take(3)
	fmt.Printf("Taken elements: %v", v)

	// Output:
	// Taken elements: [1 2 3]
}

func ExampleQuery_Take_all() {
	v := From([]T{1, 2, 3, 4, 5}).Take(42)
	fmt.Printf("Taken elements: %v", v)

	// Output:
	// Taken elements: [1 2 3 4 5]
}

func ExampleQuery_Where_greaterThan() {
	where := func(e T) bool {
		return e.(int) > 3
	}
	v := From([]T{1, 2, 3, 4, 5}).Where(where)
	fmt.Printf("Where: %v\n", v)

	// Output:
	// Where: [4 5]
}

func ExampleQuery_Where_lessThan() {
	where := func(e T) bool {
		return e.(int) > 3
	}
	v := From([]T{1, 2, 3}).Where(where)
	fmt.Printf("Where: %v\n", v)

	// Output:
	// Where: []
}
