package main

type Node struct {
	Value interface{}
	Next  *Node
}

type Stack struct {
	top *Node
}

func NewStack() *Stack {
	return &Stack{}
}

func (s *Stack) Push(value interface{}) {
	newNode := &Node{Value: value}
	newNode.Next = s.top
	s.top = newNode
}

func (s *Stack) Pop() interface{} {
	if s.IsEmpty() {
		panic("Empty")
	}
	value := s.top.Value
	newTop := s.top.Next
	s.top.Next = nil
	s.top = newTop
	return value
}

func (s *Stack) IsEmpty() bool {
	return s.top == nil
}

func (s *Stack) Peek() interface{} {
	if s.IsEmpty() {
		panic("Stack is empty")
	}
	return s.top.Value
}
