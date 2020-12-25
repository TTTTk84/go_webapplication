package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
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
	err := t.temp1.Execute(w, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := newRoom()
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)

	go r.run()

	// webサーバを開始する
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
