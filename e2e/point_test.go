package e2e

import (
	"strings"
	"testing"
)

func TestNewPoint_Defaults(t *testing.T) {
	// Y is required, so we must supply it.
	p, err := NewPoint(WithY(1))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.X != 0 || p.Y != 1 || p.Name != "origin" {
		t.Fatalf("defaults wrong: %+v", p)
	}
}

func TestNewPoint_MissingRequired(t *testing.T) {
	_, err := NewPoint()
	if err == nil {
		t.Fatal("expected error for missing required Y")
	}
	if !strings.Contains(err.Error(), "Y is required") {
		t.Fatalf("wrong error: %v", err)
	}
}

func TestNewPoint_Options(t *testing.T) {
	p, err := NewPoint(
		WithX(3),
		WithY(4),
		WithName("hello"),
		WithEmail("a@b.com"),
		WithKind("a"),
		WithTags([]string{"x", "y"}),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.X != 3 || p.Y != 4 || p.Name != "hello" {
		t.Fatalf("options not applied: %+v", p)
	}
	if len(p.Tags) != 2 || p.Tags[0] != "x" {
		t.Fatalf("tags wrong: %+v", p.Tags)
	}
}

func TestNewPoint_BadEmail(t *testing.T) {
	_, err := NewPoint(WithY(1), WithEmail("not-an-email"))
	if err == nil {
		t.Fatal("expected error for bad email")
	}
}

func TestNewPoint_BadEnum(t *testing.T) {
	_, err := NewPoint(WithY(1), WithKind("z"))
	if err == nil {
		t.Fatal("expected error for bad enum")
	}
}

func TestNewPoint_NameTooLong(t *testing.T) {
	long := strings.Repeat("x", 65)
	_, err := NewPoint(WithY(1), WithName(long))
	if err == nil {
		t.Fatal("expected error for name too long")
	}
}

func TestNewPoint_EmptyEnumAllowed(t *testing.T) {
	p, err := NewPoint(WithY(1)) // no Kind set
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Kind != "" {
		t.Fatalf("expected empty Kind, got %q", p.Kind)
	}
}
