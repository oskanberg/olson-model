package main

// Node represents a node in the PLGMN
type Node struct {
	t, tNext bool
	id       int
}

func (s *Node) GetState() bool {
	return s.t
}

func (s *Node) SetState(state bool) {
	if !s.tNext {
		s.tNext = state
	}
}

func (s *Node) Step() {
	s.t = s.tNext
}

func (s *Node) SetId(id int) {
	s.id = id
}

func (s *Node) GetId() int {
	return s.id
}
