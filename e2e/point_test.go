package e2e

import "testing"

func TestNewPoint_Defaults(t *testing.T) {
	p := NewPoint()
	if p.X != 0 || p.Y != 0 || p.Name != "origin" {
		t.Fatalf("defaults wrong: %+v", p)
	}
}

func TestNewPoint_Options(t *testing.T) {
	p := NewPoint(
		WithX(3),
		WithY(4),
		WithName("hello"),
		WithTags([]string{"a", "b"}),
	)
	if p.X != 3 || p.Y != 4 || p.Name != "hello" {
		t.Fatalf("options not applied: %+v", p)
	}
	if len(p.Tags) != 2 || p.Tags[0] != "a" {
		t.Fatalf("tags wrong: %+v", p.Tags)
	}
	if p.ID != 0 {
		t.Fatalf("readonly ID should be zero: %v", p.ID)
	}
}
