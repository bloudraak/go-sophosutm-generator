package main

import (
    "log"
    "regexp"
    "strings"

    "github.com/getkin/kin-openapi/openapi2"
    "github.com/getkin/kin-openapi/openapi3"
)

var renameMap = map[string]string{
    "1to1nat": "one to one nat",
}

type ModelBuilder struct {
    model        *Model
    skipPatterns []string
}

func (b *ModelBuilder) Build(m map[*RemoteSchemaReference]*openapi2.T, only string) *Model {
    schemaMap := m
    b.buildStereoTypes()

    clientClass := b.model.CreateClass("Client").
        WithStereotype(b.model.Stereotype(clientStereoType)).
        WithPackageImport("crypto/tls").
        WithPackageImport("encoding/base64").
        WithPackageImport("net/http")

    const operationsSuffix = " operations"
    for k, v := range schemaMap {
        if only != "" {
            match, err := regexp.MatchString(only, k.Name)
            if err != nil {
                log.Fatal(err)
            }
            if !match {
                continue
            }
        }

        b.buildClassesFromDefinitions(v)

        operationsClass := b.model.CreateClass(k.Name+operationsSuffix).
            WithExtendedProperty("schema", v).
            WithExtendedProperty("name", k).
            WithStereotype(b.model.Stereotype(dataOperationsStereoType))

        clientClass.CreateProperty(k.Name + operationsSuffix).
            WithType(operationsClass.Name())

        for path, pathItem := range v.Paths {
            if strings.HasSuffix(path, "usedby") {
                continue
            }
            for method, operation := range pathItem.Operations() {
                className := operation.Tags[0]
                className = strings.ReplaceAll(className, "/", " ")
                className = strings.ReplaceAll(className, ".", " ")

                for k, v := range renameMap {
                    className = strings.ReplaceAll(className, k, v)
                }

                className = className + operationsSuffix

                class := b.model.Class(className)
                if class == nil {
                    class = b.model.CreateClass(className).
                        WithExtendedProperty("schema", v.Paths[path]).
                        WithExtendedProperty("path", className).
                        WithStereotype(b.model.Stereotype(dataOperationsStereoType)).
                        WithPackageImport("bytes").
                        WithPackageImport("encoding/json").
                        WithPackageImport("net/http").
                        WithPackageImport("strings").
                        WithPackageImport("errors").
                        WithPackageImport("strconv")

                    operationsClass.CreateProperty(className).
                        WithType(className)
                }
                var operationName string
                switch strings.ToLower(method) {
                case "get":
                    if strings.HasSuffix(path, "usedby") {
                        operationName = "UsedBy"
                    } else {
                        response := findSuccessfulResponse(operation)
                        if response == nil {
                            log.Fatal("No successful response found for GET operation")
                        }
                        if returnsArray(response) {
                            operationName = "List"
                        } else {
                            operationName = "Get"
                        }

                        if operationName == "" {
                            log.Fatalln("No operation name found for", method, path)
                        }
                    }

                case "post":
                    operationName = "Create"
                case "put":
                    operationName = "Update"
                case "delete":
                    operationName = "Delete"
                case "patch":
                    operationName = "Patch"
                default:
                    operationName = "Unknown"
                }

                o := class.CreateOperation(operationName).
                    WithExtendedProperty("schema", operation).
                    WithExtendedProperty("path", join(v.BasePath, path)).
                    WithExtendedProperty("method", method)

                response := findSuccessfulResponse(operation)
                if response != nil {
                    lowerValue := 1
                    upperValue := 1
                    returnType := "unknown"
                    if response.Schema == nil {
                        continue
                    }
                    if response.Schema.Ref != "" {
                        returnType = b.trimDefinitionsPrefix(response.Schema.Ref)
                    } else if response.Schema.Value != nil {
                        switch response.Schema.Value.Type {
                        case "array":
                            lowerValue = 0
                            upperValue = -1
                            if response.Schema.Value.Items.Ref != "" {
                                returnType = b.trimDefinitionsPrefix(response.Schema.Value.Items.Ref)
                            } else {
                                returnType = response.Schema.Value.Items.Value.Type
                            }
                        case "object":
                            clsName := path
                            clsName = strings.ReplaceAll(clsName, "/", " ")
                            cls := b.buildClassFromSchema(clsName, response.Schema, false)
                            returnType = cls.Name()
                        default:
                            returnType = response.Schema.Value.Type
                        }
                    }
                    o.CreateReturnType("result", returnType).
                        WithLowerValue(lowerValue).
                        WithUpperValue(upperValue)
                }

                o.WithReturnType("err", "error")

                for _, param := range operation.Parameters {
                    parameterType := "String"
                    if param.Schema != nil {
                        if param.Schema.Ref != "" {
                            parameterType = b.trimDefinitionsPrefix(param.Schema.Ref)
                        } else {
                            parameterType = param.Schema.Value.Type
                        }
                    } else {
                        parameterType = param.Type
                    }

                    paramName := strings.TrimPrefix(param.Name, "X-Restd-")

                    o.CreateParameter(paramName).
                        WithExtendedProperty("schema", param).
                        WithType(parameterType).
                        AsRequired(param.Required)
                }

            }
        }

    }

    for _, c := range b.model.Classes() {
        extendedProperty, exists := c.ExtendedProperty("schema")
        if !exists {
            continue
        }

        if c.HasStereoType(dataOperationsStereoType) {

        } else if c.HasStereoType(dataSpecStereoType) {

            schema, ok := extendedProperty.(*openapi3.SchemaRef)
            if !ok {
                log.Fatalln("not a schema")
            }

            for propertyName, property := range schema.Value.Properties {
                typeName := "String"
                lowerValue := 1
                upperValue := 1
                if property.Value != nil {
                    // TODO: base64 encoded binary data
                    switch property.Value.Type {
                    case "array":
                        lowerValue = 0
                        upperValue = -1
                        if property.Value.Items != nil {
                            if property.Value.Items.Ref != "" {
                                typeName = property.Value.Items.Ref
                            } else {
                                typeName = property.Value.Items.Value.Type
                            }
                        } else {
                            typeName = property.Value.Items.Ref
                        }
                    case "object":
                        typeName = "interface{}"
                    default:
                        typeName = property.Value.Type
                    }
                } else {
                    typeName = property.Ref
                }

                c.CreateProperty(propertyName).
                    WithType(typeName).
                    WithLowerValue(lowerValue).
                    WithUpperValue(upperValue).
                    AsOptional().
                    WithExtendedProperty("schema", property).
                    WithStereoType(b.model.Stereotype(serializableStereoType))
            }
        }
    }

    for _, c := range b.model.Classes() {
        name := c.Name()
        if c.HasStereoType(dataOperationsStereoType) {
            if !strings.HasSuffix(name, operationsSuffix) {
                name = name + operationsSuffix
            }
        }
        c.SetName(name)
    }
    return b.model
}

func (b *ModelBuilder) shouldSkip(k *RemoteSchemaReference) bool {
    var shouldSkip bool
    if b.skipPatterns != nil {
        for _, s := range b.skipPatterns {
            match, err := regexp.MatchString(s, k.Name)
            if err != nil {
                log.Fatal(err)
            }
            if match {
                shouldSkip = true
                break
            }
        }
    }
    return shouldSkip
}

func (b *ModelBuilder) trimDefinitionsPrefix(returnType string) string {
    if strings.HasPrefix(returnType, "#/definitions/") {
        returnType = strings.TrimPrefix(returnType, "#/definitions/")
    }
    return returnType
}

func (b *ModelBuilder) buildClassesFromDefinitions(v *openapi2.T) {
    for name, definition := range v.Definitions {
        b.buildClassFromSchema(name, definition, true)

    }
}

func (b *ModelBuilder) buildClassFromSchema(name string, definition *openapi3.SchemaRef, isDefinition bool) *Class {
    c := b.model.CreateClass(name).
        WithExtendedProperty("schema", definition).
        WithExtendedProperty("name", name).
        WithStereotype(b.model.Stereotype(dataSpecStereoType)).
        WithStereotype(b.model.Stereotype(serializableStereoType)).
        WithPackageImport("encoding/json").
        WithPackageImport("log")

    // if isDefinition {
    //     c.CreateProperty("_ref").
    //         WithType("string").
    //         WithStereoType(b.model.Stereotype(serializableStereoType))
    //     c.CreateProperty("_locked").
    //         WithType("string").
    //         WithStereoType(b.model.Stereotype(serializableStereoType))
    //     c.CreateProperty("_type").
    //         WithType("string").
    //         WithStereoType(b.model.Stereotype(serializableStereoType))
    // }

    return c

}

func (b *ModelBuilder) buildStereoTypes() *Model {
    return b.model.WithStereotype(dataOperationsStereoType).
        WithStereotype(dataSpecStereoType).
        WithStereotype(serializableStereoType).
        WithStereotype(clientStereoType)
}

func NewModelBuilder() *ModelBuilder {
    return &ModelBuilder{
        model: NewModel(),
    }
}
