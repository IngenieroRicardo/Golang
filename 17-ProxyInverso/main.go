package main

import (
        "bytes"
        "fmt"
        "io"
        "net/http"
        "net/http/httputil"
        "net/url"
        "strings"
)

var tunel = "http://192.100.1.2"

func handleRequest(w http.ResponseWriter, r *http.Request) {
        target, _ := url.Parse(tunel) 
        proxy := httputil.NewSingleHostReverseProxy(target)

        fmt.Println("CLIENTE:", r.RemoteAddr)

        proxy.ModifyResponse = func(res *http.Response) error {
                body, err := io.ReadAll(res.Body)
                if err != nil {
                        return err
                }
                err = res.Body.Close()
                if err != nil {
                        return err
                }
                sbody := string(body)

                if strings.Contains(sbody, tunel) {
                        sbody = strings.ReplaceAll(sbody, tunel, res.Request.Header.Get("Origin"))
                }
                res.Body = io.NopCloser(bytes.NewReader([]byte(sbody)))
                res.Header.Del("Content-Length")
                res.Header.Set("Transfer-Encoding", "chunked")
                return nil
        }
        r.Header.Add("X-Forwarded-For", r.RemoteAddr)
        proxy.ServeHTTP(w, r)
}

func main() {
        http.HandleFunc("/", handleRequest)
        fmt.Println("Server listening on port 5004")
        http.ListenAndServe(":5004", nil)
}

