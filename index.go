package main

import (
	"fmt"
	"net/http"
	"os"
	"phpi/myreflect"
	"strconv"
)

func main() {
	// 静的ファイル配信際に捜査するディレクトリ
	var publicRoot string = "C:\\Users\\senbiki\\public"

	// マルチプレクサの生成
	// var server *http.ServeMux = http.NewServeMux()
	var server *http.ServeMux = new(http.ServeMux)

	// マルチプレクサオブジェクトのメソッド検証
	var methodList []string
	var err error
	methodList, err = myreflect.GetObjectMethods(server)
	if err != nil {
		fmt.Println(err)
		fmt.Println("*http.ServeMux型オブジェクトのメソッド一覧の取得に失敗しました。")
		os.Exit(255)
	}
	for key, value := range methodList {
		fmt.Println(strconv.Itoa(key) + " => " + value)
	}
	// publicなメソッドは以下一覧
	// 0 =>
	// 1 =>
	// 2 =>
	// 3 =>
	// 4 => Handle
	// 5 => HandleFunc
	// 6 => Handler
	// 7 => ServeHTTP

	// 以下より、Muliti Plexerを使ってルーティングを実装する
	// (1)各ユーザーのアカウントページ
	server.HandleFunc("/account/", func(writer http.ResponseWriter, request *http.Request) {
		// アクセスされたURL
		var reqeustedURL = request.URL.Path
		fmt.Println(reqeustedURL + "へアクセスされています。")

		// 以下よりクライアント側へ任意のレスポンスヘッダーを返却する
		var header http.Header
		header = writer.Header()
		header.Set("Content-RequestedURL", reqeustedURL)
		header.Set("Content-Route", "/account/")
		fmt.Fprint(writer, "Welcome to the page named /account/.")
	})

	// (2)会員登録ページを想定する
	server.HandleFunc("/register/", func(writer http.ResponseWriter, request *http.Request) {
		var requestedURL string = request.URL.Path
		var header http.Header = writer.Header()
		fmt.Println(requestedURL + "へアクセスされています。")
		// この時、仮にhttps://～/register/～のURLにアクセスされた場合に、静的ファイルを返却する場合
		// ただし、http://～/regiseter/へのアクセスはシステムで制御する
		if requestedURL == "/register/" {
			header.Set("Content-RequestedURL", requestedURL)
			header.Set("Content-Route", "/register/")
			fmt.Fprint(writer, "Welcome to the page named /register/.")
		} else {
			fmt.Println("/register/～")
			fmt.Println(request.URL.RawQuery)
			fmt.Println(request.URL.RawPath)
			// 静的ファイルを返却する場合
			var filesHandler http.Handler
			var fixHandler http.Handler
			// 指定したディレクトリをpublicルートとする
			filesHandler = http.FileServer(http.Dir(publicRoot))
			fixHandler = http.StripPrefix("/", filesHandler)
			// ストリームに書き込み
			fixHandler.ServeHTTP(writer, request)
		}
	})

	// (3)URLルートへのアクセス
	server.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		var requestedURL string = request.URL.Path
		var header http.Header = writer.Header()
		if requestedURL == "/" {
			fmt.Println(myreflect.GetObjectMethods(request))
			// 任意のHTTPレスポンスヘッダーを返却する
			writer.Header().Set("Content-Original", "My-Original-Header")
			writer.Header().Set("Content-Error", "1024")
			writer.Header().Set("A", "a")
			header.Set("A", "aaa")
			fmt.Fprintf(writer, "Hello world")
		} else {
			fmt.Println("/～")
			// ルーティング設定した以外のURLにアクセスされた場合は静的ファイルを返却する
			var filesHandler http.Handler
			var fixHandler http.Handler
			filesHandler = http.FileServer(http.Dir(publicRoot))
			// 任意のヘッダーを返却する
			header.Set("Content-Route", "/")
			header.Set("Content-RequestedURL", requestedURL)
			fixHandler = http.StripPrefix("/", filesHandler)
			// ストリームに書き込み
			fixHandler.ServeHTTP(writer, request)
		}
	})

	err = http.ListenAndServeTLS(":8080",
		"C:\\Users\\senbiki\\go\\bin\\server.crt",
		"C:\\Users\\senbiki\\go\\bin\\private.key",
		server,
	)
	if err != nil {
		fmt.Println(err)
	}
}
