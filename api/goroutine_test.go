package api

import (
	"testing"
)

// TestSharingWithAtomic tests the sharingWithAtomic function.
func TestSharingWithAtomic(t *testing.T) {
	result := sharingWithAtomic()

	if result < 0 || result >= 300 {
		t.Errorf("sharingWithAtomic() returned %d; expected result in range 0-299", result)
	}
}

// TestSharingWithMutex tests the sharingWithMutex function.
func TestSharingWithMutex(t *testing.T) {
	result := sharingWithMutex()

	if result < 0 || result >= 300 {
		t.Errorf("sharingWithMutex() returned %d; expected result in range 0-299", result)
	}
}

// TestSharingWithChannel tests the sharingWithChannel function.
func TestSharingWithChannel(t *testing.T) {
	result := sharingWithChannel()

	if result < 0 || result >= 300 {
		t.Errorf("sharingWithChannel() returned %d; expected result in range 0-299", result)
	}
}
