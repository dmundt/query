# query

[![GoDoc](https://godoc.org/github.com/dmundt/query?status.svg)](https://godoc.org/github.com/dmundt/query) [![Build Status](https://dev.azure.com/dmundt/query/_apis/build/status/dmundt.query?branchName=master)](https://dev.azure.com/dmundt/query/_build/latest?definitionId=3&branchName=master) [![Build Status](https://github.com/dmundt/query/workflows/Go/badge.svg)](https://github.com/dmundt/query/actions) [![Build Status](https://travis-ci.com/dmundt/query.svg?branch=master)](https://travis-ci.com/dmundt/query) [![Coverage Status](https://coveralls.io/repos/github/dmundt/query/badge.svg?branch=master)](https://coveralls.io/github/dmundt/query?branch=master) [![Codebeat Badge](https://codebeat.co/badges/369691b5-4735-405b-a83d-a61835e346d0)](https://codebeat.co/projects/github-com-dmundt-query-master) [![Go Report Card](https://goreportcard.com/badge/github.com/dmundt/query)](https://goreportcard.com/report/github.com/dmundt/query) [![Donate](https://img.shields.io/badge/Donate-Paypal-blue.svg)](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=Z5BSACFWB8DRE&source=url)

Simple query language written in Go inspired by Dart's [Iterable&lt;E>](https://api.dartlang.org/stable/2.2.0/dart-core/Iterable-class.html) with cascaded method invocation:

- [Any()](https://godoc.org/github.com/dmundt/query#Query.Any)
- [At()](https://godoc.org/github.com/dmundt/query#Query.At)
- [Contains()](https://godoc.org/github.com/dmundt/query#Query.Contains)
- [Every()](https://godoc.org/github.com/dmundt/query#Query.Every)
- [Expand()](https://godoc.org/github.com/dmundt/query#Query.Expand)
- [First()](https://godoc.org/github.com/dmundt/query#Query.First)
- [Fold()](https://godoc.org/github.com/dmundt/query#Query.Fold)
- [ForEach()](https://godoc.org/github.com/dmundt/query#Query.ForEach)
- [From()](https://godoc.org/github.com/dmundt/query#From)
- [IsEmpty()](https://godoc.org/github.com/dmundt/query#Query.IsEmpty)
- [Join()](https://godoc.org/github.com/dmundt/query#Query.Join)
- [Last()](https://godoc.org/github.com/dmundt/query#Query.Last)
- [MapTo()](https://godoc.org/github.com/dmundt/query#Query.MapTo)
- [Reduce()](https://godoc.org/github.com/dmundt/query#Query.Reduce)
- [Skip()](https://godoc.org/github.com/dmundt/query#Query.Skip)
- [Sort()](https://godoc.org/github.com/dmundt/query#Query.Sort)
- [String()](https://godoc.org/github.com/dmundt/query#Query.String)
- [Take()](https://godoc.org/github.com/dmundt/query#Query.Task)
- [Where()](https://godoc.org/github.com/dmundt/query#Query.Where)

## Installation

```golang
go get github.com/dmundt/query
```

## Quickstart

See [examples](https://godoc.org/github.com/dmundt/query#pkg-examples) at godoc.org.

## Example

The following example implements a book database.

Create the tables:

```golang
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
```

Implement String():

```golang
func (a AuthorTitleYear) String() string {
    return fmt.Sprintf("{%v: %v (%v)}", a.Author, a.Title, a.Year)
}
```

Populate all tables:

```golang
func Example() {
    // Populate the authors table:
    authors := []T{
        Author{1, "Austen, Jane"},
        Author{2, "Brontë, Emily"},
        Author{3, "Hunter, Rachel"},
    }

    // Populate the books table:
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

    // Map authors to books.
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
```

Print all authors, title of their books, published between 1804 and 1815:

```golang
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
}
```

Executing the query yields the following output:

```golang
    // Output:
    // [{Austen, Jane: Emma (1815)} {Austen, Jane: Mansfield Park (1814)} {Austen, Jane: Pride & Prejudice (1813)} {Austen, Jane: Sense & Sensibility (1811)} {Hunter, Rachel: The Schoolmistress (1811)} {Hunter, Rachel: Family Annals (1807)} {Hunter, Rachel: The Unexpected Legacy (1804)}]
}
```

## Stability and Compatibility

The query package is considered stable. We will make every effort to ensure API compatibility in future releases.

## Semantic versioning

Package query uses [semantic versioning](https://semver.org/ "semantic versioning") for satisfying dependency requirements of [Go Modules](https://blog.golang.org/using-go-modules/ "golang modules").

## License

Package query is covered by a [MIT](https://github.com/dmundt/query/blob/master/LICENSE) license:

```text
MIT License

Copyright (c) 2019 Daniel Mundt

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
