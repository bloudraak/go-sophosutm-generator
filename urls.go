package main

import "strings"

func join(paths ...string) string {
    if paths == nil {
        return ""
    }
    sb := strings.Builder{}
    for _, path := range paths {
        if path == "" {
            continue
        }
        if sb.Len() > 0 {
            if sb.String()[sb.Len()-1] != '/' && path[0] != '/' {
                sb.WriteByte('/')
            }
        }
        if path != "/" {
            sb.WriteString(path)
        }
    }
    return sb.String()
}
