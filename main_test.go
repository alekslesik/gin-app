package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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
		"<h2>My Books</h2>",
		fmt.Sprintf("<li>%s -- %s</li>", books[0].Title, books[0].Author),
		fmt.Sprintf("<li>%s -- %s</li>", books[1].Title, books[1].Author),
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
			t.Parallel()

			db := freshDb(t)
			books := createBooks(t, db, tC.count)

			w := getHasStatus(t, db, "/books/", http.StatusOK)
			body := w.Body.String()

			fragments := []string{
				"<h2>My Books</h2>",
			}

			for _, book := range books {
				fragments = append(fragments, fmt.Sprintf("<li>%s -- %s</li>", book.Title, book.Author))
			}
			bodyHasFragments(t, body, fragments)
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

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/books/", nil)
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
func freshDb(t *testing.T, path ...string) *gorm.DB {
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
