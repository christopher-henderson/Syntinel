package tests

import (
	"syntinel_executor/utils/structures"
	"testing"
)

func TestQueuePush(t *testing.T) {
	q := structures.NewQueue()
	value := 5
	q.Push(value)
}

func TestQueuePop(t *testing.T) {
	q := structures.NewQueue()
	value := 5
	q.Push(value)
	if v := q.Pop(); v != value {
		t.Errorf("Expected %v from popping the queue, got %v", value, v)
	}
}

func TestQueuePopEmpty(t *testing.T) {
	q := structures.NewQueue()
	if v := q.Pop(); v != nil {
		t.Errorf("Expected nil from popping an empty queue, got %v", v)
	}
}

func TestQueuePushPopPop(t *testing.T) {
	q := structures.NewQueue()
	value := 5
	q.Push(5)
	if v := q.Pop(); v != value {
		t.Errorf("Expected %v from popping the queue, got %v", value, v)
	}
	if v := q.Pop(); v != nil {
		t.Errorf("Expected nil from popping an empty queue, got %v", v)
	}
}

func TestQueuePeek(t *testing.T) {
	q := structures.NewQueue()
	value := 5
	q.Push(5)
	if v := q.Peek(); v != value {
		t.Errorf("Expected %v from popping the queue, got %v", value, v)
	}
	// Peek again to make sure it wasn't removed by peeking.
	if v := q.Peek(); v != value {
		t.Errorf("Expected %v from popping the queue, got %v", value, v)
	}
}

func TestQueueLen(t *testing.T) {
	q := structures.NewQueue()
	if len := q.Len(); len != 0 {
		t.Errorf("Expected 0 elements in the queue, got %v", len)
	}
	q.Push(5)
	if len := q.Len(); len != 1 {
		t.Errorf("Expected 1 element in the queue, got %v", len)
	}
	q.Push(5)
	if len := q.Len(); len != 2 {
		t.Errorf("Expected 2 elements in the queue, got %v", len)
	}
	q.Pop()
	if len := q.Len(); len != 1 {
		t.Errorf("Expected 1 element in the queue, got %v", len)
	}
	q.Pop()
	if len := q.Len(); len != 0 {
		t.Errorf("Expected 0 element in the queue, got %v", len)
	}
}
