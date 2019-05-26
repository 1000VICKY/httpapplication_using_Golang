package main

import (
	"flag"
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
	"time"
)

var err error
var userRegex *regexp.Regexp
var matchedValue [][]string

type MyStringer interface {
	String() string
}

func getCommandOptionFromIndex(i int) string {
	// 指定したindexのコマンドオプションを取得する
	flag.Parse()
	var optionList []string
	optionList = flag.Args()
	if len(optionList) > 0 {
		return optionList[i]
	} else {
		return ""
	}
}

func main() {
	var s MyStringer = nil
	// 起動ポートを指定
	var portNumber string = ""
	portNumber = getCommandOptionFromIndex(0)
	if len(portNumber) == 0 {
		// ポート番号の指定が無い場合は 8080をポート番号とする
		portNumber = "8080"
	}
	fmt.Println(portNumber)
	// URLパターンマッチング
	userRegex, err = regexp.Compile("^/account/id/([0-9]+)/?$")
	// 静的ファイル配信の際に捜査するディレクトリ
	var publicRoot string = "C:/Users/senbiki/public"
	// テンプレートファイルのディレクトリ
	var templatePath string = "C:/Users/senbiki/public"
	s = userRegex
	fmt.Println("MyStringer")
	fmt.Println(s)
	fmt.Println("MyStringer")
	// マルチプレクサの生成
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

	// 外部ファイルに記述した関数でハンドラを登録
	server.HandleFunc("/account/user/", GetUserInfo)

	// GETパラメータに渡されたURLにHTTPリクエストを送信する
	server.HandleFunc("/GetReqestedURLDate", func(writer http.ResponseWriter, request *http.Request) {
		// 指定されたURLにアクセスする
		var query string = request.URL.RawQuery
		var queries []string = strings.Split(query, "&")
		var parseQueries map[string]string = make(map[string]string)
		for _, value := range queries {
			_value := strings.Split(value, "=")
			parseQueries[_value[0]] = _value[1]
		}
		for key, value := range parseQueries {
			fmt.Println(key + "=>" + value)
		}
		// 取得したい画像URLを設定する
		var imagePath *string = new(string)
		var ok *bool = new(bool)
		*imagePath, *ok = parseQueries["imagePath"]
		// HTTPリクエストを実行
		var _res *http.Response = nil
		var _error error = nil
		if *ok == true {
			_res, _error = http.Get(*imagePath)
			if _error != nil {
				fmt.Println(_error)
				return
			}
			var buffer []byte = make([]byte, 1024)
			var total []byte = make([]byte, 0)
			for {
				n, _error := _res.Body.Read(buffer)
				if _error != nil {
					break
				}
				// 1byteでも読み取りできた場合
				if n > 0 {
					total = append(total, buffer[0:n]...)
				}
			}
			//
			fmt.Fprint(writer, string(total))
		}
	})

	// クライアントへサードパーティーCookieを送信する
	server.HandleFunc("/FunctionToSendCookie/", func(writer http.ResponseWriter, request *http.Request) {
		// アクセス中のHost名を取得
		var _host string = request.Host
		if _host == "adserver:8080" {
			// アクセスしてきたブラウザへタイムスタンプを値とするクッキーを返却する
			var _time time.Time = time.Now()
			var _unixtime int64 = _time.Unix()
			var _unixtimeString string = strconv.Itoa(int(_unixtime))
			var _header http.Header = writer.Header()
			var _requestHeader http.Header = request.Header
			var _requestCookie string = _requestHeader.Get("Cookie")
			// GETパラメータを取得
			var _getQuery string = request.URL.RawQuery
			// アクセス時のHTTPメソッドを取得
			var _method string = request.Method
			_header.Add("Set-Cookie", "RequestMehod="+_method+";path=/; Expires=Mon 31 Dec 2025 23:59:59 GMT")
			_header.Add("Set-Cookie", "Third-Party-Cookie="+_unixtimeString+";path=/; Expires=Mon 31 Dec 2025 23:59:59 GMT")
			_header.Add("Set-Cookie", "RequestHost="+_host+";path=/; Expires=Mon 31 Dec 2025 23:59:59 GMT")
			if len(_requestCookie) > 0 {
				// 二回目以降のアクセスでブラウザからのクッキーを取得する
				fmt.Println("クライアントからCookieを取得しました。")
				fmt.Println("クライアントから送信されたCookie => " + _requestCookie)
				fmt.Println("GETパラメータ => " + _getQuery)
			} else {
				fmt.Println("Cookieのヘッダーがないので初回のアクセスです。")
				fmt.Println("GETパラメータ => " + _getQuery)
			}
		} else {
			// 404エラーを返却
			writer.WriteHeader(404)
			writer.Header().Set("Content-Type", "text/html")
			fmt.Fprintln(writer, "規定のホストではありません。")
		}
	})

	// サードパーティクッキーを検証する
	server.HandleFunc("/ValidateCookie/", func(writer http.ResponseWriter, request *http.Request) {
		var _host string = request.Host
		// 規定のホストからのみ実行を許可する
		if _host == "adserver:8080" {
			var _header http.Header = request.Header
			// 生のCookieを取得する
			var _cookie string = _header.Get("Cookie")
			var _cookieList []string = request.Header["Cookie"]
			fmt.Println(_cookieList)
			if len(_cookie) > 0 {
				fmt.Println("クライアントからCookieを取得しました。")
				fmt.Println("クライアントから送信されたCookie => " + _cookie)
			} else {
				fmt.Println("Cookieのヘッダーがないので初回のアクセスです。")
			}
			return
		} else {
			fmt.Fprintln(writer, "規定のホストではありません。")
			return
		}
	})

	// 以下より、Muliti Plexerを使ってルーティングを実装する
	// (1)各ユーザーのアカウントページ
	server.HandleFunc("/account/", func(writer http.ResponseWriter, request *http.Request) {
		// アクセスされたURL
		var requestedURL = request.URL.Path
		// http://sample.com/account/{mainCategory}/{subCategory}/
		pattern1, err := regexp.Compile("^/account/([0-9a-zA-Z_]+)/([0-9a-zA-Z_]+)/?$")
		if err != nil {
			fmt.Println("正規表現のコンパイルに失敗しました。")
			fmt.Println(err)
			return
		}
		// 正規表現にマッチした値を取得する
		var matchedValue [][]string
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
		} else if matchedValue = pattern1.FindAllStringSubmatch(requestedURL, -1); len(matchedValue) > 0 {
			// 指定した正規表現にマッチするURLへのアクセスの場合
			fmt.Println(matchedValue)
			var mainCategory = matchedValue[0][1]
			var subCategory = matchedValue[0][2]
			// システム的には上記2つのパラメータにマッチする情報を返却する
			fmt.Println("mainCategory => " + mainCategory)
			fmt.Println("subCategory => " + subCategory)
			return
		} else {
			var _header http.Header = request.Header
			_header.Set("Content-Type", "text/html; charset=UTF-8")
			_header.Set("Cookie", "Secret-Cookie=極秘Cookie")
			fmt.Println("仕様外のURLにアクセスしています。")
			return
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
		var _header http.Header = writer.Header()
		var _host string = request.Host
		var _method string = request.Method
		// GETパラメータを取得
		var query string = request.URL.RawQuery
		if requestedURL == "/" {
			// テンプレートオブジェクトの作成
			var _template *template.Template = nil
			_template = template.Must(template.ParseFiles(templatePath + "/index.tpl"))
			// テンプレートに渡すパラメータ
			var data map[string]string = make(map[string]string)
			data["ImagePath"] = "https://fukuoka.nasse.com/image/index/newimages/10989/b20170712_1.jpg/1024"
			data["requestedURL"] = "/"
			// ドメイン名をテンプレートを渡す
			data["query"] = query
			data["host"] = _host
			data["method"] = _method
			// 任意のHTTPレスポンスヘッダーを返却する
			_header.Add("Set-Cookie", "RequestedURL=/;path=/")
			_header.Set("Content-Type", "text/html;charset=UTF-8")
			_header.Set("Content-Error", "My-Originaol-Error")
			_template.Execute(writer, data)
			fmt.Println(data)
			return
		} else {
			// ルーティング設定した以外のURLにアクセスされた場合は静的ファイルを返却する
			var filesHandler http.Handler
			var fixHandler http.Handler
			filesHandler = http.FileServer(http.Dir(publicRoot))
			// 任意のヘッダーを返却する
			writer.Header().Set("Set-Cookie", "RequestedURL="+requestedURL+";path=/")
			writer.Header().Set("Content-Type", "text/html;charset=UTF-8")
			writer.Header().Set("Content-Host", _host)
			writer.Header().Set("Content-Method", _method)
			writer.Header().Set("Content-RequestedURL", requestedURL)
			fixHandler = http.StripPrefix("/", filesHandler)
			// ストリームに書き込み
			fixHandler.ServeHTTP(writer, request)
		}
	})

	err = http.ListenAndServeTLS("127.0.0.1:"+portNumber,
		"C:/Users/senbiki/go/bin/server.crt",
		"C:/Users/senbiki/go/bin/private.key",
		server,
	)
	if err != nil {
		fmt.Println(err)
	}
}
