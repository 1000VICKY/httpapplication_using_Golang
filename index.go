package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"phpi/myreflect"
	"regexp"
	"strconv"
	"strings"
)

var err error
var userRegex *regexp.Regexp
var matchedValue [][]string

func main() {
	// URLパターンマッチング
	userRegex, err = regexp.Compile("^/account/id/([0-9]+)/?$")
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
		// fmt.Fprint(writer, "Welcome to the page named /account/.")

		pattern1, err := regexp.Compile("^/account/([0-9a-zA-Z_]+)/([0-9a-zA-Z_]+)/?$")
		if err != nil {

		}
		if reqeustedURL == "/account/" {
			var getParamter string = request.URL.RawQuery
			var replacedParameter string = strings.Replace(getParamter, "imagePath=", "", -1)
			res, _ := http.Get(replacedParameter)
			binary, _ := ioutil.ReadAll(res.Body)
			fmt.Fprint(writer, string(binary))
			// var getParamter string = request.URL.RawQuery
			// var reg *regexp.Regexp
			// reg, _ = regexp.Compile("imagePath=(.*)$")
			// matchedValue = reg.FindAllStringSubmatch(getParamter, -1)
			// if len(matchedValue) > 0 {
			// 	res, _ := http.Get(matchedValue[0][1])
			// 	binary, _ := ioutil.ReadAll(res.Body)
			// 	// fmt.Println(string(binary))
			// 	header.Set("Content-Type", "image/jpeg")
			// 	fmt.Fprint(writer, string(binary))
			// }

		} else if matchedValue = pattern1.FindAllStringSubmatch(reqeustedURL, -1); len(matchedValue) > 0 {
			/** /acccount/Acategory/Bgenre のようなURLへのアクセス時 */
			var mainCategory string
			var subCategory string
			mainCategory = matchedValue[0][1]
			subCategory = matchedValue[0][2]
			fmt.Fprintln(writer, "以下が閲覧中の情報です<br>")
			fmt.Fprint(writer, mainCategory)
			fmt.Fprint(writer, subCategory)
		} else {
			var userID int
			var err error
			matchedValue = userRegex.FindAllStringSubmatch(reqeustedURL, -1)
			if len(matchedValue) > 0 {
				userID, err = strconv.Atoi(matchedValue[0][1])
				if err == nil {
					fmt.Fprint(writer, "閲覧中のユーザーIDは"+strconv.Itoa(userID)+"です")
				} else {
					header.Set("Content-Type", "text/html;charset=UTF-8")
					fmt.Println("ユーザIDの取得に失敗しました。")
				}
			}
		}
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
