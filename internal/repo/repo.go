// Package repo parses and validates references to GitHub repositories.
package repo

import (
	"errors"
	"fmt"
	"strings"
)

// Ref is a GitHub repository the user added, identified by owner/name.
type Ref struct {
	Owner string
	Name  string
}

// String returns the reference as owner/name.
func (r Ref) String() string {
	return r.Owner + "/" + r.Name
}

var (
	// ErrBadFormat means the input is not owner/name or a https://github.com/owner/name URL.
	ErrBadFormat = errors.New("expected owner/name or a GitHub URL")
	// ErrBadChar means owner or name has a character GitHub does not allow.
	ErrBadChar = errors.New("invalid character")
)

// ParseError reports the input that failed to parse. Use errors.Is with
// ErrBadFormat or ErrBadChar to tell the reasons apart.
type ParseError struct {
	Input string
	Err   error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("%v: %q", e.Err, e.Input)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}

// Parse turns "owner/name" or "https://github.com/owner/name[.git]" into a Ref.
func Parse(input string) (Ref, error) {
	path := strings.TrimSpace(input)
	path = strings.TrimPrefix(path, "https://github.com/")
	path = strings.TrimRight(path, "/")
	path = strings.TrimSuffix(path, ".git")

	parts := strings.Split(path, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return Ref{}, &ParseError{Input: input, Err: ErrBadFormat}
	}
	if !validPart(parts[0]) || !validPart(parts[1]) {
		return Ref{}, &ParseError{Input: input, Err: ErrBadChar}
	}
	return Ref{Owner: parts[0], Name: parts[1]}, nil
}

func validPart(part string) bool {
	if part == "" || part == "." || part == ".." {
		return false
	}
	for _, c := range part {
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9':
		case c == '-', c == '_', c == '.':
		default:
			return false
		}
	}
	return true
}
