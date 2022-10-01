package main

type Stereotype struct {
	name string
}

func NewStereotype(name string) *Stereotype {
	return &Stereotype{name: name}
}

func (s *Stereotype) Name() string {
	return s.name
}

func (s *Stereotype) WithName(name string) *Stereotype {
	s.name = name
	return s
}
