package application

import (
	"labyrinths/internal/domain"
)

type Item struct {
	coord domain.Point
	priority int
	index int
}

type PriorityQueue []*Item

func (priorityQueue PriorityQueue) Len() int {
	return len(priorityQueue)
}

func (priorityQueue PriorityQueue) Less(i, j int) bool {
	return priorityQueue[i].priority < priorityQueue[j].priority
}

func (priorityQueue PriorityQueue) Swap(i, j int) {
	priorityQueue[i], priorityQueue[j] = priorityQueue[j], priorityQueue[i]
	priorityQueue[i].index = i
	priorityQueue[j].index = j
}

func (priorityQueue *PriorityQueue) Push(x interface{}) {
	n := len(*priorityQueue)
	item := x.(*Item)
	item.index = n
	*priorityQueue = append(*priorityQueue, item)
}

func (priorityQueue *PriorityQueue) Pop() interface{} {
	old := *priorityQueue
	n := len(old)
	item := old[n - 1]
	old[n - 1] = nil
	item.index = -1
	*priorityQueue = old[0:n-1]

	return item
}
