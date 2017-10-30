package main

import (
        "testing"
        "reflect"
        "os"
        "bufio"
        "net/http"
        "io/ioutil"
        "fmt"
        "path/filepath"
)

const addr = "localhost:8081"

type tHandler struct {
        path string
}

func (th *tHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html")
        w.WriteHeader(http.StatusOK)
        data, err := ioutil.ReadFile(th.path)
        if err != nil {
                panic(err)
        }
        w.Header().Set("Content-Length", fmt.Sprint(len(data)))
        fmt.Fprint(w, string(data))
}


func startServer(paths []string) *http.Server {
        mux := http.NewServeMux()
        for _, path := range paths {
                h1 := &tHandler{path}
                mux.Handle("/" + filepath.Base(path), h1)
        }
        server := http.Server{Handler: mux, Addr: addr}
        go func() { server.ListenAndServe() }()
        return &server

}

func TestSimpleCyclic(t *testing.T) {

}




func TestBrokenLink(t *testing.T) {
        srv := startServer([]string{"tests/testBrokenLink/1.html"})
        expUrls := make(map[string]bool)
        expUrls["http://localhost:8081/1.html"] = true
        if !testParams(t, 1, false, "http://localhost:8081/1.html", &expUrls) {
        }
        defer srv.Shutdown(nil)
}

func TestExternalLink(t *testing.T) {

}


func testParams(t *testing.T, limit int, noLimit bool, url string, expected *map[string]bool) bool {
        out := make(chan []string, 100)
        go func() {
                <-out
        }()
        resUrls := startCrawling(out, limit, noLimit, url)
        return reflect.DeepEqual(*resUrls, *expected)
}

func testTomParams(t *testing.T, limit int, noLimit bool) bool {
        url := "http://tomblomfield.com/"
        file, err := os.Open("tests/testTom.txt")
        if(err != nil) {
                t.Fatal("Cannot open test file")
        }
        defer file.Close()
        testUrls := make(map[string]bool)
        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
                testUrls[scanner.Text()] = true
        }
        if err = scanner.Err(); err != nil {
                t.Fatal("Error while reading file")
        }
        return testParams(t, limit, noLimit, url, &testUrls)
}

func TestTomDefault(t *testing.T) {
        if !testTomParams(t, 15, false) {
                t.Error("Result URLs are not equal to test URLs")
        }
}
func TestTomLimitOne(t *testing.T) {
        if !testTomParams(t, 1, false) {
                t.Error("Result URLs are not equal to test URLs")
        }
}
func TestTomNoLimit(t *testing.T) {
        if !testTomParams(t, 1, true) {
                t.Error("Result URLs are not equal to test URLs")
        }
}
