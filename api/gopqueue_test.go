package api_test

import (
	"testing"

	"example.com/m/api"
)

func TestGopqueue(t *testing.T) {
	// 測試 1: 正常情況
	pq := api.NewConcurrentPriorityQueue()

	pq.Enqueue("task1", 3)
	pq.Enqueue("task2", 1)
	pq.Enqueue("task3", 2)

	expectedOrder := []string{"task2", "task3", "task1"}
	actualOrder := []string{}

	for pq.Len() > 0 {
		if value, ok := pq.Dequeue(); ok {
			actualOrder = append(actualOrder, value)
		} else {
			t.Errorf("Dequeue failed when it shouldn't")
		}
	}

	if len(actualOrder) != len(expectedOrder) {
		t.Fatalf("Expected %d tasks but got %d", len(expectedOrder), len(actualOrder))
	}

	for i, expected := range expectedOrder {
		if actualOrder[i] != expected {
			t.Errorf("Expected task %s at position %d but got %s", expected, i, actualOrder[i])
		}
	}

	// 測試 2: 空佇列情況
	if value, ok := pq.Dequeue(); ok {
		t.Errorf("Dequeue should fail on an empty queue, but got %s", value)
	}
	if pq.Len() != 0 {
		t.Errorf("Queue length should be 0 for an empty queue, but got %d", pq.Len())
	}
}
