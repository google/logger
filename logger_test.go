package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewLogger(t *testing.T) {
	var buf1 bytes.Buffer
	l1, err := NewLogger("test1", &buf1)
	if err != nil {
		t.Fatal("Could not create logger")
	}
	// Subsequent runs of Init shouldn't change defaultLogger.
	var buf2 bytes.Buffer
	l2, err := NewLogger("test2", &buf2)
	if err != nil {
		t.Fatal("Could not create logger")
	}

	// Check log output.
	l1.Info("logger #1")
	l2.Info("logger #2")

	tests := []struct {
		out  string
		want int
	}{
		{buf1.String(), 1},
		{buf2.String(), 1},
	}

	for i, tt := range tests {
		got := len(strings.Split(strings.TrimSpace(tt.out), "\n"))
		if got != tt.want {
			t.Errorf("logger %d wrong number of lines, want %d, got %d", i+1, tt.want, got)
		}
	}
}
