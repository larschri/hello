package main

import (
	"fmt"
	"testing"
)

func Hello(s string) []string {
	return []string{
		"x " + s + " y",
		fmt.Sprintf("x %s y", s),
	}
}

func BenchmarkHello(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Hello("hello")
	}
}
