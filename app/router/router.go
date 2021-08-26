package router

import (
	"context"
	"fmt"
	"github.com/chenjiandongx/ginprom"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"github.com/v-mars/frame/db"
	"github.com/v-mars/frame/pkg/utils"
	"github.com/v-mars/sys/app/logger"
	"github.com/v-mars/sys/app/middleware"
	"github.com/v-mars/sys/app/models/sys"
	"github.com/v-mars/sys/app/sys/api"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	Engine *gin.Engine
)

func InitRouter(RunMode string, addr string)  {
	//RunMode := gin.DebugMode
	gin.SetMode(RunMode) //调试模式
	Engine = gin.New()

	// Prom
	Engine.Use(ginprom.PromMiddleware(&ginprom.PromOpts{}))

	if strings.ToLower(RunMode) == "debug" {
		Engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		Engine.GET("/api/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}


	gin.DefaultWriter = io.MultiWriter(logger.Gin.Writer()) // os.Stdout, logger.Gin.Writer()

	//错误日志审计
	Engine.Use(middleware.ErrorLogHandler())

	// 不存在方法
	Engine.NoMethod(middleware.NoMethodHandler())
	// 不存在路由
	Engine.NoRoute(middleware.NoRouteHandler())

	// 崩溃恢复
	//Engine.Use(gin.Recovery())
	Engine.Use(middleware.RecoveryMiddleware())
	Engine.Use(middleware.CORSMiddleware())
	// 日志
	Engine.Use(gin.Logger())

	// 登录验证 及信息提取
	var notCheckLoginUrlArr = []string{
		"/favicon.ico","/ping","/swager/*","/api/v1/swagger/*","/debug/pprof","/metrics","/api/code",
		"/gin/routes",
		"/api/v1/login","/api/v1/token","/api/v1/refresh-token","/api/login","/api/user/logout",
		"/api/v1/kubernetes/pod_log",
		"/api/v1/kubernetes/pod_terminal",
		"/api/v1/docker/container/terminal",
		"/api/v1/docker/container/log",
		"/api/v1/audit/ssh/replay",
		"/api/v1/fort/host/terminal",
		}
	Engine.Use(middleware.UserAuthMiddleware(middleware.AllowPathPrefixSkipper(notCheckLoginUrlArr...)))

	// 权限验证
	var notCheckPermissionUrlArr = []string{"/api/v1/user-info", "/api/menu/menubuttonlist"}
	notCheckPermissionUrlArr = append(notCheckPermissionUrlArr, notCheckLoginUrlArr...)
	Engine.Use(middleware.CasbinMiddleware(middleware.AllowPathPrefixSkipper(notCheckPermissionUrlArr...)))

	//setEmbedWeb(Engine)

	// 注册路由
	if err := RegisterRouter(Engine); err != nil{
		log.Fatal(err)
	}
	// Prometheus
	prometheusMonitor(Engine)
	//profile(Engine)

	//pprof.Register(Engine)

	srv := &http.Server{
		Addr:    addr,
		Handler: Engine,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()


	msg := fmt.Sprintf("服务启动成功，运行模式：%s，版本号：%s，进程号：%d", RunMode, "release", os.Getpid())
	fmt.Println(utils.Green(msg))
	fmt.Println(utils.Green("访问地址 http://"+addr))
	//log.Printf(utils.Green(msg, 1))
	//InitUpdatePermissionByGinRoutes()
	fmt.Println(utils.Green("[*] Waiting for messages. To exit press CTRL+C"))
	//select {}

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit,
		os.Interrupt, os.Kill,
		syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGKILL,
		syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP,
		syscall.SIGABRT)  // 此处不会阻塞
	sig:=<-quit  // 阻塞在此，当接收到上述两种信号时才会往下执行
	logger.Logger.Info(fmt.Sprintf("get signal %s, application will shutdown.", sig))
	logger.Logger.Info("start shutdown server...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {log.Fatal("server shutdown: ", err)}

	log.Println("server exited")
}


func InitUpdateApiByGinRoutes()  {
	perms := api.NewService(db.DB)
	var res []sys.Api
	routes := Engine.Routes()
	for _,v := range routes{
		var name string
		pathArray:=strings.Split(strings.Trim(v.Path, "/"), "/")
		if len(pathArray)>=4{
			name = fmt.Sprintf("%s:%s:%s", pathArray[2],pathArray[3], strings.ToLower(v.Method))
		} else if len(pathArray)==3{
			name = fmt.Sprintf("%s:%s", pathArray[2], strings.ToLower(v.Method))
		}else if len(pathArray)==2{
			name = fmt.Sprintf("%s:%s", pathArray[1], strings.ToLower(v.Method))
		} else if len(pathArray)==1{
			name = fmt.Sprintf("%s:%s", pathArray[0], strings.ToLower(v.Method))
		}
		res = append(res, sys.Api{Name: name,Method: v.Method, Path: v.Path})
	}
	err := perms.UpdateApi(res)
	if err!=nil{
		logger.Logger.Errorf("InitUpdateApiByGinRoutes: %s", err)
	}
}