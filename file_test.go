package bget

import "testing"

func TestHumanSize(t *testing.T) {
	v := HumanSize(1044)
	if v != "1 KB" {
		t.Fatalf("HumanSize(1044) = %s", v)
	}
}
