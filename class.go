package main

type Class struct {
    name               string
    extendedProperties map[string]interface{}
    stereotypes        []*Stereotype
    properties         map[string]*Property
    operations         map[string]*Operation
    imports            []*PackageImport
}

func (c *Class) Properties() []*Property {
    var properties []*Property
    for _, p := range c.properties {
        properties = append(properties, p)
    }
    return properties
}

func NewClass(name string) *Class {
    return &Class{
        name:               name,
        extendedProperties: make(map[string]interface{}),
        properties:         make(map[string]*Property),
        operations:         make(map[string]*Operation),
    }
}

func (c *Class) ExtendedProperty(key string) (interface{}, bool) {
    v, ok := c.extendedProperties[key]
    return v, ok
}

func (c *Class) ExtendedProperties() map[string]interface{} {
    return c.extendedProperties
}

func (c *Class) Name() string {
    return c.name
}

func (c *Class) SetName(name string) {
    c.name = name
}

func (c *Class) WithName(name string) *Class {
    c.name = name
    return c
}
func (c *Class) WithExtendedProperty(key string, value interface{}) *Class {
    c.extendedProperties[key] = value
    return c
}

func (c *Class) WithProperty(name string) *Class {
    c.properties[name] = NewProperty(name)
    return c
}

func (c *Class) CreateProperty(name string) *Property {
    p := NewProperty(name)
    c.properties[name] = p
    return p
}

func (c *Class) HasStereoType(name string) bool {
    for _, s := range c.stereotypes {
        if s.Name() == name {
            return true
        }
    }
    return false
}

func (c *Class) Stereotypes() []*Stereotype {
    return c.stereotypes
}

func (c *Class) WithStereotype(s *Stereotype) *Class {
    c.stereotypes = append(c.stereotypes, s)
    return c
}

func (c *Class) CreateOperation(name string) *Operation {
    o := NewOperation(name)
    c.operations[name] = o
    return o
}

func (c *Class) Operations() []*Operation {
    var operations []*Operation
    for _, o := range c.operations {
        operations = append(operations, o)
    }
    return operations
}

func (c *Class) WithPackageImport(s string) *Class {
    c.imports = append(c.imports, NewPackageImport(c, s))
    return c
}

func (c *Class) PackageImports() []*PackageImport {
    return c.imports
}
