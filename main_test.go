package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestBookIndexRefactored(t *testing.T) {
	t.Parallel()

	db := freshDb(t)
	books := createBooks(t, db, 2)

	w := getHasStatus(t, db, "/books/", http.StatusOK)
	body := w.Body.String()
	fragments := []string{
		`<h1>My Books</h1>`,
		`<ul class="books">`,
		fmt.Sprintf(`<span class="title">%s</span>`, books[0].Title),
		// fmt.Sprintf(`<span class="author">%s</span>`, books[0].Author),
	}

	bodyHasFragments(t, body, fragments)
}

// Table Test sample
func TestBookIndexTable(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc  string
		count int
	}{
		{"empty", 0},
		{"single", 1},
		{"multiple", 10},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			tC := tC
			t.Parallel()

			db := freshDb(t)
			books := createBooks(t, db, tC.count)

			w := getHasStatus(t, db, "/books/", http.StatusOK)
			body := w.Body.String()

			fragments := []string{
				"<h1>My Books</h1>",
			}

			for _, book := range books {
				fragments = append(fragments, fmt.Sprintf(`<span class="title">%s</span>`, book.Title))
			}
			bodyHasFragments(t, body, fragments)
		})
	}
}

func TestBookNewGet(t *testing.T) {
	t.Parallel()
	tcs := []struct {
		name string
	}{
		{"basic"},
	}

	for i := range tcs {
		tc := &tcs[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := freshDb(t)
			w := getHasStatus(t, db, "/books/new", http.StatusOK)
			body := w.Body.String()
			fragments := []string{
				"<h2>Add a Book</h2>",
				`<form action="/books/new" method="POST">`,
				`<input type="text" name="title" id="title"`,
				`<input type="text" name="author" id="author"`,
				`<button type="submit"`,
			}
			bodyHasFragments(t, body, fragments)
		})
	}
}

func TestBookNewPost(t *testing.T) {
	t.Parallel()

	dropTable := func(t *testing.T, db *gorm.DB) {
		err := db.Migrator().DropTable("books")
		if err != nil {
			t.Fatalf("error dropping table 'books': %s", err)
		}
	}

	tcs := []struct {
		name      string
		data      gin.H
		setup     func(*testing.T, *gorm.DB)
		status    int
		fragments []string
	}{
		{
			name:   "nominal",
			data:   gin.H{"title": "my book", "author": "me"},
			status: http.StatusFound,
		},
		{
			// This makes the field validation fail because the
			// author is empty.
			name:   "empty_author",
			data:   gin.H{"title": "1"},
			status: http.StatusBadRequest,
			fragments: []string{
				"Author is required, but was empty",
			},
		},
		{
			// This makes the field validation fail because the
			// title is empty.
			name:   "empty_title",
			data:   gin.H{"author": "9"},
			status: http.StatusBadRequest,
			fragments: []string{
				"Title is required, but was empty",
			},
		},
		{
			// This makes the field validation fail because both
			// title and author are empty.
			name:   "empty",
			data:   gin.H{},
			status: http.StatusBadRequest,
			fragments: []string{
				"Author is required, but was empty",
				"Title is required, but was empty",
			},
		},
		{
			name:   "db_error",
			data:   gin.H{"title": "a", "author": "b"},
			setup:  dropTable,
			status: http.StatusInternalServerError,
		},
	}

	for i := range tcs {
		tc := &tcs[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := freshDb(t)
			if tc.setup != nil {
				tc.setup(t, db)
			}
			w := postHasStatus(t, db, "/books/new", &tc.data,
				tc.status)

			if tc.fragments != nil {
				body := w.Body.String()
				bodyHasFragments(t, body, tc.fragments)
			}

			if tc.status == http.StatusFound {
				// Make sure the record is in the db.
				books := []Book{}
				result := db.Find(&books)
				if result.Error != nil {
					t.Fatalf("error fetching books: %s", result.Error)
				}
				if result.RowsAffected != 1 {
					t.Fatalf("expected 1 row affected, got %d",
						result.RowsAffected)
				}
				if tc.data["title"] != books[0].Title {
					t.Fatalf("expected title '%s', got '%s",
						tc.data["title"], books[0].Title)
				}
				if tc.data["author"] != books[0].Author {
					t.Fatalf("expected author '%s', got '%s",
						tc.data["author"], books[0].Author)
				}

				// Check the redirect location.
				url, err := w.Result().Location()
				if err != nil {
					t.Fatalf("location check error: %s", err)
				}

				if url.String() != "/books/" {
					t.Errorf("expected location '/books/', got '%s'",
						url.String())
				}

				w := getCookieHasStatus(t, db, url.String(), w.Result(), http.StatusOK)
				fragments := []string{
					fmt.Sprintf("New book &#39;%s&#39; saved successfully.", books[0].Title),
				}

				bodyHasFragments(t, w.Body.String(), fragments)

			}
		})
	}
}

// Helpers
//* This just moves the loop at the bottom of the test function into its own function. There are a few things to note here:
//* 1.This is marked as a t.Helper() – failures will be reported at the line number of the test case function instead of inside this function, which can make diagnosing a failing test easier.
//* 2.Test helper functions should not return errors! Don’t create extra boilerplate in your tests with error handling. Helper functions should check errors and fail instead of passing errors up the stack.
//* 3.Both 1 and 2 are why every test helper function should take a *testing.T as the first argument.
//* 4.I try to name my verification-type helper functions in the form thingHasProperty or thingIsStatus. I think this makes it easier to read what test code is checking.

// Check that template includes fragment
func bodyHasFragments(t *testing.T, body string, fragments []string) bool {
	t.Helper()
	for _, fragment := range fragments {
		if !strings.Contains(body, fragment) {
			t.Fatalf("expected body to contain '%s', got %s", fragment, body)
		}
	}

	return true
}

// Checking status code
func getHasStatus(t *testing.T, db *gorm.DB, path string, status int) *httptest.ResponseRecorder {
	t.Helper()

	w := httptest.NewRecorder()
	ctx, router := gin.CreateTestContext(w)
	setupRouter(router, db)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		t.Errorf("got error: %s", err)
	}

	router.ServeHTTP(w, req)
	if status != w.Code {
		t.Errorf("expected response code %d, got %d", status, w.Code)
	}
	return w
}

// Create slice of books
func createBooks(t *testing.T, db *gorm.DB, count int) []*Book {
	books := []*Book{}
	t.Helper()
	for i := 0; i < count; i++ {
		b := &Book{
			Title:  fmt.Sprintf("Book%03d", i),
			Author: fmt.Sprintf("Author%03d", i),
		}
		if err := db.Create(b).Error; err != nil {
			t.Fatalf("error creating book: %s", err)
		}

		books = append(books, b)
	}

	return books
}

// Create new db
func freshDb(t testing.TB, path ...string) *gorm.DB {
	t.Helper()

	var dbUri string

	// Note: path can be specified in an individual test for debugging
	// purposes -- so the db file can be inspected after the test runs.
	// Normally it should be left off so that a truly fresh memory db is
	// used every time.
	if len(path) == 0 {
		dbUri = ":memory:"
	} else {
		dbUri = path[0]
	}

	db, err := gorm.Open(sqlite.Open(dbUri), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error opening memory db: %s", err)
	}
	if err := setupDatabase(db); err != nil {
		t.Fatalf("Error setting up db: %s", err)
	}
	return db
}

func postHasStatus(t *testing.T, db *gorm.DB, path string, h *gin.H, status int) *httptest.ResponseRecorder {
	t.Helper()

	data := url.Values{}

	for k, vi := range *h {
		v := vi.(string)

		data.Set(k, v)
	}

	w := httptest.NewRecorder()
	ctx, router := gin.CreateTestContext(w)
	setupRouter(router, db)

	req, err := http.NewRequestWithContext(ctx, "POST", path, strings.NewReader(data.Encode()))
	if err != nil {
		t.Errorf("got error: %s", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)
	responseHasCode(t, w, status)
	return w
}

func responseHasCode(t *testing.T, w *httptest.ResponseRecorder, expected int) {
	if expected != w.Code {
		t.Errorf("expected response code %d, got %d", expected, w.Code)
	}
}

func FuzzPaginate(f *testing.F) {
	f.Add("1", 100, 10)
	f.Add("5", 0, 50)
	f.Add("10", 250, 50)

	f.Fuzz(func(t *testing.T, pageStr string, n, per int) {
		p, err := paginate(pageStr, n, per)
		if err != nil {
			// TODO: verify pageStr is invalid int
			return
		}
		if p == nil {
			t.Fatal("p is nil")
		}
		if p.Page < 1 || p.Page > p.Count {
			t.Fatalf("p.Page is %d (count %d)", p.Page, p.Count)
		}
		if p.Count <= 0 {
			t.Fatalf("p.Count is %d", p.Count)
		}
		if p.Offset < 0 || p.Offset > n {
			t.Fatalf("p.Offset is %d (n=%d, per=%d)", p.Offset, n, per)
		}
		if p.Page == 1 {
			if p.Prev != 0 {
				t.Fatalf("p.Page is %d but p.Prev is not zero (%d)", p.Page, p.Prev)
			}
		} else if p.Prev+1 != p.Page {
			t.Fatalf("prev %d+1 != %d", p.Prev, p.Page)
		}
		if p.Page == p.Count {
			if p.Next != 0 {
				t.Fatalf("p.Page is %d but p.Next is not zero (%d)", p.Page, p.Next)
			}
		} else if p.Next-1 != p.Page {
			t.Fatalf("next %d-1 != %d", p.Next, p.Page)
		}
	})
}

func getCookieHasStatus(t *testing.T, db *gorm.DB, path string, r *http.Response, status int) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	ctx, router := gin.CreateTestContext(w)
	os.Setenv("AKLATAN_SESSION_KEY", "dummy")
	setupRouter(router, db)
	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		t.Errorf("got error: %s", err)
	}

	if r != nil {
		req.Header["Cookie"] = r.Header["Set-Cookie"]
	}

	router.ServeHTTP(w, req)
	if status != w.Code {
		t.Errorf("expected response code %d, got %d", status, w.Code)
	}
	return w
}
