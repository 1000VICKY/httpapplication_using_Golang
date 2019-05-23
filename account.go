package main

import (
	"fmt"
	"html/template"
	"net/http"
	"phpi/myreflect"
)

func AccountIndex(writer http.ResponseWriter, request *http.Request) {
	var myTemplate *template.Template
	myTemplate = template.New("account/index.tpl")
	myTemplate = template.Must(myTemplate.Parse("html"))

	var param map[string]string = nil
	var header http.Header
	header = writer.Header()
	header.Set("Content-Type", "text/html")
	myTemplate.Execute(writer, param)
}

func GetUserInfo(writer http.ResponseWriter, request *http.Request) {
	var param map[string]string = make(map[string]string)
	t := template.Must(template.ParseFiles("account/user.tpl"))
	param["userName"] = "ユーザー名"
	fmt.Println(myreflect.
		GetObjectMethods(writer))
	writer.Header().Set("Content-Type", "text/html;charset=UTF-8")
	t.Execute(writer, param)
}
