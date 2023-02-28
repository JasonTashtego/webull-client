module quantfu.com/webull-client

require (
	github.com/eclipse/paho.mqtt.golang v1.4.1
	github.com/google/uuid v1.3.0
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.3.0
	golang.org/x/oauth2 v0.0.0-20221014153046-6fdb5e3db783
	quantfu.com/webull/client v0.0.0-00010101000000-000000000000
	quantfu.com/webull/openapi v0.0.0

)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.0.0-20221017152216-f25eb7ecb193 // indirect
	golang.org/x/sync v0.1.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

go 1.19

replace quantfu.com/webull/openapi => ../openapi

replace quantfu.com/webull/client => ../client
