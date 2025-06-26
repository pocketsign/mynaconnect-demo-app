package main

import (
	"net/http"
)

const (
	SessionCookieName = "myna_connect_session"
)

// Session はセッション情報を保持する構造体
type Session struct {
	ID string
}

// SetSession はセッション情報を保存する
func SetSession(w http.ResponseWriter, sessionData *Session) error {
	cookie := &http.Cookie{
		Name:     SessionCookieName,
		Value:    sessionData.ID,
		Path:     "/",
		MaxAge:   3600, // 1時間
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
	return nil
}

// GetSession はセッション情報を取得する
func GetSession(r *http.Request) (*Session, error) {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			return nil, nil // セッションが存在しない
		}
		return nil, err
	}

	if cookie.Value == "" {
		return nil, nil
	}

	return &Session{
		ID: cookie.Value,
	}, nil
}
