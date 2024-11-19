package api

import (
	"container/heap"
	"fmt"
	"sync"
)

// Item 儲存佇列元素和優先級
type Item struct {
	Value    string // 元素值
	Priority int    // 優先級（數字越小優先級越高）
	Index    int    // 元素索引（由 heap.Interface 管理）
}

// PriorityQueue 實現一個佇列
type PriorityQueue []*Item

// Len 取得佇列長度
func (pq PriorityQueue) Len() int { return len(pq) }

// Less 定義優先順序
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

// Swap 交換元素
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

// Push 插入元素
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.Index = n
	*pq = append(*pq, item)
}

// Pop 移除最高優先級元素
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // 避免記憶體洩漏
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}

// ConcurrentPriorityQueue 支援併發操作的優先佇列
type ConcurrentPriorityQueue struct {
	queue PriorityQueue
	lock  sync.Mutex
}

// NewConcurrentPriorityQueue 創建新的佇列
func NewConcurrentPriorityQueue() *ConcurrentPriorityQueue {
	return &ConcurrentPriorityQueue{queue: make(PriorityQueue, 0)}
}

// Enqueue 插入元素
func (cpq *ConcurrentPriorityQueue) Enqueue(value string, priority int) {
	cpq.lock.Lock()
	defer cpq.lock.Unlock()

	heap.Push(&cpq.queue, &Item{Value: value, Priority: priority})
}

// Dequeue 移除最高優先級元素
func (cpq *ConcurrentPriorityQueue) Dequeue() (string, bool) {
	cpq.lock.Lock()
	defer cpq.lock.Unlock()

	if cpq.queue.Len() == 0 {
		return "", false
	}

	item := heap.Pop(&cpq.queue).(*Item)
	return item.Value, true
}

// Len 取得佇列長度
func (cpq *ConcurrentPriorityQueue) Len() int {
	cpq.lock.Lock()
	defer cpq.lock.Unlock()

	return cpq.queue.Len()
}

// 主函數測試
func main() {
	pq := NewConcurrentPriorityQueue()

	pq.Enqueue("task1", 3)
	pq.Enqueue("task2", 1)
	pq.Enqueue("task3", 2)

	for pq.Len() > 0 {
		if value, ok := pq.Dequeue(); ok {
			fmt.Println(value)
		}
	}
}
