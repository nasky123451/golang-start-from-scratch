package api

import (
	pq "github.com/jupp0r/go-priority-queue"
	"github.com/oleiade/lane"
)

// 使用 go-priority-queue 實作
func UseGoPriorityQueue() []string {
	pq := pq.New()
	pq.Insert("task1", 3.0)
	pq.Insert("task2", 1.0)
	pq.Insert("task3", 2.0)

	results := []string{}
	for pq.Len() > 0 {
		item, _ := pq.Pop()
		results = append(results, item.(string))
	}
	return results
}

// 使用 lane 實作
func UseLane() []string {
	pq := lane.NewPQueue(lane.MINPQ)
	pq.Push("task1", 3)
	pq.Push("task2", 1)
	pq.Push("task3", 2)

	results := []string{}
	for !pq.Empty() {
		item, _ := pq.Pop()
		results = append(results, item.(string))
	}
	return results
}
