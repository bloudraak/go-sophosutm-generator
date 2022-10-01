package main

type ReturnType struct {
    name               string
    extendedProperties map[string]interface{}
    returnType         string
    upperValue         int
    lowerValue         int
}

func (r *ReturnType) Type() string {
    return r.returnType
}

func (r *ReturnType) SetType(returnType string) {
    r.returnType = returnType
}

func (r *ReturnType) WithType(returnType string) *ReturnType {
    r.returnType = returnType
    return r
}

func (r *ReturnType) WithExtendedProperty(key string, value interface{}) *ReturnType {
    r.extendedProperties[key] = value
    return r
}

func (r *ReturnType) ExtendedProperty(key string) (interface{}, bool) {
    v, ok := r.extendedProperties[key]
    return v, ok
}

func (r *ReturnType) ExtendedProperties() map[string]interface{} {
    return r.extendedProperties
}

func NewReturnType() *ReturnType {
    return &ReturnType{
        extendedProperties: make(map[string]interface{}),
    }
}

func (r *ReturnType) Name() string {
    return r.name
}

func (r *ReturnType) SetName(name string) {
    r.name = name
}

func (r *ReturnType) WithName(name string) *ReturnType {
    r.name = name
    return r
}

func (r *ReturnType) LowerValue() int {
    return r.lowerValue
}

func (r *ReturnType) SetLowerValue(lowerValue int) {
    r.lowerValue = lowerValue
}

func (r *ReturnType) WithLowerValue(lowerValue int) *ReturnType {
    r.lowerValue = lowerValue
    return r
}

func (r *ReturnType) UpperValue() int {
    return r.upperValue
}

func (r *ReturnType) SetUpperValue(upperValue int) {
    r.upperValue = upperValue
}

func (r *ReturnType) WithUpperValue(upperValue int) *ReturnType {
    r.upperValue = upperValue
    return r
}
