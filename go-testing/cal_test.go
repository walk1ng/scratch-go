package main

import (
	"fmt"
	"os"
	"testing"
)

func setup() {
	fmt.Println("before all tests")
}

func teardown() {
	fmt.Println("after all tests")
}

func TestAdd(t *testing.T) {
	if ans := Add(1, 2); ans != 3 {
		t.Errorf("expected 1+2=3, but got %d", ans)
	}

	if ans := Add(-1, -2); ans != -3 {
		t.Errorf("expected -1+-2=3, but got %d", ans)
	}
	if ans := Add(-1, 0); ans != -1 {
		t.Errorf("expected -1+0=-1, but got %d", ans)
	}
}

type testcase struct {
	name       string
	a, b, want int
}

func createtestcases(t *testing.T, c *testcase) {
	t.Run(c.name, func(t *testing.T) {
		t.Helper()
		if got := Mul(c.a, c.b); got != c.want {
			t.Errorf("Mul(%d,%d) = %v, want %v", c.a, c.b, got, c.want)
		}
	})
}

func TestMul(t *testing.T) {
	createtestcases(t, &testcase{"pos", 2, 3, 6})
	createtestcases(t, &testcase{"neg", 2, -3, 6})
	createtestcases(t, &testcase{"zero", 0, 3, 0})
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}
