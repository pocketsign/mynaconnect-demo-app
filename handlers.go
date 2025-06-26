package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	mynaconnectv1 "buf.build/gen/go/pocketsign/apis/protocolbuffers/go/pocketsign/mynaconnect/v1"
	"connectrpc.com/connect"
)

// handleCSS は共通CSSファイルを配信する
func handleCSS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	http.ServeFile(w, r, "templates/common.css")
}

// handleIndex はトップページを表示する
func handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/index.html")
}

// handleStart はマイナコネクトのセッションを開始する
func handleStart(w http.ResponseWriter, r *http.Request) {
	// セッション作成リクエスト
	resp, err := client.CreateSelfPersonalDataRequestSession(context.Background(), connect.NewRequest(&mynaconnectv1.CreateSelfPersonalDataRequestSessionRequest{
		ServiceId: serviceID,
		// 照会条件として今日の日付を設定
		InquiryCond: &mynaconnectv1.SelfPersonalDataInquiryCondition{
			Condition: &mynaconnectv1.SelfPersonalDataInquiryCondition_Date{
				Date: &mynaconnectv1.SelfPersonalDataInquiryCondition_DateCondition{
					Date: time.Now().Format("20060102"),
				},
			},
		},
		RedirectUri: callbackBaseURL + "/callback",
	}))
	if err != nil {
		log.Printf("セッション作成エラー: %v", err)
		renderError(w, "セッション作成エラー", "マイナコネクトとの通信に失敗しました。", err.Error())
		return
	}

	// セッション情報を保存
	if err := SetSession(w, &Session{
		ID: resp.Msg.SessionId,
	}); err != nil {
		log.Printf("セッション保存エラー: %v", err)
		renderError(w, "セッション保存エラー", "セッションの保存に失敗しました。", err.Error())
		return
	}

	log.Printf("セッション作成成功: %s", resp.Msg.SessionId)

	// マイナコネクトにリダイレクト
	http.Redirect(w, r, resp.Msg.RedirectUri, http.StatusSeeOther)
}

// handleCallback はマイナコネクトからのコールバックを処理する
func handleCallback(w http.ResponseWriter, r *http.Request) {
	// セッション情報を取得
	session, err := GetSession(r)
	if err != nil {
		renderError(w, "セッションエラー", "セッションの取得に失敗しました。", err.Error())
		return
	}

	if session == nil {
		renderError(w, "セッションエラー", "セッションが見つかりません。", "")
		return
	}

	log.Printf("コールバック受信: セッションID %s", session.ID)

	// ステータスを取得して適切なページにリダイレクト
	resp, err := client.GetSelfPersonalDataRequestStatus(context.Background(), connect.NewRequest(&mynaconnectv1.GetSelfPersonalDataRequestStatusRequest{
		SessionId: session.ID,
	}))
	if err != nil {
		renderError(w, "ステータス取得エラー", "ステータスの取得に失敗しました。", err.Error())
		return
	}

	// ステータスに応じてリダイレクト
	switch resp.Msg.Status {
	case mynaconnectv1.SelfPersonalDataRequestStatus_SELF_PERSONAL_DATA_REQUEST_STATUS_SUCCESS:
		http.Redirect(w, r, "/result", http.StatusSeeOther)
	case mynaconnectv1.SelfPersonalDataRequestStatus_SELF_PERSONAL_DATA_REQUEST_STATUS_NEED_TO_VERIFY_USER:
		http.Redirect(w, r, "/verify", http.StatusSeeOther)
	case mynaconnectv1.SelfPersonalDataRequestStatus_SELF_PERSONAL_DATA_REQUEST_STATUS_ERROR:
		renderError(w, "処理エラー", "処理中にエラーが発生しました。", fmt.Sprintf("Reason: %s", resp.Msg.ErrorReason))
	case mynaconnectv1.SelfPersonalDataRequestStatus_SELF_PERSONAL_DATA_REQUEST_STATUS_EXPIRED:
		renderError(w, "セッション期限切れ", "セッションの有効期限が切れました。最初からやり直してください。", "")
	case mynaconnectv1.SelfPersonalDataRequestStatus_SELF_PERSONAL_DATA_REQUEST_STATUS_PENDING:
		renderError(w, "処理中", "まだ処理が完了していません。しばらくしてから再度お試しください。", "")
	default:
		renderError(w, "不明なステータス", "不明なステータスが返されました。", fmt.Sprintf("Status: %v", resp.Msg.Status))
	}
}

// handleResult は取得結果を表示する
func handleResult(w http.ResponseWriter, r *http.Request) {
	// クッキーからセッションIDを取得
	sessionID, err := getSessionID(r)
	if err != nil {
		renderError(w, "セッションエラー", "セッションの取得に失敗しました。", err.Error())
		return
	}

	// 結果を取得
	resp, err := client.GetSelfPersonalDataRequestResult(context.Background(), connect.NewRequest(&mynaconnectv1.GetSelfPersonalDataRequestResultRequest{
		SessionId: sessionID,
	}))
	if err != nil {
		renderError(w, "結果取得エラー", "結果の取得に失敗しました。", err.Error())
		return
	}

	// rawデータを整形
	var rawData string
	if resp.Msg.Raw != "" {
		rawData = resp.Msg.Raw
	} else {
		rawData = "データなし"
	}

	// parsedフィールドの各Anyを対応するmessage型に変換
	parsedDataList := []string{}
	for _, any := range resp.Msg.Parsed {
		message, _ := any.UnmarshalNew()
		messageJSON, _ := json.MarshalIndent(message, "", "  ")
		parsedDataList = append(parsedDataList, string(messageJSON))
	}

	if err := templates.ExecuteTemplate(w, "result.html", map[string]any{
		"RawData":        rawData,
		"ParsedDataList": parsedDataList,
	}); err != nil {
		http.Error(w, "テンプレートのレンダリングに失敗しました", http.StatusInternalServerError)
		log.Printf("テンプレートエラー: %v", err)
	}
}

// handleVerify は本人確認画面を表示する
func handleVerify(w http.ResponseWriter, r *http.Request) {
	// クッキーからセッションIDを取得
	sessionID, err := getSessionID(r)
	if err != nil {
		renderError(w, "セッションエラー", "セッションの取得に失敗しました。", err.Error())
		return
	}

	// 本人確認情報を取得
	resp, err := client.GetUserIdentity(context.Background(), connect.NewRequest(&mynaconnectv1.GetUserIdentityRequest{
		SessionId: sessionID,
	}))
	if err != nil {
		renderError(w, "本人確認情報取得エラー", "本人確認情報の取得に失敗しました。", err.Error())
		return
	}

	// レスポンス全体をpretty JSON形式で表示
	identityJSON, _ := json.MarshalIndent(resp.Msg, "", "  ")

	if err := templates.ExecuteTemplate(w, "verify.html", map[string]any{
		"SessionID":    sessionID,
		"IdentityJSON": string(identityJSON),
	}); err != nil {
		http.Error(w, "テンプレートのレンダリングに失敗しました", http.StatusInternalServerError)
		log.Printf("テンプレートエラー: %v", err)
	}
}

// handleVerifyIdentity は本人確認の結果を処理する
func handleVerifyIdentity(w http.ResponseWriter, r *http.Request) {
	// クッキーからセッションIDを取得
	sessionID, err := getSessionID(r)
	if err != nil {
		renderError(w, "セッションエラー", "セッションの取得に失敗しました。", err.Error())
		return
	}

	// 結果を送信（デモアプリでは常に承認）
	resp, err := client.SubmitUserIdentityVerificationResult(context.Background(), connect.NewRequest(&mynaconnectv1.SubmitUserIdentityVerificationResultRequest{
		SessionId: sessionID,
		Ok:        true,
	}))
	if err != nil {
		renderError(w, "結果送信エラー", "本人確認結果の送信に失敗しました。", err.Error())
		return
	}

	// マイナコネクトにリダイレクト
	http.Redirect(w, r, resp.Msg.GetRedirectUri(), http.StatusSeeOther)
}

// renderError はエラーページを表示する
func renderError(w http.ResponseWriter, title, message, details string) {
	if err := templates.ExecuteTemplate(w, "error.html", map[string]any{
		"Title":   title,
		"Message": message,
		"Details": details,
	}); err != nil {
		http.Error(w, "エラーページの表示に失敗しました", http.StatusInternalServerError)
		log.Printf("テンプレートエラー: %v", err)
	}
}

// getSessionID はクッキーからセッションIDを取得する
func getSessionID(r *http.Request) (string, error) {
	// クッキーから取得
	session, err := GetSession(r)
	if err != nil {
		return "", err
	}

	if session == nil {
		return "", fmt.Errorf("session not found")
	}

	return session.ID, nil
}
