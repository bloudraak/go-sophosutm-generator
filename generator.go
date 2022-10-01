package main

import (
    "log"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "strings"

    "github.com/getkin/kin-openapi/openapi2"
    "github.com/iancoleman/strcase"
)

type GoGeneratorOptions struct {
    PackageName string
}

type GoGenerator struct {
    writer      *IndentedTextWriter
    path        string
    packageName string
}

func (g *GoGenerator) Writer() *IndentedTextWriter {
    if g.writer == nil {
        log.Fatalln("Writer is not initialized")
    }
    return g.writer
}

func (g *GoGenerator) generateClass(class *Class) {
    w := g.writer
    w.Println("package ", g.packageName)
    w.Println()
    if len(class.PackageImports()) > 0 {
        w.Println("import (")
        w.Indent()
        for _, p := range class.PackageImports() {
            w.Println(strconv.Quote(p.importedPackage))
        }
        w.Outdent()
        w.Println(")")
        w.Println()
    }
    if class.HasStereoType(clientStereoType) {
        w.Println("var DefaultUserAgent = \"go-sophosutm\"")
        w.Println()
    }
    name := class.Name()
    w.Println("type ", g.toTypeName(name), " struct {")
    w.Indent()
    for _, p := range class.Properties() {
        if p.HasStereoType(serializableStereoType) {
            w.Print(g.toPublicFieldName(p.Name()), " ")
        } else {
            w.Print(g.toFieldName(p.Name()), " ")
        }
        if p.LowerValue() != p.UpperValue() {
            w.Print("[]")
        }
        if g.ptr(p.Type(), p.Required()) {
            w.Print("*")
        }
        w.Print(g.toGoType(p.Type()))

        if p.HasStereoType(serializableStereoType) {
            w.Print(" `json:\"" + p.Name() + ",omitempty\"`")
        }
        w.Println()
    }

    if class.HasStereoType(clientStereoType) {
        w.Println("baseUri string")
        w.Println("httpClient *http.Client")
        w.Println("userAgent string")
        w.Println("credentials Credentials")
    }

    if class.HasStereoType(dataOperationsStereoType) {
        w.Println("client *Client")
    }

    w.Outdent()
    w.Println("}")
    w.Println()

    if class.HasStereoType(dataOperationsStereoType) {
        w.Print("func new", g.toTypeName(name))
        w.Print("(client *Client)")
        w.Print("(*", g.toTypeName(name), ", error) {")
        w.Println()
        w.Indent()
        if len(class.Properties()) > 0 {
            w.Println("var err error")
        }
        w.Println("result := new(", g.toTypeName(name), ")")

        w.Println("result.client = client")
        for _, p := range class.Properties() {
            w.Print("result.", g.toFieldName(p.Name()))
            if strings.HasSuffix(p.Type(), " operations") {
                w.Println(", err = new", g.toGoType(p.Type()), "(client)")
                w.Println("if err != nil {")
                w.Indent()
                w.Println("return nil, err")
                w.Outdent()
                w.Println("}")
            } else {
                w.Println(" = {}")
            }
        }
        w.Println("return result, nil")
        w.Outdent()
        w.Println("}")
        w.Println()
    }

    if class.HasStereoType(dataSpecStereoType) {
        w.Println(`func (c *`, g.toTypeName(name), `) String() string {`)
        w.Indent()
        w.Println(`b, err := json.MarshalIndent(c, "", "    ")`)
        w.Println(`if err != nil {`)
        w.Indent()
        w.Println(`log.Fatalln("Failed to marshal `, g.toTypeName(name), `", err)`)
        w.Outdent()
        w.Println(`}`)
        w.Println(`return string(b)`)
        w.Outdent()
        w.Println(`}`)
        w.Println()
    }

    if class.HasStereoType(clientStereoType) {

        w.Println("type Credentials interface {")
        w.Indent()
        w.Println("GetAuthorizationHeaderValue() string")
        w.Outdent()
        w.Println("}")
        w.Println()
        w.Println("type TokenCredentials struct {")
        w.Indent()
        w.Println("Token string")
        w.Outdent()
        w.Println("}")
        w.Println()
        w.Println("func (c *TokenCredentials) GetAuthorizationHeaderValue() string {")
        w.Indent()
        w.Println(`return "Basic " + base64.StdEncoding.EncodeToString([]byte("token:" + c.Token))`)
        w.Outdent()
        w.Println("}")
        w.Println()
        w.Println("type UsernamePasswordCredentials struct {")
        w.Indent()
        w.Println("Username string")
        w.Println("Password string")
        w.Outdent()
        w.Println("}")
        w.Println()
        w.Println("func (c *UsernamePasswordCredentials) GetAuthorizationHeaderValue() string {")
        w.Indent()
        w.Println(`return "Basic " + base64.StdEncoding.EncodeToString([]byte(c.Username + ":" + c.Password))`)
        w.Outdent()
        w.Println("}")
        w.Println()

        w.Print("func new", g.toTypeName(name))
        w.Print("(credentials Credentials, baseUri string, userAgent string, insecureSkipVerify bool)")
        w.Print("(*", g.toTypeName(name), ", error) {")
        w.Println()
        w.Indent()
        w.Println("var err error")
        w.Println("client := new(", g.toTypeName(name), ")")
        w.Println("client.baseUri = baseUri")
        w.Println("client.credentials = credentials")
        w.Println("transport := &http.Transport{")
        w.Indent()
        w.Println("TLSClientConfig: &tls.Config{")
        w.Indent()
        w.Println("InsecureSkipVerify: insecureSkipVerify,")
        w.Outdent()
        w.Println("},")
        w.Outdent()
        w.Println("}")
        w.Println("client.httpClient = &http.Client{")
        w.Indent()
        w.Println("Transport:     transport,")
        w.Println("CheckRedirect: nil,")
        w.Println("Jar:           nil,")
        w.Println("Timeout:       0,")
        w.Outdent()
        w.Println("}")
        w.Println("client.userAgent = userAgent")
        for _, p := range class.Properties() {
            w.Print("client.", g.toFieldName(p.Name()), ", err = ")
            if strings.HasSuffix(p.Type(), " operations") {
                w.Println("new", g.toGoType(p.Type()), "(client)")
            } else {
                w.Print("{},")
            }
            w.Println("if err != nil {")
            w.Indent()
            w.Println("return nil, err")
            w.Outdent()
            w.Println("}")
        }
        w.Println("return client, nil")

        w.Outdent()
        w.Println("}")
        w.Println()

        w.Println("func (c *", g.toTypeName(name), ") BaseUri() string {")
        w.Indent()
        w.Println("return c.baseUri")
        w.Outdent()
        w.Println("}")
        w.Println()
        w.Println("func (c *", g.toTypeName(name), ") WithBaseUri(value string) *", g.toTypeName(name), " {")
        w.Indent()
        w.Println("c.baseUri = value")
        w.Println("return c")
        w.Outdent()
        w.Println("}")
        w.Println()
        w.Println("func (c *", g.toTypeName(name), ") UserAgent() string {")
        w.Indent()
        w.Println("return c.userAgent")
        w.Outdent()
        w.Println("}")
        w.Println()
        w.Println("func (c *", g.toTypeName(name), ") WithUserAgent(value string) *", g.toTypeName(name), " {")
        w.Indent()
        w.Println("c.userAgent = value")
        w.Println("return c")
        w.Outdent()
        w.Println("}")
        w.Println()
        w.Println("func (c *", g.toTypeName(name), ") HttpClient() *http.Client {")
        w.Indent()
        w.Println("return c.httpClient")
        w.Outdent()
        w.Println("}")
        w.Println()
    }

    for _, p := range class.Properties() {
        if !p.HasStereoType(serializableStereoType) {
            g.generateGetter(class, p)
        }
    }

    for _, o := range class.Operations() {
        w.Print("func (c *", g.toTypeName(name), ") ", g.toOperationName(o.Name()), "(")
        for i, p := range o.Parameters() {
            if i > 0 {
                w.Print(", ")
            }
            w.Print(g.toFieldName(p.Name()), " ")
            if p.LowerValue() != p.UpperValue() {
                w.Print("[]")
            }
            if g.ptr(p.Type(), p.Required()) {
                w.Print("*")
            }
            w.Print(g.toGoType(p.Type()))
        }
        w.Print(") ")

        returnTypes := o.ReturnTypes()
        if len(returnTypes) > 0 {
            w.Print("(")
            for i, t := range returnTypes {
                if i > 0 {
                    w.Print(", ")
                }
                w.Print(g.toFieldName(t.Name()), " ")
                if t.LowerValue() != t.UpperValue() {
                    w.Print("[]")
                }
                if g.ptr(t.Type(), false) {
                    w.Print("*")
                }
                w.Print(g.toGoType(t.Type()))
            }
            w.Print(") ")
        }

        w.Println("{")
        w.Indent()
        w.Println("// ", o.extendedProperties["method"], " ", o.extendedProperties["path"])
        value, ok := o.extendedProperties["schema"]
        if !ok {
            log.Fatalln("Missing schema for operation", o.Name())
        }
        schema := value.(*openapi2.Operation)

        w.Println("var response *http.Response")
        w.Println("var request *http.Request")
        for _, p := range o.Parameters() {
            schema := p.extendedProperties["schema"].(*openapi2.Parameter)
            if schema.In == "body" {
                w.Println()
                w.Println("var buffer *bytes.Buffer")
                break
            }
        }
        w.Println()
        w.Println("client := c.client.httpClient")
        w.Println("credentials := c.client.credentials")
        w.Println()
        w.Println("baseUri := c.client.baseUri")
        w.Println("userAgent := c.client.userAgent")
        w.Println("if strings.HasSuffix(baseUri, \"/\") {")
        w.Indent()
        w.Println("baseUri = baseUri[:len(baseUri)-1]")
        w.Outdent()
        w.Println("}")
        w.Println("url := strings.Builder{}")
        w.Println("url.WriteString(baseUri)")
        path := o.extendedProperties["path"].(string)

        // find fist { in path
        i := strings.Index(path, "{")
        if i >= 0 {
            for i >= 0 {
                w.Println("url.WriteString(\"", path[:i], "\")")
                path = path[i:]
                j := strings.Index(path, "}")
                if j < 0 {
                    log.Fatalln("Invalid path", o.extendedProperties["path"])
                }
                paramName := path[1:j]
                w.Println("url.WriteString(", paramName, ")")
                path = path[j+1:]
                i = strings.Index(path, "{")
            }

        } else {
            // no parameters
            w.Println("url.WriteString(\"", path, "\")")
        }

        body := "nil"
        for _, p := range o.Parameters() {
            schema := p.extendedProperties["schema"].(*openapi2.Parameter)
            if schema.In == "body" {
                w.Println()
                w.Println("buffer = &bytes.Buffer{}")
                w.Println("err = json.NewEncoder(buffer).Encode(", g.toFieldName(p.Name()), ")")
                w.Println("if err != nil {")
                w.Indent()
                w.Println("return")
                w.Outdent()
                w.Println("}")

                body = "buffer"
            }
        }

        w.Println()
        w.Print("request, err = http.NewRequest(http.Method")
        method := o.extendedProperties["method"].(string)
        w.Print(strcase.ToCamel(strings.ToLower(method)))
        w.Print(", url.String(), ")
        w.Print(body)
        w.Println(")")
        w.Println("if err != nil {")
        w.Indent()
        w.Println("return")
        w.Outdent()
        w.Println("}")
        w.Println()

        w.Println("if userAgent == \"\" {")
        w.Indent()
        w.Println("userAgent = \"go-sophosutm\"")
        w.Outdent()
        w.Println("}")
        w.Println("request.Header.Set(\"User-Agent\", userAgent)")
        w.Println("request.Header.Set(\"Authorization\", credentials.GetAuthorizationHeaderValue())")

        if len(schema.Consumes) > 0 {
            w.Println("request.Header.Set(\"Content-Type\", \"", schema.Consumes[0], "\")")
        }
        if len(schema.Produces) > 0 {
            w.Println("request.Header.Set(\"Accept\", \"", schema.Produces[0], "\")")
        }

        query := false
        for _, p := range o.Parameters() {
            parameter, ok := p.extendedProperties["schema"].(*openapi2.Parameter)
            if !ok {
                continue
            }
            switch parameter.In {
            case "query":
                if !query {
                    w.Println("query := request.URL.Query()")
                }
                w.Println("query.Add(\"", parameter.Name, "\", ", g.toFieldName(p.Name()), ")")
                query = true
            case "header":
                if parameter.Required {
                    w.Println("request.Header.Set(\"", parameter.Name, "\", ", g.toFieldName(p.Name()), ")")
                } else {
                    w.Println("if ", g.toFieldName(p.Name()), " != nil {")
                    w.Indent()
                    w.Println("request.Header.Set(\"", parameter.Name, "\", *", g.toFieldName(p.Name()), ")")
                    w.Outdent()
                    w.Println("}")
                }

            case "path":
                // already handled
            case "body":
                // already handled
            default:
                log.Fatalln("Unsupported parameter location", parameter.In)
            }

            if query {
                w.Println("request.URL.RawQuery = query.Encode()")
            }

        }

        w.Println()
        w.Println("response, err = client.Do(request)")
        w.Println("if err != nil {")
        w.Indent()
        w.Println("return")
        w.Outdent()
        w.Println("}")
        w.Println("defer response.Body.Close()")
        w.Println()
        w.Println("switch response.StatusCode {")
        w.Indent()
        for c, r := range schema.Responses {
            statusCode, err := strconv.Atoi(c)
            if err != nil {
                log.Fatalln("Invalid status code", c)
            }
            w.Println("case ", statusCode, ":")
            w.Indent()
            if statusCode >= 200 && statusCode < 300 {
                if len(o.ReturnTypes()) > 0 {
                    for _, t := range o.ReturnTypes() {
                        if t.Name() == "err" {
                            continue
                        }
                        w.Print("var data ")
                        if t.LowerValue() != t.UpperValue() {
                            w.Print("[]")
                        }
                        if g.ptr(t.Type(), false) {
                            w.Print("*")
                        }
                        w.Print(g.toGoType(t.Type()))
                        w.Println()
                        w.Println("err = json.NewDecoder(response.Body).Decode(&data)")
                        w.Println("if err != nil {")
                        w.Indent()
                        w.Println("return")
                        w.Outdent()
                        w.Println("}")
                        w.Println("result = data")
                    }
                }

            } else {
                e, ok := HttpErrors[statusCode]
                if !ok {
                    w.Println("err = Err", g.toTypeName(r.Description))
                } else {
                    w.Println("err = Err", e)
                }
                g.generateDefaultReturn(w, returnTypes, true)

            }
            w.Outdent()
        }

        w.Println("default:")
        w.Indent()
        w.Println(`err = errors.New("Unexpected status code " + strconv.Itoa(response.StatusCode))`)
        g.generateDefaultReturn(w, returnTypes, true)
        w.Outdent()
        w.Println("}")
        w.Println("return")
        w.Outdent()
        w.Println("}")
        w.Println()
    }
}

func (g *GoGenerator) generateDefaultReturn(w *IndentedTextWriter, returnTypes []*ReturnType, hasError bool) {
    w.Print("return ")
    for i, t := range returnTypes {
        if i > 0 {
            w.Print(", ")
        }
        if t.Type() == "error" && hasError {
            w.Print("err")
        } else {
            w.Print(defaultValue(t.Type()))
        }
    }
    w.Println()
}

func (g *GoGenerator) ptr(typeName string, required bool) bool {
    ptr := false
    if !g.isInterfaceType(typeName) {
        if !required {
            ptr = true
        }
        if !g.isPrimitiveType(typeName) {
            ptr = true
        }
    }
    return ptr
}

func defaultValue(t string) string {
    return "nil"
}

func (g *GoGenerator) toFieldName(n string) string {
    return strcase.ToLowerCamel(n)
}

func (g *GoGenerator) toPublicFieldName(n string) string {
    return strcase.ToCamel(n)
}

func (g *GoGenerator) toTypeName(name string) string {
    return strcase.ToCamel(name)
}

func (g *GoGenerator) toGoType(s string) string {
    switch strings.ToLower(s) {
    case "string":
        return "string"
    case "integer":
        return "int"
    case "boolean":
        return "bool"
    case "interface{}":
        return "interface{}"
    case "error":
        return "error"
    default:
        return g.toTypeName(s)
    }
}

func (g *GoGenerator) Generate(m *Model) {
    g.reset("errors.go")
    g.generateErrors(m)
    classes := m.Classes()
    sort.Sort(classes)

    for _, class := range classes {
        g.reset(strcase.ToSnake(class.Name()) + ".go")
        g.generateClass(class)
    }

    for _, class := range classes {
        g.reset(strcase.ToSnake(class.Name()) + "_test.go")
        g.generateTests(class)
    }
}

func (g *GoGenerator) toOperationName(name string) string {
    return strcase.ToCamel(name)
}

func (g *GoGenerator) isPrimitiveType(s string) bool {
    switch strings.ToLower(s) {
    case "string":
        return true
    case "integer":
        return true
    case "boolean":
        return true
    default:
        return false
    }
}

func (g *GoGenerator) isInterfaceType(s string) bool {
    switch strings.ToLower(s) {
    case "interface{}":
        return true
    case "error":
        return true
    default:
        if strings.HasPrefix(s, "map[") {
            return true
        }
        if strings.HasPrefix(s, "[]") {
            return true
        }
        if strings.HasPrefix(s, "*") {
            return true
        }
        return false
    }
}

func (g *GoGenerator) generateGetter(class *Class, p *Property) {
    w := g.writer
    methodName := strcase.ToCamel(p.Name())
    methodName = strings.TrimSuffix(methodName, "Operations")

    methodName = removeCommonPrefix(methodName, strings.TrimSuffix(strcase.ToCamel(class.Name()), "Operations"))

    w.Print("func (c *", g.toTypeName(class.Name()), ") ", methodName, "() ")
    if p.LowerValue() != p.UpperValue() {
        w.Print("[]")
    }
    if g.ptr(p.Type(), p.Required()) {
        w.Print("*")
    }
    w.Println(g.toGoType(p.Type()), "{")
    w.Indent()
    w.Print("return c.")
    if p.HasStereoType(serializableStereoType) {
        w.Print(g.toPublicFieldName(p.Name()))
    } else {
        w.Print(g.toFieldName(p.Name()))
    }
    w.Println()
    w.Outdent()
    w.Println("}")
    w.Println()
}

func removeCommonPrefix(left string, right string) string {
    if left == right {
        return ""
    }
    if left == "" {
        return ""
    }

    if right == "" {
        return left
    }

    m := min(len(left), len(right))
    var i int
    for i = 0; i < m && left[i] == right[i]; i++ {

    }

    if i == 0 {
        return left
    }
    return left[i:]
}

func min(l int, r int) int {
    if l > r {
        return r
    }
    return l
}

var HttpErrors = map[int]string{
    400: "BadRequest",
    401: "Unauthorized",
    403: "Forbidden",
    404: "NotFound",
    405: "MethodNotAllowed",
    406: "NotAcceptable",
    408: "RequestTimeout",
    409: "Conflict",
    410: "Gone",
    411: "LengthRequired",
    412: "PreconditionFailed",
    413: "RequestEntityTooLarge",
    414: "RequestURITooLong",
    415: "UnsupportedMediaType",
    416: "RequestedRangeNotSatisfiable",
    417: "ExpectationFailed",
    418: "Teapot",
    422: "UnprocessableEntity",
    423: "Locked",
    424: "FailedDependency",
    425: "UnorderedCollection",
    426: "UpgradeRequired",
    428: "PreconditionRequired",
    429: "TooManyRequests",
    431: "RequestHeaderFieldsTooLarge",
    451: "UnavailableForLegalReasons",
    500: "InternalServerError",
    501: "NotImplemented",
    502: "BadGateway",
    503: "ServiceUnavailable",
    504: "GatewayTimeout",
    505: "HTTPVersionNotSupported",
    506: "VariantAlsoNegotiates",
    507: "InsufficientStorage",
    508: "LoopDetected",
    509: "BandwidthLimitExceeded",
    510: "NotExtended",
    511: "NetworkAuthenticationRequired",

}

func (g *GoGenerator) generateErrors(m *Model) {
    errors := make(map[int]bool)
    for _, class := range m.Classes() {
        for _, operation := range class.Operations() {
            value, ok := operation.extendedProperties["schema"]
            if !ok {
                log.Fatalln("Missing schema for operation", operation.Name())
            }
            schema := value.(*openapi2.Operation)
            for c := range schema.Responses {
                statusCode, err := strconv.Atoi(c)
                if err != nil {
                    log.Fatalln("Invalid status code", c)
                }
                if statusCode >= 200 && statusCode < 300 {
                    continue
                }
                errors[statusCode] = true
            }
        }
    }

    w := g.writer
    w.Println("package ", g.packageName)
    w.Println()
    w.Println("import (")
    w.Indent()
    w.Println(`"errors"`)
    w.Outdent()
    w.Println(")")
    w.Println()

    for statusCode := range errors {
        message, ok := HttpErrors[statusCode]
        if !ok {
            log.Fatalln("Missing error message for status code", statusCode)
        }
        w.Println("var Err", message, " = new", message, "Error()")
    }
    if len(errors) > 0 {
        g.writer.Println()
    }

    for statusCode := range errors {
        message, ok := HttpErrors[statusCode]
        if !ok {
            log.Fatalln("Missing error message for status code", statusCode)
        }

        description := strings.ReplaceAll(strcase.ToKebab(message), "-", " ")
        w.Println("func new", message, "Error() error { return errors.New(\"", description, "\") }")
    }
    if len(errors) > 0 {
        g.writer.Println()
    }
}

func (g *GoGenerator) reset(path string) error {
    err := os.MkdirAll(g.path, 0755)
    if err != nil {
        return err
    }

    f, err := os.Create(filepath.Join(g.path, path))
    if err != nil {
        return err
    }
    g.writer = NewIndentedTextWriter(f, "   ")
    return nil
}

func (g *GoGenerator) generateTests(class *Class) {
    operations := class.Operations()

    w := g.writer
    w.Println("package ", g.packageName)
    w.Println()
    if len(operations) == 0 {
        return
    }

    w.Println("import (")
    w.Indent()
    w.Println(`"os"`)
    w.Println(`"testing"`)
    w.Println(`"github.com/stretchr/testify/assert"`)
    w.Println(`"github.com/stretchr/testify/require"`)
    w.Println(`"github.com/stretchr/testify/suite"`)
    w.Outdent()
    w.Println(")")
    w.Println()

    w.Println("type ", strcase.ToCamel(class.Name()), "TestSuite struct {")
    w.Indent()
    w.Println("suite.Suite")
    w.Println()
    w.Println("client *", strcase.ToCamel(class.Name()))
    w.Outdent()
    w.Print("}")
    w.Println()
    w.Println()
    w.Println("func Test", strcase.ToCamel(class.Name()), "TestSuite(t *testing.T) {")
    w.Indent()
    w.Println("suite.Run(t, new(", strcase.ToCamel(class.Name()), "TestSuite))")
    w.Outdent()
    w.Println("}")
    w.Println()
    w.Println("func (s *", strcase.ToCamel(class.Name()), "TestSuite) SetupSuite() {")
    w.Indent()
    w.Println()
    w.Println(`username := os.Getenv("SOPHOS_UTM_USERNAME")`)
    w.Println(`password := os.Getenv("SOPHOS_UTM_PASSWORD")`)
    w.Println(`token := os.Getenv("SOPHOS_UTM_TOKEN")`)
    w.Println(`baseUri := os.Getenv("SOPHOS_UTM_HOST")`)
    w.Println(`insecureSkipVerify := os.Getenv("SOPHOS_UTM_INSECURE_SKIP_VERIFY") == "true"`)
    w.Println(`userAgent := os.Getenv("SOPHOS_UTM_USERAGENT")`)
    w.Println()
    w.Println(`var credentials Credentials`)
    w.Println(`if token != "" {`)
    w.Indent()
    w.Println(`credentials = &TokenCredentials{`)
    w.Indent()
    w.Println(`Token: token,`)
    w.Outdent()
    w.Println(`}`)
    w.Outdent()
    w.Println(`}`)
    w.Println()
    w.Println(`if username != "" && password != "" {`)
    w.Indent()
    w.Println(`credentials = &UsernamePasswordCredentials{`)
    w.Indent()
    w.Println(`Username: username,`)
    w.Println(`Password: password,`)
    w.Outdent()
    w.Println(`}`)
    w.Outdent()
    w.Println(`}`)
    w.Println(`require.NotNil(s.T(), credentials, "No credentials provided")`)
    w.Println()
    w.Println(`client, err := newClient(credentials, baseUri, userAgent, insecureSkipVerify)`)
    w.Println(`require.NoError(s.T(), err, "failed to create client")`)
    w.Println(`require.NotNil(s.T(), client, "No client created")`)
    w.Println()
    w.Println(`s.client = &`, strcase.ToCamel(class.Name()), `{client: client}`)
    w.Outdent()
    w.Println("}")
    w.Println()
    w.Println("func (s *", strcase.ToCamel(class.Name()), "TestSuite) SetupTest() {")
    w.Indent()
    w.Outdent()
    w.Println("}")
    w.Println()
    w.Println("func (s *", strcase.ToCamel(class.Name()), "TestSuite) TearDownTest() {")
    w.Indent()
    w.Outdent()
    w.Println("}")
    w.Println()
    w.Println("func (s *", strcase.ToCamel(class.Name()), "TestSuite) TearDownSuite() {")
    w.Indent()
    w.Outdent()
    w.Println("}")
    w.Println()

    g.buildOperationsTests(class)
}

func (g *GoGenerator) buildOperationsTests(class *Class) {
    w := g.writer
    operations := class.Operations()
    for _, operation := range operations {
        operationName := strings.ToLower(strcase.ToCamel(operation.Name()))
        if strings.HasPrefix(operationName, "list") {
            w.Println(
                `func (s *`, strcase.ToCamel(class.Name()), `TestSuite) Test`,
                strcase.ToCamel(operation.Name()),
                `() {`,
            )
            w.Indent()
            w.Println(`actual, err := s.client.`, strcase.ToCamel(operation.Name()), `()`)
            w.Println(`assert.NoError(s.T(), err, "failed to list")`)
            w.Println(`assert.NotEmpty(s.T(), actual, "empty list")`)
            w.Println(`s.T().Logf("Actual: %+v\n", actual)`)
            w.Outdent()
            w.Println(`}`)
            w.Println()

        } else if strings.HasPrefix(operationName, "get") {
            w.Println(
                `func (s *`, strcase.ToCamel(class.Name()), `TestSuite) Test`,
                strcase.ToCamel(operation.Name()),
                `() {`,
            )
            w.Indent()
            w.Println(`s.T().Skip("Not implemented")`)
            w.Outdent()
            w.Println(`}`)
            w.Println()

        } else if strings.HasPrefix(operationName, "create") {
            w.Println(
                `func (s *`, strcase.ToCamel(class.Name()), `TestSuite) Test`,
                strcase.ToCamel(operation.Name()),
                `() {`,
            )
            w.Indent()
            w.Println(`s.T().Skip("Not implemented")`)
            w.Outdent()
            w.Println(`}`)
            w.Println()
        } else if strings.HasPrefix(operationName, "update") {
            w.Println(
                `func (s *`, strcase.ToCamel(class.Name()), `TestSuite) Test`,
                strcase.ToCamel(operation.Name()),
                `() {`,
            )
            w.Indent()
            w.Println(`s.T().Skip("Not implemented")`)
            w.Outdent()
            w.Println(`}`)
            w.Println()

        } else if strings.HasPrefix(operationName, "delete") {
            w.Println(
                `func (s *`, strcase.ToCamel(class.Name()), `TestSuite) Test`,
                strcase.ToCamel(operation.Name()),
                `() {`,
            )
            w.Indent()
            w.Println(`s.T().Skip("Not implemented")`)
            w.Outdent()
            w.Println(`}`)
            w.Println()
        }
    }
}

func (g *GoGenerator) buildListOperationTest(class *Class, operation *Operation) {
    w := g.writer
    w.Println(`actual, err := s.client.`, strcase.ToCamel(operation.Name()), `()`)
    w.Println(`assert.NoError(s.T(), err, "failed to list")`)
    w.Println(`assert.NotEmpty(s.T(), actual, "empty list")`)
}

func NewGoGenerator(path, packageName string) *GoGenerator {
    return &GoGenerator{
        path:        path,
        packageName: packageName,
    }
}
