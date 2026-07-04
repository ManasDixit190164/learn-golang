package mathpkg

import "testing"

func TestAdd(t *testing.T) {
    got := Add(2, 3)
    want := 5

    if got != want {
        t.Fatalf("Add(2, 3) = %d; want %d", got, want)
    }
}

func TestMultiply(t *testing.T) {
    got := Multiply(2, 3)
    want := 6

    if got != want {
        t.Fatalf("Multiply(2, 3) = %d; want %d", got, want)
    }
}
