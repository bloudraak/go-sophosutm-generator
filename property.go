package main

type Property struct {
    name               string
    propertyType       string
    required           bool
    extendedProperties map[string]interface{}
    lowerValue         int
    upperValue         int
    stereotypes        map[string]*Stereotype
}

func (p *Property) LowerValue() int {
    return p.lowerValue
}

func (p *Property) SetLowerValue(lowerValue int) {
    p.lowerValue = lowerValue
}

func (p *Property) WithLowerValue(lowerValue int) *Property {
    p.lowerValue = lowerValue
    return p
}

func (p *Property) UpperValue() int {
    return p.upperValue
}

func (p *Property) SetUpperValue(upperValue int) {
    p.upperValue = upperValue
}

func (p *Property) WithUpperValue(upperValue int) *Property {
    p.upperValue = upperValue
    return p
}

func (p *Property) Required() bool {
    return p.required
}

func (p *Property) SetRequired(required bool) {
    p.required = required
}

func (p *Property) AsRequired() *Property {
    p.required = true
    return p
}

func (p *Property) AsOptional() *Property {
    p.required = true
    return p
}

func (p *Property) Type() string {
    return p.propertyType
}

func (p *Property) SetType(kind string) {
    p.propertyType = kind
}

func (p *Property) Name() string {
    return p.name
}

func (p *Property) SetName(name string) {
    p.name = name
}

func (p *Property) WithName(name string) *Property {
    p.name = name
    return p
}

func (p *Property) WithType(s string) *Property {
    p.propertyType = s
    return p
}

func (p *Property) WithExtendedProperty(key string, value interface{}) *Property {
    p.extendedProperties[key] = value
    return p
}

func (p *Property) HasStereoType(s string) bool {
    _, ok := p.stereotypes[s]
    return ok
}

func (p *Property) WithStereoType(s *Stereotype) *Property {
    p.stereotypes[s.Name()] = s
    return p
}

func NewProperty(name string) *Property {
    return &Property{
        name:               name,
        required:           false,
        propertyType:       "string",
        extendedProperties: make(map[string]interface{}),
        stereotypes:        make(map[string]*Stereotype),
    }
}
