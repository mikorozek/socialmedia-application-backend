package auth

import "testing"

func TestSimpleAddition(t *testing.T) {
    if 2+2 != 4 {
        t.Error("Basic math failed: 2 + 2 != 4")
    }
}

func TestSimpleEquality(t *testing.T) {
    expected := 2
    if expected != 2 {
        t.Errorf("Expected %d to equal 2", expected)
    }
}
