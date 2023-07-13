package main

import "testing"

func test1(t *testing.T) {
	// ruleid: error-failnow-fatal
	t.Error("abc")
	t.FailNow()
}

func test2(t *testing.T) {
	// ruleid: error-failnow-fatal
	t.Errorf("%s", "abc")
	t.FailNow()
}

func test3(t *testing.T) {
	// ruleid: error-failnow-fatal
	t.Log("abc")
	t.FailNow()
}

func test4(t *testing.T) {
	// ruleid: error-failnow-fatal
	t.Logf("%s", "abc")
	t.FailNow()
}
