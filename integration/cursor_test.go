package integration

import (
	"testing"
)

func TestComputeCursorCol(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want int
	}{
		{
			name: "Simple prompt",
			data: []byte("λ "),
			want: 2,
		},
		{
			name: "ANSI color prompt",
			data: []byte("\033[32mλ \033[0m"),
			want: 2,
		},
		{
			name: "Carriage return and move right",
			data: []byte("λ \033[140G...22 chars right prompt...\r\033[2C"),
			want: 2,
		},
		{
			name: "OSC sequence before prompt",
			data: []byte("\033]0;iris on fix/menu-debouncing\007λ "),
			want: 2,
		},
		{
			name: "CSI Horizontal Absolute",
			data: []byte("abc\033[10Gde"),
			want: 11,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ComputeCursorCol(tt.data)
			if got != tt.want {
				t.Errorf("ComputeCursorCol(%q) = %d, want %d", tt.data, got, tt.want)
			}
		})
	}
}
