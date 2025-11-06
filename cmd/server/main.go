package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"cloudsaver/embed"
	"cloudsaver/internal/api/handler"
	"cloudsaver/internal/api/middleware"
	"cloudsaver/internal/config"
	"cloudsaver/internal/pkg/cloud"
	"cloudsaver/internal/repository"
	"cloudsaver/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("加载配置失败，使用默认配置: %v", err)
		cfg = &config.Config{
			Server: config.ServerConfig{
				Port: 8080,
				Mode: "release",
			},
			Database: config.DatabaseConfig{
				Path: "./data/cloudsaver.db",
			},
			JWT: config.JWTConfig{
				Secret:      "your-secret-key-change-in-production",
				ExpireHours: 168,
			},
		}
	}

	// 初始化数据库
	db, err := config.InitDatabase(cfg.Database.Path)
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 初始化仓库
	userRepo := repository.NewUserRepository(db)
	settingRepo := repository.NewSettingRepository(db)

	// 初始化服务
	userService := service.NewUserService(userRepo, settingRepo, cfg.JWT.Secret, cfg.JWT.ExpireHours)
	settingService := service.NewSettingService(settingRepo)
	cloud115Service := cloud.NewCloud115Service(settingRepo)
	quarkService := cloud.NewQuarkService(settingRepo)

	// 初始化处理器
	userHandler := handler.NewUserHandler(userService)
	settingHandler := handler.NewSettingHandler(settingService)
	cloud115Handler := handler.NewCloud115Handler(cloud115Service)
	quarkHandler := handler.NewQuarkHandler(quarkService)

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 创建路由
	router := gin.New()

	// 注册全局中间件
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggerMiddleware())

	// API路由组
	api := router.Group("/api")
	{
		// 用户路由（无需认证）
		api.POST("/user/login", userHandler.Login)
		api.POST("/user/register", userHandler.Register)

		// 需要认证的路由
		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware(cfg.JWT.Secret, userRepo))
		{
			// 设置路由
			auth.GET("/setting/get", settingHandler.Get)
			auth.POST("/setting/save", settingHandler.Save)

			// 115网盘路由
			auth.GET("/cloud115/share-info", cloud115Handler.GetShareInfo)
			auth.GET("/cloud115/folders", cloud115Handler.GetFolderList)
			auth.POST("/cloud115/save", cloud115Handler.SaveFile)

			// 夸克网盘路由
			auth.GET("/quark/share-info", quarkHandler.GetShareInfo)
			auth.GET("/quark/folders", quarkHandler.GetFolderList)
			auth.POST("/quark/save", quarkHandler.SaveFile)
		}
	}

	// 嵌入前端静态资源
	distFS, err := fs.Sub(embed.FrontendAssets, "dist")
	if err != nil {
		log.Printf("警告: 无法加载前端资源: %v", err)
	} else {
		// 使用 StaticFS 为 assets 文件提供服务
		router.StaticFS("/assets", http.FS(distFS))
		
		// 处理 SPA 路由
		router.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path

			// API路由已经处理，这里不需要再处理
			if strings.HasPrefix(path, "/api/") {
				c.JSON(http.StatusNotFound, gin.H{"error": "API not found"})
				return
			}

			// 尝试读取文件
			filePath := strings.TrimPrefix(path, "/")
			if filePath == "" {
				filePath = "index.html"
			}

			// 尝试打开文件
			file, err := distFS.Open(filePath)
			if err != nil {
				// 文件不存在，返回 index.html（SPA路由）
				data, err := fs.ReadFile(distFS, "index.html")
				if err != nil {
					c.String(http.StatusNotFound, "index.html not found")
					return
				}
				c.Data(http.StatusOK, "text/html; charset=utf-8", data)
				return
			}
			file.Close()

			// 文件存在，读取并返回
			data, err := fs.ReadFile(distFS, filePath)
			if err != nil {
				// 读取失败，返回 index.html
				data, _ = fs.ReadFile(distFS, "index.html")
				c.Data(http.StatusOK, "text/html; charset=utf-8", data)
				return
			}

			// 根据文件扩展名设置 Content-Type
			contentType := "text/plain"
			switch {
			case strings.HasSuffix(filePath, ".html"):
				contentType = "text/html; charset=utf-8"
			case strings.HasSuffix(filePath, ".css"):
				contentType = "text/css; charset=utf-8"
			case strings.HasSuffix(filePath, ".js"):
				contentType = "application/javascript; charset=utf-8"
			case strings.HasSuffix(filePath, ".json"):
				contentType = "application/json; charset=utf-8"
			case strings.HasSuffix(filePath, ".png"):
				contentType = "image/png"
			case strings.HasSuffix(filePath, ".jpg"), strings.HasSuffix(filePath, ".jpeg"):
				contentType = "image/jpeg"
			case strings.HasSuffix(filePath, ".svg"):
				contentType = "image/svg+xml"
			case strings.HasSuffix(filePath, ".ico"):
				contentType = "image/x-icon"
			}

			c.Data(http.StatusOK, contentType, data)
		})
	}

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("服务器启动在端口 %d", cfg.Server.Port)
	if err := router.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
