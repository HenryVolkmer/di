package di

import "testing"

func assertEquals(expected,actual string,t *testing.T) {
    if expected != actual {
        t.Fatalf("Failed asserting that two strings are equal.\n--- Expected\n+++ Actual\n@@ @@\n-'%s'\n+'%s'", expected, actual)
    }
}