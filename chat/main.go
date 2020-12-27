package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/TTTTk84/go_webapplication/trace"
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
			template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	// 出力
	err := t.temp1.Execute(w, r)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	var addr = flag.String("addr", ":8080", "アドレス")
	flag.Parse()
	r := newRoom()
	r.tracer = trace.New(os.Stdout)
	http.Handle("/", &templateHandler{filename: "chat.html"})
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
