package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Task struct {
	Identifier int
	Priority   int
}

type Scheduler struct {
	heap       []Task
	index      map[int]int
	priorities map[int]int
}

func NewScheduler() Scheduler {
	return Scheduler{
		heap:       make([]Task, 0),
		index:      make(map[int]int),
		priorities: make(map[int]int),
	}
}

func (s *Scheduler) AddTask(task Task) {
	if _, exists := s.index[task.Identifier]; exists {
		s.ChangeTaskPriority(task.Identifier, task.Priority)

		return
	}

	s.heap = append(s.heap, task)
	idx := len(s.heap) - 1
	s.index[task.Identifier] = idx
	s.priorities[task.Identifier] = task.Priority

	s.siftUp(idx)
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	idx, exists := s.index[taskID]
	if !exists {
		return
	}

	oldPriority := s.priorities[taskID]
	s.priorities[taskID] = newPriority

	if newPriority > oldPriority {
		s.siftUp(idx)
	} else if newPriority < oldPriority {
		s.siftDown(idx)
	}
}

func (s *Scheduler) GetTask() Task {
	if len(s.heap) == 0 {
		return Task{}
	}

	root := s.heap[0]
	delete(s.index, root.Identifier)
	delete(s.priorities, root.Identifier)

	lastIdx := len(s.heap) - 1
	if lastIdx > 0 {
		s.swap(0, lastIdx)
		s.heap = s.heap[:lastIdx]
		s.siftDown(0)
	} else {
		s.heap = s.heap[:0]
	}

	return root
}

func (s *Scheduler) siftUp(idx int) {
	for idx > 0 {
		parent := (idx - 1) / 2
		if s.priorities[s.heap[idx].Identifier] <= s.priorities[s.heap[parent].Identifier] {
			break
		}
		s.swap(idx, parent)
		idx = parent
	}
}

func (s *Scheduler) siftDown(idx int) {
	n := len(s.heap)
	for {
		left := 2*idx + 1
		right := 2*idx + 2
		largest := idx

		if left < n && s.priorities[s.heap[left].Identifier] > s.priorities[s.heap[largest].Identifier] {
			largest = left
		}
		if right < n && s.priorities[s.heap[right].Identifier] > s.priorities[s.heap[largest].Identifier] {
			largest = right
		}

		if largest == idx {
			break
		}

		s.swap(idx, largest)
		idx = largest
	}
}

func (s *Scheduler) swap(i, j int) {
	s.heap[i], s.heap[j] = s.heap[j], s.heap[i]
	s.index[s.heap[i].Identifier] = i
	s.index[s.heap[j].Identifier] = j
}

func TestTrace(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, task)

	scheduler.ChangeTaskPriority(1, 100)

	task = scheduler.GetTask()
	assert.Equal(t, task1, task)

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)
}
