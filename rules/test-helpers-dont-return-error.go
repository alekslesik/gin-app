package main

import (
	"testing"
)

// ruleid: test-helpers-dont-return-error
func testHelper1(t *testing.T) error {
	t.Helper()
	return nil
}

// ok: test-helpers-dont-return-error
func testHelper2(t *testing.T) int {
	t.Helper()
	return 0
}

// ok: test-helpers-dont-return-error
func testHelper2(t *testing.T) {
	t.Helper()
}