package main

import (
	"testing"
)

func TestStack_PushPop_LinkedList(t *testing.T) {
	s := NewStack()
	s.Push(1)
	s.Push(2)
	s.Push(3)

	if s.top.Value != 3 {
		t.Errorf("Incorrect top element: expected %v, got %v", 3, s.top.Value)
	}

	value1 := s.Pop()
	if value1 != 3 {
		t.Errorf("Incorrect value popped: expected %v, got %v", 3, value1)
	}

	value2 := s.Pop()
	if value2 != 2 {
		t.Errorf("Incorrect value popped: expected %v, got %v", 2, value2)
	}

	value3 := s.Pop()
	if value3 != 1 {
		t.Errorf("Incorrect value popped: expected %v, got %v", 1, value2)
	}

	if s.top != nil {
		t.Errorf("Top pointer should be nil after popping all elements")
	}
}

func TestStack_IsEmpty_LinkedList(t *testing.T) {
	s := NewStack()
	if !s.IsEmpty() {
		t.Errorf("Stack should be empty by default")
	}

	s.Push(1)
	if s.IsEmpty() {
		t.Errorf("Stack should not be empty after pushing an element")
	}
}

func TestStack_Peek_LinkedList(t *testing.T) {
	s := NewStack()
	s.Push(1)
	s.Push(2)

	if s.Peek() != 2 {
		t.Errorf("Incorrect value peeked: expected %v, got %v", 2, s.Peek())
	}

	if s.top.Value != 2 {
		t.Errorf("Top element modified after peek: expected %v, got %v", 2, s.top.Value)
	}
}
