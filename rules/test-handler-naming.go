package main

import "testing"

// Doesn't match because it doesn't use getHasStatus/postHasStatus.
// ok: test-handler-naming
func testFoo(t *testing.T) {}

// Has all the parts.
// ok: test-handler-naming
func testThingIndexGet(t *testing.T) {
	getHasStatus(t)
}

// Has all the parts, plus extra.
// ok: test-handler-naming
func testThingIndexGetStuff(t *testing.T) {
	getHasStatus(t)
}

// Has Thing, but missing crud+method.
// ruleid: test-handler-naming
func testFoo(t *testing.T) {
	getHasStatus(t)
}

// Missing "Thing"
// ruleid: test-handler-naming
func testIndexGet(t *testing.T) {
	getHasStatus(t)
}

// Missing method
// ruleid: test-handler-naming
func testThingIndex(t *testing.T) {
	getHasStatus(t)
}

// Wrong order for crud+method
// ruleid: test-handler-naming
func testThingGetIndex(t *testing.T) {
	getHasStatus(t)
}
