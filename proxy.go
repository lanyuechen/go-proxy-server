package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

type Origin struct {
	Scheme string `json: "scheme"`
	Host   string `json: "host"`
	Path   string `json: "path"`
}

type ProxyConfig struct {
	Test string  `json: "test"`
	To   *Origin `json: "to"`
}

func midProxy(configDir string) func(http.ResponseWriter, *http.Request) {
	proxyConfig := []*ProxyConfig{}
	Load(configDir, &proxyConfig)

	for _, c := range proxyConfig {
		log.Printf("[proxy] %s >> %s://%s%s", c.Test, c.To.Scheme, c.To.Host, c.To.Path)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		for _, c := range proxyConfig {
			reg, _ := regexp.Compile(c.Test)
			if reg.MatchString(r.URL.Path) {
				Proxy(url.URL{
					Scheme: c.To.Scheme,
					Host:   c.To.Host,
					Path:   r.URL.Path,
				}, w, r)
			}
		}
	}
}

func Proxy(target url.URL, w http.ResponseWriter, r *http.Request) {
	// we need to buffer the body if we want to read it here and send it
  // in the request. 
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// you can reassign the body if you need to parse it as multipart
	// r.Body = ioutil.NopCloser(bytes.NewReader(body))

	// create a new url from the raw RequestURI sent by the client
	url := fmt.Sprintf("%s://%s%s", target.Scheme, target.Host, r.RequestURI)

	proxyReq, _ := http.NewRequest(r.Method, url, bytes.NewReader(body))

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	w.Write(body)
}

func Load(filename string, v interface{}) {
	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, v)
	if err != nil {
		return
	}
}
