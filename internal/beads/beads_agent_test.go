package beads

import (
	"fmt"
	"testing"
)

func TestIsDoltOrWispError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
		{
			name: "normal error",
			err:  fmt.Errorf("bd create: UNIQUE constraint failed"),
			want: false,
		},
		{
			name: "nil pointer dereference",
			err:  fmt.Errorf("bd create: panic: runtime error: invalid memory address nil pointer dereference"),
			want: true,
		},
		{
			name: "SIGSEGV",
			err:  fmt.Errorf("bd create: signal SIGSEGV: segmentation violation"),
			want: true,
		},
		{
			name: "panic prefix",
			err:  fmt.Errorf("bd create: panic: some dolt error"),
			want: true,
		},
		{
			name: "runtime error",
			err:  fmt.Errorf("bd create: runtime error: index out of range"),
			want: true,
		},
		{
			name: "signal in error",
			err:  fmt.Errorf("bd create: signal: killed"),
			want: true,
		},
		{
			name: "wisps table issue",
			err:  fmt.Errorf("bd create: table 'wisps' does not exist"),
			want: true,
		},
		{
			name: "DoltDB error",
			err:  fmt.Errorf("github.com/dolthub/dolt/go/libraries/doltcore/doltdb.(*DoltDB).SetCrashOnFatalError"),
			want: true,
		},
		{
			name: "not found error",
			err:  fmt.Errorf("bd create: not found"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isDoltOrWispError(tt.err)
			if got != tt.want {
				t.Errorf("isDoltOrWispError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}
