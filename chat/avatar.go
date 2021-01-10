package main

import (
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// ErrNoAvatarはAvatarインスタンスがアバターのURLを返すことができない場合に発生するエラー
var ErrNoAvatarURL = errors.New("chat: アバターのURLを取得できません")

// Avatarはユーザのプロフィール画像を表す型
// gravatarを使う場合と使わない場合で分けることができる
type Avatar interface {
	// GetAvatarURLは指定されたクライアントのアバターのURLを返す
	// 問題が発生した場合にはエラーを返す。URLが取得できなかったら、ErrNoAvatarURLを返す.
	GetAvatarURL(ChatUser) (string, error)
}

// まとめる
type TryAvatars []Avatar
func (a TryAvatars) GetAvatarURL(u ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(u); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}


// ローカルにある画像を使う
type FileSystemAvatar struct{}
var UseFileSystemAvatar FileSystemAvatar
func (FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	files, err := ioutil.ReadDir("avatars")
	if err != nil {
		return "", ErrNoAvatarURL
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fname := file.Name()
		if u.UniqueID() == strings.TrimSuffix(fname, filepath.Ext(fname)) {
			return "/avatars/" + fname, nil
		}
	}
	return "", ErrNoAvatarURL
}


// プロバイダごとの画像を使う
type AuthAvatar struct{}
var UseAuthAvatar AuthAvatar
func (AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()
	if len(url) == 0 {
		return "", ErrNoAvatarURL
	}
	return u.AvatarURL(), nil
}



// gravatarにupされている画像を使う
type GravatarAvatar struct{}
var UseGravatar GravatarAvatar
func (GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}
