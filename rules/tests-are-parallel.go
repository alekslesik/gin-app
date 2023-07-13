package main

import (
	"log"
	"testing"
)

// ruleid: tests-are-parallel
func test1(t *testing.T) {
	// Note: parallel must be called.
}

// ruleid: tests-are-parallel
func test2(t *testing.T) {
	// Note: parallel has to be called first.
	log.Print("abc")
	t.Parallel()
}

// ok: tests-are-parallel
func test3(t *testing.T) {
	// Note: parallel is called first. This is ok.
	t.Parallel()
	log.Print("abc")
}

// ruleid: tests-are-parallel
func testHelper1(t *testing.T) {
	// Note: helper has to be called first.
	log.Print("abc")
	t.Helper()
}

// ok: tests-are-parallel
func testHelper2(t *testing.T) {
	// Note: helper is called first. This is ok.
	t.Helper()
	log.Print("abc")
}