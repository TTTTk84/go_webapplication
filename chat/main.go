package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"chat/trace"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

//template用構造体
type templateHandler struct {
	once sync.Once 						// 関数を一度だけ呼び出したいときに使う
	filename string						//
	temp1 *template.Template	//
}




// chat.html出力
func (t *templateHandler) ServeHTTP (w http.ResponseWriter, r *http.Request) {
	// 一度だけ実行
	t.once.Do(func() {
		// template.Must(template.New("name").Parse("text"))
		t.temp1 =
			template.Must(template.ParseFiles(filepath.Join("chat/templates", t.filename)))
	})
	data := map[string]interface{} {
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	// 出力
	err := t.temp1.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var addr = flag.String("addr", ":8080", "アドレス")
	flag.Parse()
	client_id := os.Getenv("GOOGLE_CLIENT_ID")
	client_secret := os.Getenv("GOOGLE_CLIENT_SECRET")

	// Gomniauthのセットアップ
	gomniauth.SetSecurityKey("sec6696123")
	gomniauth.WithProviders(
		facebook.New(client_id, client_secret, "http://localhost:8080/auth/callback/facebook"),
		github.New(client_id, client_secret,"http://localhost:8080/auth/callback/github"),
		google.New(client_id, client_secret,"http://localhost:8080/auth/callback/google"),
	)

	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	//http.HandleはServe.HTTPを実装していないとだめだが、http.HandleFuncはなくてもよい。
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)

	go r.run()

	// webサーバを開始する
	log.Println("webサーバーを開始します。ポート: ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}

/*
http.Handle(pattern string, handler http.Handler)は、

type Handler interface {
  ServeHTTP(ResponseWriter, *Request)
}

で、handlerは、ServeHTTPをメソッドに定義して入ればなんでもよい.
*/
