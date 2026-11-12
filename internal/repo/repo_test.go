package repo_test

import (
	"errors"
	"testing"

	"github.com/pedromctech/fix-factory/internal/repo"
)

func TestParsesOwnerSlashName(t *testing.T) {
	r, err := repo.Parse("octo/hello-world")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Owner != "octo" || r.Name != "hello-world" {
		t.Fatalf("got %+v", r)
	}
	if r.String() != "octo/hello-world" {
		t.Fatalf("String() = %q", r.String())
	}
}

func TestParsesGitHubURL(t *testing.T) {
	r, err := repo.Parse("https://github.com/octo/hello.rs.git")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.String() != "octo/hello.rs" {
		t.Fatalf("String() = %q", r.String())
	}
}

func TestRejectsBadInput(t *testing.T) {
	cases := []struct {
		input string
		want  error
	}{
		{"octo", repo.ErrBadFormat},
		{"a/b/c", repo.ErrBadFormat},
		{"octo/../x", repo.ErrBadFormat},
		{"", repo.ErrBadFormat},
		{"octo/re po", repo.ErrBadChar},
		{"octo/..", repo.ErrBadChar},
	}
	for _, c := range cases {
		_, err := repo.Parse(c.input)
		if !errors.Is(err, c.want) {
			t.Errorf("Parse(%q) error = %v, want %v", c.input, err, c.want)
		}
	}
}
