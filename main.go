package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"

	"buf.build/gen/go/pocketsign/apis/connectrpc/go/pocketsign/mynaconnect/v1/mynaconnectv1connect"
	spdv1 "buf.build/gen/go/pocketsign/apis/protocolbuffers/go/pocketsign/mynaconnect/spd/v1"
	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoregistry"
)

var (
	client          mynaconnectv1connect.SelfPersonalDataServiceClient
	templates       = template.Must(template.ParseGlob("templates/*.html"))
	mynaConnectHost = "https://mynaconnect.mock.p8n.app"
	apiToken        = os.Getenv("MYNA_CONNECT_API_TOKEN")  // PocketSign Platformから取得したテナントAPIトークンを設定
	serviceID       = os.Getenv("MYNA_CONNECT_SERVICE_ID") // PocketSign Platformで作成した自己情報取得APIのAPIサービスIDを設定
	callbackBaseURL = "http://localhost:3000"
)

func main() {
	// マイナコネクトクライアントの作成
	client = mynaconnectv1connect.NewSelfPersonalDataServiceClient(
		http.DefaultClient,
		mynaConnectHost,
		connect.WithInterceptors(
			connect.UnaryInterceptorFunc(
				func(unaryFunc connect.UnaryFunc) connect.UnaryFunc {
					return func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
						request.Header().Set("Authorization", "Bearer "+apiToken)
						return unaryFunc(ctx, request)
					}
				},
			),
		),
	)

	// ルートの設定
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handleIndex)
	mux.HandleFunc("GET /start", handleStart)
	mux.HandleFunc("GET /callback", handleCallback)
	mux.HandleFunc("GET /result", handleResult)
	mux.HandleFunc("GET /verify", handleVerify)
	mux.HandleFunc("POST /verify-identity", handleVerifyIdentity)
	mux.HandleFunc("GET /common.css", handleCSS)

	// サーバーの起動
	log.Printf("サーバーを起動します: http://localhost:3000")
	log.Printf("マイナコネクトホスト: %s", mynaConnectHost)
	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatalf("サーバーの起動に失敗しました: %v", err)
	}
}

func init() {
	// 実際のプロダクトでは必要になる型のみ登録すれば問題ない
	for _, msgType := range []proto.Message{
		&spdv1.TM00000000000001{},
		&spdv1.TM00000000000002{},
		&spdv1.TM00000000000003{},
		&spdv1.TM00000000000004{},
		&spdv1.TM00000000000008{},
		&spdv1.TM00000000000009{},
		&spdv1.TM00000000000010{},
		&spdv1.TM00000000000012{},
		&spdv1.TM00000000000015{},
		&spdv1.TM00000000000016{},
		&spdv1.TM00000000000017{},
		&spdv1.TM00000000000018{},
		&spdv1.TM00000000000019{},
		&spdv1.TM00000000000020{},
		&spdv1.TM00000000000021{},
		&spdv1.TM00000000000023{},
		&spdv1.TM00000000000024{},
		&spdv1.TM00000000000025{},
		&spdv1.TM00000000000026{},
		&spdv1.TM00000000000029{},
		&spdv1.TM00000000000030{},
		&spdv1.TM00000000000031{},
		&spdv1.TM00000000000033{},
		&spdv1.TM00000000000034{},
		&spdv1.TM00000000000035{},
		&spdv1.TM00000000000036{},
		&spdv1.TM00000000000037{},
		&spdv1.TM00000000000038{},
		&spdv1.TM00000000000039{},
		&spdv1.TM00000000000043{},
		&spdv1.TM00000000000044{},
		&spdv1.TM00000000000045{},
		&spdv1.TM00000000000046{},
		&spdv1.TM00000000000047{},
		&spdv1.TM00000000000049{},
		&spdv1.TM00000000000050{},
		&spdv1.TM00000000000051{},
		&spdv1.TM00000000000052{},
		&spdv1.TM00000000000053{},
		&spdv1.TM00000000000054{},
		&spdv1.TM00000000000055{},
		&spdv1.TM00000000000056{},
		&spdv1.TM00000000000057{},
		&spdv1.TM00000000000058{},
		&spdv1.TM00000000000059{},
		&spdv1.TM00000000000064{},
		&spdv1.TM00000000000065{},
		&spdv1.TM00000000000068{},
		&spdv1.TM00000000000069{},
		&spdv1.TM00000000000074{},
		&spdv1.TM00000000000075{},
		&spdv1.TM00000000000078{},
		&spdv1.TM00000000000080{},
		&spdv1.TM00000000000081{},
		&spdv1.TM00000000000082{},
		&spdv1.TM00000000000083{},
		&spdv1.TM00000000000084{},
		&spdv1.TM00000000000085{},
		&spdv1.TM00000000000086{},
		&spdv1.TM00000000000087{},
		&spdv1.TM00000000000089{},
		&spdv1.TM00000000000090{},
		&spdv1.TM00000000000091{},
		&spdv1.TM00000000000092{},
		&spdv1.TM00000000000093{},
		&spdv1.TM00000000000094{},
		&spdv1.TM00000000000095{},
		&spdv1.TM00000000000096{},
		&spdv1.TM00000000000097{},
		&spdv1.TM00000000000098{},
		&spdv1.TM00000000000099{},
		&spdv1.TM00000000000100{},
		&spdv1.TM00000000000101{},
		&spdv1.TM00000000000102{},
		&spdv1.TM00000000000103{},
		&spdv1.TM00000000000104{},
		&spdv1.TM00000000000105{},
		&spdv1.TM00000000000106{},
		&spdv1.TM00000000000107{},
		&spdv1.TM00000000000108{},
	} {
		if _, err := protoregistry.GlobalTypes.FindMessageByName(msgType.ProtoReflect().Type().Descriptor().FullName()); err != nil {
			protoregistry.GlobalTypes.RegisterMessage(msgType.ProtoReflect().Type())
		}
	}
}
