package main

import (
    "encoding/json"
    "flag"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strconv"
    "strings"

    "github.com/getkin/kin-openapi/openapi2"
)

type RemoteSchemaReference struct {
    Description string `json:"description"`
    Name        string `json:"name"`
    Link        string `json:"link"`
}

const clientStereoType = "Client"
const dataOperationsStereoType = "Operations"
const dataSpecStereoType = "Data"
const serializableStereoType = "Serializable"

func main() {

    var url string
    flag.StringVar(
        &url,
        "url",
        "",
        "the hostnane of the the Sophos UTM appliance, e.g. https://firewall.example.net:4444/",
    )
    var path string
    flag.StringVar(&path, "path", "generated", "The path where files will be generated, such as 'go-sophosutm'")
    var packageName string
    flag.StringVar(&packageName, "package", "sophosutm", "The package name, such as 'sophosutm'")

    flag.Parse()

    if url == "" {
        log.Fatalln("usage: generator -url <hostname>")
    }
    // if url doesn't start with http then add prefix it with https
    if !strings.HasPrefix(url, "http") {
        url = "https://" + url
    }
    if strings.HasSuffix(url, "/") {
        url = url[:len(url)-1]
    }

    schemas, err := loadDefinitions(url)
    if err != nil {
        log.Fatal(err)
    }

    var schemaMap = make(map[*RemoteSchemaReference]*openapi2.T)
    schemaMap, err = buildSchemaMap(schemas, url)
    if err != nil {
        log.Fatal(err)
    }

    mb := NewModelBuilder()
    var model = mb.Build(schemaMap)

    g := NewGoGenerator(path, packageName)

    g.Generate(model)

}

func returnsArray(response *openapi2.Response) bool {
    if response.Schema != nil {
        if response.Schema.Value != nil {
            if response.Schema.Value.Type == "array" {
                return true

            }
        }
    }
    return false
}

func findSuccessfulResponse(operation *openapi2.Operation) *openapi2.Response {
    for code, response := range operation.Responses {
        c, err := strconv.Atoi(code)
        if err != nil {
            continue
        }
        if c >= 200 && c < 300 {
            return response
        }
    }
    return nil
}

func loadDefinitions(url string) ([]*RemoteSchemaReference, error) {
    b, err := downloadBytes(url + "/api/definitions/")
    if err != nil {
        log.Fatal(err)
    }

    var schemas []*RemoteSchemaReference
    err = json.Unmarshal(b, &schemas)
    if err != nil {
        log.Fatal(err)
    }
    return schemas, err
}

func buildSchemaMap(schemas []*RemoteSchemaReference, url string) (
    result map[*RemoteSchemaReference]*openapi2.T,
    err error,
) {
    var schemaMap = make(map[*RemoteSchemaReference]*openapi2.T)
    for _, schema := range schemas {
        b2, err2 := downloadBytes(url + schema.Link)
        if err2 != nil {
            log.Fatal(err)
        }

        os.MkdirAll("schemas", 0755)
        err2 = ioutil.WriteFile("schemas/"+schema.Name+".json", b2, 0644)
        if err2 != nil {
            log.Fatal(err)
        }

        var s openapi2.T
        err2 = json.Unmarshal(b2, &s)
        if err2 != nil {
            log.Fatal(err)
        }
        schemaMap[schema] = &s
    }
    return schemaMap, err
}

func ShouldSkip(document *RemoteSchemaReference) bool {
    return !strings.HasPrefix(document.Name, "aaa")
}

func downloadBytes(url string) ([]byte, error) {
    client := http.Client{
        CheckRedirect: func(r *http.Request, via []*http.Request) error {
            r.URL.Opaque = r.URL.Path
            return nil
        },
    }

    resp, err := client.Get(url)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    b, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
    return b, err
}
