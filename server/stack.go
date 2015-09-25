package main

import (
  "errors"  
)
  
var ErrStackEmpty = errors.New("Stack is empty")  

// Stack is the data structure representation of the classic stack
type Stack struct {
  top *Node
  size int
}

type node struct {
  value interface{}
  next *node
}

// Len return the size of the stack
func (s *Stack) Len() int {
  return s.size
}


// IsEmpty returns true if the size of the stack is 0.
func (s *Stack) IsEmpty() bool {
  return s.size == 0
}

// Push adds this value on top of the stack.
func (s *Stack) Push(val interface{}) {
  s.top = &node{val, s.top}
  s.size++
}

// Peek returns the value on top of the stack or err if the stack is empty.
func (s *Stack) Peek() (interface{}, error) {
  if s.size == 0{
    return nil, ErrStackEmpty
  }
  return s.top.value
}

// Pop removes and returns the value on top of the stack, or nil.
func (s *Stack) Pop() (interface{}, error) {
  if s.size == 0{
    return nil, ErrStackEmpty 
  }
  s.size--
  var val interface{}
  val, s.top = s.top.value, s.top.next
  return val, nil
}
