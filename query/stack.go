package query

type stack struct {
	top  *element
	size int
}

type element struct {
	node *ConditionNode
	next *element
}

func (s *stack) len() int {
	return s.size
}

func (s *stack) push(node *ConditionNode) {
	s.top = &element{node, s.top}
	s.size++
}

func (s *stack) pop() (node *ConditionNode) {
	if s.size == 0 {
		return nil
	}

	node, s.top = s.top.node, s.top.next
	s.size--
	return node
}
