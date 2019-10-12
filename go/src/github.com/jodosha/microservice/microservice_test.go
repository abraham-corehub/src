package microservice

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	port      = "8080"
	httpsPort = "8081"
)

//
// Unit Tests
//

func TestHandler(t *testing.T) {
	expected := []byte("Hello World")

	req, err := http.NewRequest("GET", buildUrl("/"), nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	handler(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("Response code was %v; want 200", res.Code)
	}

	if bytes.Compare(expected, res.Body.Bytes()) != 0 {
		t.Errorf("Response body was '%v'; want '%v'", expected, res.Body)
	}
}

//
// Integration Tests
//

func TestHTTPServer(t *testing.T) {
	srv := NewServer(port)
	go srv.ListenAndServe()
	defer srv.Close()
	time.Sleep(100 * time.Millisecond)

	res, err := http.Get(buildUrl("/"))

	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Response code was %v; want 200", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	expected := []byte("Hello World")

	if bytes.Compare(expected, body) != 0 {
		t.Errorf("Response body was '%v'; want '%v'", expected, body)
	}
}

func TestHTTPSServer(t *testing.T) {
	srv := NewServer(httpsPort)
	go srv.ListenAndServeTLS("cert.pem", "key.pem")
	defer srv.Close()
	time.Sleep(100 * time.Millisecond)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	res, err := client.Get(buildSecureUrl("/"))

	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Response code was %v; want 200", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	expected := []byte("Hello World")

	if bytes.Compare(expected, body) != 0 {
		t.Errorf("Response body was '%v'; want '%v'", expected, body)
	}
}

//
// Private Helpers
//

func buildUrl(path string) string {
	return urlFor("http", port, path)
}

func buildSecureUrl(path string) string {
	return urlFor("https", httpsPort, path)
}

func urlFor(scheme string, serverPort string, path string) string {
	return scheme + "://localhost:" + serverPort + path
}
