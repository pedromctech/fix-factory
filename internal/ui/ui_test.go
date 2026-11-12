package ui_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/pedromctech/fix-factory/internal/repo"
	"github.com/pedromctech/fix-factory/internal/ui"
)

type fakeAdder struct {
	added []repo.Ref
	err   error
}

func (f *fakeAdder) Add(r repo.Ref) error {
	if f.err != nil {
		return f.err
	}
	f.added = append(f.added, r)
	return nil
}

func TestAddsValidRepos(t *testing.T) {
	tests := []struct {
		in   string
		want repo.Ref
		out  string
	}{
		{"octo/hello-world\n", repo.Ref{Owner: "octo", Name: "hello-world"}, "added octo/hello-world"},
		{"https://github.com/octo/hello.git\n", repo.Ref{Owner: "octo", Name: "hello"}, "added octo/hello"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			f := &fakeAdder{}
			var out bytes.Buffer
			if err := ui.Run(strings.NewReader(tt.in), &out, f); err != nil {
				t.Fatal(err)
			}
			if len(f.added) != 1 || f.added[0] != tt.want {
				t.Fatalf("added = %v, want [%v]", f.added, tt.want)
			}
			if !strings.Contains(out.String(), tt.out) {
				t.Errorf("output %q lacks %q", out.String(), tt.out)
			}
		})
	}
}

func TestRejectsBadInputAndContinues(t *testing.T) {
	for _, bad := range []string{"not-a-repo", "a/b/c", "../x"} {
		t.Run(bad, func(t *testing.T) {
			f := &fakeAdder{}
			var out bytes.Buffer
			if err := ui.Run(strings.NewReader(bad+"\nocto/ok\n"), &out, f); err != nil {
				t.Fatal(err)
			}
			if len(f.added) != 1 || f.added[0].String() != "octo/ok" {
				t.Fatalf("added = %v, want only octo/ok", f.added)
			}
			if !strings.Contains(out.String(), "error") {
				t.Errorf("no error printed: %q", out.String())
			}
		})
	}
}

func TestIgnoresBlankLines(t *testing.T) {
	f := &fakeAdder{}
	var out bytes.Buffer
	if err := ui.Run(strings.NewReader("\n   \n\t\n"), &out, f); err != nil {
		t.Fatal(err)
	}
	if len(f.added) != 0 {
		t.Errorf("added = %v", f.added)
	}
	if strings.Contains(out.String(), "error") {
		t.Errorf("unexpected error output: %q", out.String())
	}
}

func TestQuitStopsLoop(t *testing.T) {
	for _, q := range []string{"quit", "exit"} {
		t.Run(q, func(t *testing.T) {
			f := &fakeAdder{}
			var out bytes.Buffer
			if err := ui.Run(strings.NewReader("octo/a\n"+q+"\nocto/b\n"), &out, f); err != nil {
				t.Fatal(err)
			}
			if len(f.added) != 1 {
				t.Errorf("added = %v, want one", f.added)
			}
		})
	}
	var out bytes.Buffer
	if err := ui.Run(strings.NewReader("octo/a"), &out, &fakeAdder{}); err != nil {
		t.Errorf("EOF should return nil, got %v", err)
	}
}

func TestAdderErrorIsReturned(t *testing.T) {
	sentinel := errors.New("boom")
	var out bytes.Buffer
	err := ui.Run(strings.NewReader("octo/a\n"), &out, &fakeAdder{err: sentinel})
	if !errors.Is(err, sentinel) {
		t.Fatalf("err = %v, want sentinel", err)
	}
}

func TestMemoryIgnoresDuplicates(t *testing.T) {
	m := ui.NewMemory()
	a := repo.Ref{Owner: "octo", Name: "a"}
	b := repo.Ref{Owner: "octo", Name: "b"}
	for _, r := range []repo.Ref{a, b, a} {
		if err := m.Add(r); err != nil {
			t.Fatal(err)
		}
	}
	got := m.List()
	if len(got) != 2 || got[0] != a || got[1] != b {
		t.Fatalf("List = %v", got)
	}
}
