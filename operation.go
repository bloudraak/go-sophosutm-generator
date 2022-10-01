package main

type Operation struct {
    name               string
    extendedProperties map[string]interface{}
    parameters         []*Parameter
    returnTypes        []*ReturnType
}

func NewOperation(name string) *Operation {
    return &Operation{
        name:               name,
        extendedProperties: make(map[string]interface{}),
    }
}

func (o *Operation) Name() string {
    return o.name
}

func (o *Operation) SetName(name string) {
    o.name = name
}

func (o *Operation) WithName(name string) *Operation {
    o.name = name
    return o
}

func (o *Operation) WithExtendedProperty(key string, value interface{}) *Operation {
    o.extendedProperties[key] = value
    return o
}

func (o *Operation) CreateParameter(name string) *Parameter {
    p := NewParameter(name)
    o.parameters = append(o.parameters, p)
    return p
}

func (o *Operation) Parameters() []*Parameter {
    var parameters []*Parameter
    for _, p := range o.parameters {
        parameters = append(parameters, p)
    }
    return parameters
}

func (o *Operation) WithParameter(name string, parameterType string) *Operation {
    p := NewParameter(name)
    p.WithType(parameterType)
    o.parameters = append(o.parameters, p)
    return o
}

func (o *Operation) ReturnTypes() []*ReturnType {
    return o.returnTypes
}

func (o *Operation) WithReturnType(name string, returnType string) *Operation {
    r := NewReturnType()
    r.WithName(name)
    r.WithType(returnType)
    o.returnTypes = append(o.returnTypes, r)
    return o
}

func (o *Operation) CreateReturnType(name string, returnType string) *ReturnType {
    r := NewReturnType()
    r.WithName(name)
    r.WithType(returnType)
    o.returnTypes = append(o.returnTypes, r)
    return r
}

func (o *Operation) ExtendedProperty(key string) (interface{}, bool) {
    v, ok := o.extendedProperties[key]
    return v, ok
}

func (o *Operation) ExtendedProperties() map[string]interface{} {
    return o.extendedProperties
}
