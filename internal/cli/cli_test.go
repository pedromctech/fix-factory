package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pedromctech/fix-factory/internal/cli"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		in        string
		wantCode  int
		wantOut   []string
		wantErr   []string
		wantNoOut bool
		wantNoErr bool
	}{
		{name: "add ok", args: []string{"add", "octo/hello-world"}, wantCode: 0, wantOut: []string{"added octo/hello-world"}, wantNoErr: true},
		{name: "add url", args: []string{"add", "https://github.com/octo/hello-world"}, wantCode: 0, wantOut: []string{"added octo/hello-world"}, wantNoErr: true},
		{name: "add path traversal", args: []string{"add", "../etc"}, wantCode: 1, wantErr: []string{"error:"}, wantNoOut: true},
		{name: "add not a repo", args: []string{"add", "not a repo"}, wantCode: 1, wantErr: []string{"error:"}, wantNoOut: true},
		{name: "interactive", in: "octo/hello\nquit\n", wantCode: 0, wantOut: []string{"added octo/hello"}},
		{name: "interactive eof", in: "", wantCode: 0},
		{name: "interactive bad then good", in: "../x\nocto/hello\n", wantCode: 0, wantOut: []string{"error:", "added octo/hello"}},
		{name: "status", args: []string{"status"}, wantCode: 0, wantOut: []string{"[ ] scan", "[ ] create-pr"}, wantNoErr: true},
		{name: "status extra args", args: []string{"status", "x"}, wantCode: 2, wantErr: []string{"usage:"}, wantNoOut: true},
		{name: "help", args: []string{"help"}, wantCode: 0, wantOut: []string{"usage:"}, wantNoErr: true},
		{name: "-h", args: []string{"-h"}, wantCode: 0, wantOut: []string{"usage:"}, wantNoErr: true},
		{name: "--help", args: []string{"--help"}, wantCode: 0, wantOut: []string{"usage:"}, wantNoErr: true},
		{name: "unknown", args: []string{"bogus"}, wantCode: 2, wantErr: []string{"usage:"}, wantNoOut: true},
		{name: "add no arg", args: []string{"add"}, wantCode: 2, wantErr: []string{"usage:"}, wantNoOut: true},
		{name: "add extra args", args: []string{"add", "a/b", "c/d"}, wantCode: 2, wantErr: []string{"usage:"}, wantNoOut: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out, errOut bytes.Buffer
			code := cli.Run(tt.args, strings.NewReader(tt.in), &out, &errOut)
			if code != tt.wantCode {
				t.Fatalf("code = %d, want %d (out=%q err=%q)", code, tt.wantCode, out.String(), errOut.String())
			}
			for _, s := range tt.wantOut {
				if !strings.Contains(out.String(), s) {
					t.Errorf("stdout %q missing %q", out.String(), s)
				}
			}
			for _, s := range tt.wantErr {
				if !strings.Contains(errOut.String(), s) {
					t.Errorf("stderr %q missing %q", errOut.String(), s)
				}
			}
			if tt.wantNoOut && out.Len() != 0 {
				t.Errorf("unexpected stdout %q", out.String())
			}
			if tt.wantNoErr && errOut.Len() != 0 {
				t.Errorf("unexpected stderr %q", errOut.String())
			}
		})
	}
}
