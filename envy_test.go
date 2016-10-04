package envy

import (
	"os"
	"testing"
)

func TestBasic(t *testing.T) {
	os.Setenv("FOO", "BAR")
	os.Setenv("FII", "BIR")

	e := New()
	e.Add("FOO")
	e.Add("FII")

	res := e.Getenv()
	expected := "BAR"
	if res != expected {
		t.Errorf("Expected %s but got %s", expected, res)
	}
}

func TestFusion(t *testing.T) {
	os.Setenv("REDIS_PORT_6379_TCP_ADDR", "127.0.0.1")
	os.Setenv("REDIS_PORT_6379_TCP_PORT", "6379")

	res := New().Merge(Join(":"), "REDIS_PORT_6379_TCP_ADDR", "REDIS_PORT_6379_TCP_PORT").
		Getenv()

	expected := "127.0.0.1:6379"
	if res != expected {
		t.Errorf("Expected %s but got %s", expected, res)
	}
}
