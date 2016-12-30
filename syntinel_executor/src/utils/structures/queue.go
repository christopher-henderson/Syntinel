package structures

import "container/list"

type Queue struct {
	list list.List
}

func NewQueue() Queue {
	return Queue{list.List{}}
}

func (q *Queue) Push(element interface{}) error {
	q.list.PushBack(element)
	return nil
}

func (q *Queue) Pop() interface{} {
	element := q.list.Front()
	if element == nil {
		return nil
	}
	q.list.Remove(element)
	return element.Value
}

func (q *Queue) Peek() interface{} {
	element := q.list.Front()
	if element == nil {
		return nil
	}
	return element.Value
}

func (q *Queue) Len() int {
	return q.list.Len()
}
