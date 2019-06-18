// Copyright 2019 Daniel Mundt. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
//
// SPDX-License-Identifier: MIT
//

package query

import (
	"fmt"
)

type Author struct {
	AuthorID int
	Name     string
}

type AuthorBook struct {
	AuthorID int
	BookID   int
}

type Book struct {
	BookID int
	Title  string
	Year   int
}

type NameBookID struct {
	Name   string
	BookID int
}

type AuthorTitleYear struct {
	Author string
	Title  string
	Year   int
}

func (a AuthorTitleYear) String() string {
	return fmt.Sprintf("{%v: %v (%v)}", a.Author, a.Title, a.Year)
}

func Example() {
	// Authors table:
	authors := []T{
		Author{1, "Austen, Jane"},
		Author{2, "Brontë, Emily"},
		Author{3, "Hunter, Rachel"},
	}

	// Books table:
	books := []T{
		Book{1, "Sense & Sensibility", 1811},
		Book{2, "Pride & Prejudice", 1813},
		Book{3, "Mansfield Park", 1814},
		Book{4, "Emma", 1815},
		Book{5, "Persuasion", 1817},
		Book{6, "Northanger Abbey", 1817},
		Book{7, "Sanditon", 1817},
		Book{8, "Wuthering Heights", 1847},
		Book{9, "Letitia, or, The Castle without a Spectre", 1801},
		Book{10, "The History of the Grubthorpe Family", 1802},
		Book{11, "Letters from Mrs Palmerstone to her Daughter, Inculcating Morality by Entertaining Narratives", 1803},
		Book{12, "The Unexpected Legacy", 1804},
		Book{13, "Family Annals", 1807},
		Book{14, "The Schoolmistress", 1811},
	}

	// Authors to books table:
	author2Books := []T{
		AuthorBook{1, 1},
		AuthorBook{1, 2},
		AuthorBook{1, 3},
		AuthorBook{1, 4},
		AuthorBook{1, 5},
		AuthorBook{1, 6},
		AuthorBook{1, 7},
		AuthorBook{2, 8},
		AuthorBook{3, 9},
		AuthorBook{3, 10},
		AuthorBook{3, 11},
		AuthorBook{3, 12},
		AuthorBook{3, 13},
		AuthorBook{3, 14},
	}

	// Print all authors, title of their books, published between 1804 and 1815:
	query := From([]T(authors)).
		Join(From(author2Books),
			func(e T) interface{} {
				return e.(Author).AuthorID
			}, func(e T) interface{} {
				return e.(AuthorBook).AuthorID
			}, func(e1, e2 interface{}) interface{} {
				return NameBookID{e1.(Author).Name, e2.(AuthorBook).BookID}
			}).
		Join(From(books),
			func(e T) interface{} {
				return e.(NameBookID).BookID
			}, func(e T) interface{} {
				return e.(Book).BookID
			}, func(e1, e2 interface{}) interface{} {
				return AuthorTitleYear{e1.(NameBookID).Name, e2.(Book).Title, e2.(Book).Year}
			}).
		Sort(
			func(e1, e2 T) bool {
				return e1.(AuthorTitleYear).Year > e2.(AuthorTitleYear).Year
			},
			func(e1, e2 T) bool {
				return e1.(AuthorTitleYear).Author < e2.(AuthorTitleYear).Author
			}).
		Where(
			func(e T) bool {
				return e.(AuthorTitleYear).Year >= 1804
			}, func(e T) bool {
				return e.(AuthorTitleYear).Year <= 1815
			})
	fmt.Printf("%v\n", query)

	// Output:
	// [{Austen, Jane: Emma (1815)} {Austen, Jane: Mansfield Park (1814)} {Austen, Jane: Pride & Prejudice (1813)} {Austen, Jane: Sense & Sensibility (1811)} {Hunter, Rachel: The Schoolmistress (1811)} {Hunter, Rachel: Family Annals (1807)} {Hunter, Rachel: The Unexpected Legacy (1804)}]
}
