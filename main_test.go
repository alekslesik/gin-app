package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

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

// This tests that a fresh database returns no rows (but no error) when
// fetching Books.
func TestBookEmpty(t *testing.T) {
	db := freshDb(t)
	books := []Book{}
	if err := db.Find(&books).Error; err != nil {
		t.Fatalf("Error querying books from fresh db: %s", err)
	}
	if len(books) != 0 {
		t.Errorf("Expected 0 books, got %d", len(books))
	}
}

func TestBookIndex(t *testing.T) {
	t.Parallel()
	w := httptest.NewRecorder()
	ctx, r := gin.CreateTestContext(w)
	setupRouter(r, freshDb(t))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/books/", nil)
	if err != nil {
		t.Errorf("got error: %s", err)
	}

	r.ServeHTTP(w, req)
	if http.StatusOK != w.Code {
		t.Fatalf("expected response code %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	expected := "<h2>My Books</h2>"

	if !strings.Contains(body, expected) {
		t.Fatalf("expected response body to contain '%s', got '%s'", expected, body)
	}
}

func TestBookIndexError(t *testing.T) {
	t.Parallel()

	db := freshDb(t)

	if err := db.Migrator().DropTable(&Book{}); err != nil {
		t.Fatalf("got error: %s", err)
	}

	_ = getHasStatus(t, db, "/books/", http.StatusInternalServerError)
}

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
	return  w
}

func TestBookIndexNominal(t *testing.T)  {
	t.Parallel()

	db := freshDb(t)

	b := &Book{Title: "Book1", Author: "Author1"}

	if err := db.Create(&b).Error; err != nil {
		t.Fatalf("error creating book: %s", err)
	}

	b = &Book{Title: "Book2", Author: "Author2"}

	if err := db.Create(&b).Error; err != nil {
		t.Fatalf("error creating book: %s", err)
	}

	w := httptest.NewRecorder()
	ctx, router := gin.CreateTestContext(w)
	setupRouter(router, db)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/books/", nil)
	if err != nil {
		t.Errorf("got error: %s", err)
	}

	router.ServeHTTP(w, req)
	if http.StatusOK != w.Code {
		t.Errorf("expected response code %d, got %d", http.StatusOK, w.Code)
	}

	body := w.Body.String()

	fragments :=[]string{
		"<h2>My Books</h2>",
		"<li>Book1 -- Author1</li>",
		"<li>Book2 -- Author2</li>",
	}

	for _, fragment := range fragments {
		if !strings.Contains(body, fragment) {
			t.Fatalf("expected body to contain '%s', got %s", fragment, body)
		}
	}
}
