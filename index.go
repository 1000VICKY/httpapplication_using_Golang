package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"phpi/myreflect"
	"reflect"
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

	// 外部ファイルに記述いした関数でハンドラを登録
	server.HandleFunc("/account/user/", GetUserInfo)

	// 以下より、Muliti Plexerを使ってルーティングを実装する
	// (1)各ユーザーのアカウントページ
	server.HandleFunc("/account/", func(writer http.ResponseWriter, request *http.Request) {
		// アクセスされたURL
		// var reqeustedURL = request.URL.Path
		// fmt.Println(reqeustedURL + "へアクセスされています。")

		// fmt.Fprint(writer, "Welcome to the page named /account/.")

		// pattern1, err := regexp.Compile("^/account/([0-9a-zA-Z_]+)/([0-9a-zA-Z_]+)/?$")
		// if err != nil {

		// }
		var requestedURL string = request.URL.Path
		if requestedURL == "/account/" {
			fmt.Println(request.Header.Get("If-Modified-Since"))
			fmt.Println(request.Header.Get("User-Agent"))
			print("===")
			var getParamter string = request.URL.RawQuery
			var replacedParameter string = strings.Replace(getParamter, "imagePath=", "", -1)
			res, _ := http.Get(replacedParameter)
			binary, _ := ioutil.ReadAll(res.Body)
			fmt.Println(myreflect.GetObjectMethods(writer))
			fmt.Println(reflect.TypeOf(writer))
			var lastModified string = res.Header.Get("Last-Modified")
			fmt.Println(lastModified)
			// レスポンスヘッダーを返却
			// 以下よりクライアント側へ任意のレスポンスヘッダーを返却する
			writer.Header().Set("Content-RequestedURL", requestedURL)
			writer.Header().Set("Content-Route", "/account/")
			writer.Header().Set("Content-Original", "unique-header")
			writer.Header().Set("Cache-Control", "max-age=640480")
			writer.Header().Set("Last-Modified", lastModified)
			writer.Header().Set("ETag", "----------")
			writer.Header().Set("Cache-Control", strconv.Itoa(60*60*24*14))
			writer.Header().Set("Pragma", "cache")
			writer.Header().Set("Content-Type", "image/jpeg")
			fmt.Fprint(writer, string(binary))
			return
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
			// return

			// } else if matchedValue = pattern1.FindAllStringSubmatch(reqeustedURL, -1); len(matchedValue) > 0 {
			// 	/** /acccount/Acategory/Bgenre のようなURLへのアクセス時 */
			// 	var mainCategory string
			// 	var subCategory string
			// 	mainCategory = matchedValue[0][1]
			// 	subCategory = matchedValue[0][2]
			// 	fmt.Fprintln(writer, "以下が閲覧中の情報です<br>")
			// 	fmt.Fprint(writer, mainCategory)
			// 	fmt.Fprint(writer, subCategory)
		} else {
			// var userID int
			// var err error
			// if len(matchedValue) > 0 {
			// 	userID, err = strconv.Atoi(matchedValue[0][1])
			// 	if err == nil {
			// 		fmt.Fprint(writer, "閲覧中のユーザーIDは"+strconv.Itoa(userID)+"です")
			// 	} else {
			// 		fmt.Println("ユーザIDの取得に失敗しました。")
			// 	}
			// }
			// return
		}
		fmt.Println("return 後")
		fmt.Fprint(writer, "return した後に、レスポンスを返却")
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
		tpl := template.Must(template.ParseFiles("./index.tpl"))
		if requestedURL == "/" {
			fmt.Println("/へアクセス")
			vParam := map[string]string{
				"ImagePath": "https://fukuoka.nasse.com/image/index/newimages/10989/b20170712_1.jpg/1024",
			}

			// 任意のHTTPレスポンスヘッダーを返却する
			// writer.Header().Set("Content-Original", "My-Original-Header")
			// writer.Header().Set("Content-Error", "1024")
			// writer.Header().Set("A", "a")
			// header.Set("A", "aaa")
			header.Set("Content-Type", "text/html;charset=UTF-8")
			tpl.Execute(writer, vParam)
			fmt.Println(vParam)
			fmt.Println(myreflect.GetObjectMethods(request))
			// fmt.Fprintf(writer, "Hello world")
			return
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
