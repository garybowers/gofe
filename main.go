package main

import (
	"crypto/tls"
	"encoding/gob"
	"flag"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	store    *sessions.CookieStore
	googUser *userData
	port     *string
	backend  *string
)

const cookieName string = "gofe"

func init() {
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 15,
		HttpOnly: true,
	}

	gob.Register(User{})
}

type User struct {
	Email         string
	Username      string
	Forename      string
	Surname       string
	Authenticated bool
}

type userData struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
	Hd            string `json:"hd"`
}

func main() {
	log.Printf("GoFE v0.0.1")
	const (
		defaultPort      = ":9090"
		defaultPortUsage = "The port the reverse proxy should listen on, default is ':9090'"
		urlUsage         = "The backend url to proxy"
	)

	port = flag.String("port", defaultPort, defaultPortUsage)
	backend = flag.String("backend", "http://10.76.42.1:8000/", urlUsage)

	flag.Parse()
	log.Printf("Listening on port: %s", *port)
	log.Printf("Proxying backend: %s", *backend)

	http.HandleFunc("/", proxy)
	http.ListenAndServe(*port, nil)
}

func proxy(res http.ResponseWriter, req *http.Request) {
	url, _ := url.Parse(*backend)
	proxy := httputil.NewSingleHostReverseProxy(url)

	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	//req.Header.Del("Origin")
	req.Host = url.Host

	//req.URL.Path = mux.Vars(req)["rest"]

	log.Printf("%s %s %s", url.Scheme, req.RemoteAddr, req.URL)

	switch url.Scheme {
	case "https":
		proxy.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	}
	proxy.ServeHTTP(res, req)
}
