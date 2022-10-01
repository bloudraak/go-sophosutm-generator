package main

type Parameter struct {
    name               string
    parameterType      string
    extendedProperties map[string]interface{}
    required           bool
    lowerValue         int
    upperValue         int
}

func (p *Parameter) LowerValue() int {
    return p.lowerValue
}

func (p *Parameter) SetLowerValue(lowerValue int) {
    p.lowerValue = lowerValue
}

func (p *Parameter) WithLowerValue(lowerValue int) *Parameter {
    p.lowerValue = lowerValue
    return p
}

func (p *Parameter) UpperValue() int {
    return p.upperValue
}

func (p *Parameter) SetUpperValue(upperValue int) {
    p.upperValue = upperValue
}

func (p *Parameter) WithUpperValue(upperValue int) *Parameter {
    p.upperValue = upperValue
    return p
}

func (p *Parameter) Required() bool {
    return p.required
}

func (p *Parameter) SetRequired(required bool) {
    p.required = required
}

func (p *Parameter) Name() string {
    return p.name
}

func (p *Parameter) SetName(name string) {
    p.name = name
}

func NewParameter(name string) *Parameter {
    return &Parameter{
        name:               name,
        extendedProperties: make(map[string]interface{}),
    }
}

func (p *Parameter) WithName(name string) *Parameter {
    p.name = name
    return p
}

func (p *Parameter) WithExtendedProperty(key string, value interface{}) *Parameter {
    p.extendedProperties[key] = value
    return p
}

func (p *Parameter) ExtendedProperty(key string) (interface{}, bool) {
    v, ok := p.extendedProperties[key]
    return v, ok
}

func (p *Parameter) ExtendedProperties() map[string]interface{} {
    return p.extendedProperties
}

func (p *Parameter) WithType(typeName string) *Parameter {
    p.parameterType = typeName
    return p
}

func (p *Parameter) Type() string {
    return p.parameterType
}

func (p *Parameter) AsRequired(required bool) *Parameter {
    p.required = required
    return p
}
