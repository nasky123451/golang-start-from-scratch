package api

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

// Helper function to create a Block with a given start and end time
func newBlock(start, end time.Time) Block {
	return Block{
		id:    uuid.New(),
		start: start,
		end:   end,
	}
}

func TestCompact(t *testing.T) {
	// Set up blocks with different start and end times
	block1 := newBlock(time.Date(2023, 9, 24, 10, 0, 0, 0, time.UTC), time.Date(2023, 9, 24, 12, 0, 0, 0, time.UTC))
	block2 := newBlock(time.Date(2023, 9, 24, 8, 0, 0, 0, time.UTC), time.Date(2023, 9, 24, 9, 0, 0, 0, time.UTC))
	block3 := newBlock(time.Date(2023, 9, 24, 11, 0, 0, 0, time.UTC), time.Date(2023, 9, 24, 13, 0, 0, 0, time.UTC))

	// Call Compact function
	compactBlock := Compact(block1, block2, block3)

	// Check that the start and end times are correct
	expectedStart := block2.start // The earliest start time
	expectedEnd := block3.end     // The latest end time

	if !compactBlock.start.Equal(expectedStart) {
		t.Errorf("expected start time %v, got %v", expectedStart, compactBlock.start)
	}

	if !compactBlock.end.Equal(expectedEnd) {
		t.Errorf("expected end time %v, got %v", expectedEnd, compactBlock.end)
	}
}

func TestEmptyCompact(t *testing.T) {
	// Test Compact with no blocks
	compactBlock := Compact()

	// Check that the start and end times are zero
	if !compactBlock.start.IsZero() {
		t.Errorf("expected zero start time, got %v", compactBlock.start)
	}

	if !compactBlock.end.IsZero() {
		t.Errorf("expected zero end time, got %v", compactBlock.end)
	}
}
