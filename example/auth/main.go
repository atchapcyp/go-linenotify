package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/utahta/go-linenotify"
)

// EDIT THIS
var (
	BaseURL      = "http://localhost"
	ClientID     = ""
	ClientSecret = ""
)

func Authorize(w http.ResponseWriter, req *http.Request) {
	c, err := linenotify.NewAuthorization(ClientID, BaseURL+"/callback")
	if err != nil {
		fmt.Fprintf(w, "error:%v", err)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "state", Value: c.State, Expires: time.Now().Add(60 * time.Second)})

	c.Redirect(w, req)
}

func Callback(w http.ResponseWriter, req *http.Request) {
	resp, err := linenotify.ParseAuthorization(req)
	if err != nil {
		fmt.Fprintf(w, "error:%v", err)
		return
	}

	state, err := req.Cookie("state")
	if err != nil {
		fmt.Fprintf(w, "error:%v", err)
		return
	}
	if resp.State != state.Value {
		fmt.Fprintf(w, "error:%v", err)
		return
	}

	c := linenotify.NewToken(resp.Code, BaseURL+"/callback", ClientID, ClientSecret)
	token, err := c.Get()
	if err != nil {
		fmt.Fprintf(w, "error:%v", err)
		return
	}

	fmt.Fprintf(w, "token:%v", token)
}

func main() {
	http.HandleFunc("/auth", Authorize)
	http.HandleFunc("/callback", Callback)

	http.ListenAndServe(":9090", nil)
}
