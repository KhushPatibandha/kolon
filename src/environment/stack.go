package environment

type Stack struct {
	Stk []*Environment
}

func NewStack() *Stack { return &Stack{Stk: []*Environment{}} }

func (s *Stack) Push(env *Environment) { s.Stk = append(s.Stk, env) }

func (s *Stack) Pop() *Environment {
	if len(s.Stk) == 1 {
		return nil
	}
	env := s.Stk[len(s.Stk)-1]
	s.Stk = s.Stk[:len(s.Stk)-1]
	return env
}

func (s *Stack) Top() *Environment { return s.Stk[len(s.Stk)-1] }

func (s *Stack) Len() int { return len(s.Stk) }
