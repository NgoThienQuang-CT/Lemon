package value

type ScopeErr int

const (
	OK = iota
	UNDEFINED
	IMMUTABLE
)

type Binding struct {
	Inner   Value
	Mutable bool
}

type Scope struct {
	store map[string]Binding
	outer *Scope
}

func NewScope() *Scope {
	store := make(map[string]Binding)
	return &Scope{store: store, outer: nil}
}

func NewEnclosedScope(outer *Scope) *Scope {
	scope := NewScope()
	scope.outer = outer
	return scope
}

func (s *Scope) Get(name string) (Value, bool) {
	binding, ok := s.store[name]
	if ok {
		return binding.Inner, ok
	}

	if s.outer != nil {
		return s.outer.Get(name)
	}

	return nil, ok
}

func (s *Scope) Set(name string, val Value) (Value, ScopeErr) {
	binding, ok := s.store[name]

	if !ok {
		if s.outer != nil {
			return s.outer.Set(name, val)
		}
		return nil, UNDEFINED
	}

	if !binding.Mutable {
		return nil, IMMUTABLE
	}

	s.store[name] = Binding{Inner: val, Mutable: true}
	return val, OK
}

func (s *Scope) Define(name string, mutable bool, val Value) (Value, bool) {
	binding, ok := s.store[name]

	if ok && binding.Mutable {
		return nil, false
	}

	s.store[name] = Binding{Inner: val, Mutable: mutable}
	return val, true
}
