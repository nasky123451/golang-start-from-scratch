package api_test

import (
	"testing"

	"example.com/m/api"
)

// 測試所有實作是否正確按優先級排序
func TestPriorityQueues(t *testing.T) {
	expected := []string{"task2", "task3", "task1"}

	// 測試 go-priority-queue
	results := api.UseGoPriorityQueue()
	if !equal(results, expected) {
		t.Errorf("go-priority-queue: expected %v, got %v", expected, results)
	}

	// 測試 lane
	results = api.UseLane()
	if !equal(results, expected) {
		t.Errorf("lane: expected %v, got %v", expected, results)
	}
}

// Helper function: 比較兩個 slice 是否相等
func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
