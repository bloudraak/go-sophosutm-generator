package main

type Model struct {
    classes     map[string]*Class
    stereotypes map[string]*Stereotype
}

func NewModel() *Model {
    return &Model{
        classes:     make(map[string]*Class),
        stereotypes: make(map[string]*Stereotype),
    }
}

func (m *Model) Classes() Classes {
    var classes Classes
    for _, c := range m.classes {
        classes = append(classes, c)
    }
    return classes
}

func (m *Model) CreateClass(name string) *Class {
    c := NewClass(name)
    m.classes[name] = c
    return c
}

func (m *Model) CreateStereotype(name string) *Stereotype {
    s := NewStereotype(name)
    m.stereotypes[name] = s
    return s
}

func (m *Model) Stereotypes() []*Stereotype {
    var stereotypes []*Stereotype
    for _, s := range m.stereotypes {
        stereotypes = append(stereotypes, s)
    }
    return stereotypes
}

func (m *Model) WithClass(name string) *Model {
    m.classes[name] = NewClass(name)
    return m
}

func (m *Model) WithStereotype(name string) *Model {
    m.stereotypes[name] = NewStereotype(name)
    return m
}

func (m *Model) Class(name string) *Class {
    if c, ok := m.classes[name]; ok {
        return c
    }
    return nil
}

func (m *Model) Stereotype(name string) *Stereotype {
    if s, ok := m.stereotypes[name]; ok {
        return s
    }
    return nil
}

func (m *Model) HasClass(name string) bool {
    _, ok := m.classes[name]
    return ok
}
