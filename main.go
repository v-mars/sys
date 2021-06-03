package main

import (
	"github.com/v-mars/sys/cmd/server"
	_ "github.com/v-mars/sys/docs"
)


// Bearer
//securityDefinitions:
//  APIKey:
//    type: apiKey
//    name: Authorization
//    in: header
//security:
//  - APIKey: []


// @title OPS Go Docs
// @version v1.0
// @description Mars System Manage api v1
// @description 200: '成功'
// @description Authorization Bearer token

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @schema Bearer

// @contact.name ocean.zhang
// @contact.url
// @contact.email 429472406@qq.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:5000

//// @BasePath /api/v1
func main(){
	/*flag.Parse()
	if h {
		flag.Usage()
		os.Exit(-1)
	}
	*/

	server.Run()
}