package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"

	gomniauthcommon "github.com/stretchr/gomniauth/common"
)

type ChatUser interface {
	UniqueID() string
	AvatarURL() string
}

type chatUser struct {
	gomniauthcommon.User
	uniqueID string
}

func (u chatUser) UniqueID() string {
	return u.uniqueID
}

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	if cookie, err := r.Cookie("auth"); err == http.ErrNoCookie || cookie.Value == "" {
		// 未認証
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		// 何らかの別のエラーが発生
		panic(err.Error())
	} else {
		// 成功。ラップされたハンドラを呼び出します
		h.next.ServeHTTP(w, r)
	}
}

// 認証処理を行うhttphandlerを返す
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// ログイン処理待ち受け
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("認証プロバイダーの取得に失敗しました:", provider, "-", err)
		}
		loginUrl, err := provider.GetBeginAuthURL(nil,nil) // 認可サーバのurlを取得
		if err != nil {
			log.Fatalln("GetBeginAuthURLの呼び出し中にエラーが発生しました:", provider, "-", err)
		}
		w.Header().Set("Location",loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect) //307

		log.Println("TODO: ログイン処理", provider)
	case "callback":
		provider, err := gomniauth.Provider(provider)
		if err != nil{
			log.Fatalln("認証プロバイダーの取得に失敗しました:", provider, "-", err)
		}

		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			log.Fatalln("認証を完了できませんでした。", provider, "-", err)
		}

		user, err := provider.GetUser(creds)
		if err != nil {
			log.Fatalln("ユーザの取得に失敗しました", provider, "-", err)
		}
		chatUser := &chatUser{User: user}

		// userIDにmd5でハッシュ化したemailアドレスを入れる
		m := md5.New()
		io.WriteString(m, strings.ToLower(user.Email()))
		chatUser.uniqueID = fmt.Sprintf("%x", m.Sum(nil))

		avatarURL, err := avatars.GetAvatarURL(chatUser)
		if err != nil {
			log.Fatalln("GetAvatarURLに失敗しました", "-", err)
		}
		/* 認証が成功して、/chatにリダイレクトするときにauthというcookieに、
		 {name: ユーザ名, avatar_url: url, email: email} という形で保存する */
		authCookieValue := objx.New(map[string]interface{}{
			"userid": chatUser.uniqueID,
			"name": user.Name(),
			"avatar_url": avatarURL,
			//"email": user.Email(),
		}).MustBase64()
		http.SetCookie(w, &http.Cookie{
			Name: "auth",
			Value: authCookieValue,
			Path: "/"})

		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)

	default:
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w,"アクション%sには非対応です", action)
	}

}
