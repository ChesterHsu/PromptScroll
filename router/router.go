package router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/swaggo/swag"
	"promptscroll/config"
	"promptscroll/controller"
	_ "promptscroll/docs"
	"promptscroll/middleware"
)

// InitRouter 初始化路由器
func InitRouter() http.Handler {
	// 載入 Swagger 設定
	swagCfg := config.LoadSwaggerConfig()

	// 設定路由器
	r := mux.NewRouter()

	// 註冊中介軟體
	r.Use(middleware.SetJSONMiddleware)
	r.Use(middleware.LogMiddleware)

	// 設定路由規則
	// 顯示 Swagger API 文件
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// 設定驗證路由
	authRouter := r.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", controller.RegisterHandler).Methods(http.MethodPost)
	authRouter.HandleFunc("/login", controller.LoginHandler).Methods(http.MethodPost)

	// 設定使用者路由
	userRouter := r.PathPrefix("/users").Subrouter()
	userRouter.Use(middleware.JWTMiddleware)
	userRouter.HandleFunc("/{id:[0-9]+}", controller.GetUserHandler).Methods(http.MethodGet)

	// 載入 Swagger 設定
	swagCfg := config.LoadSwaggerConfig()

	// 顯示 Swagger API 文件
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// 設定 Swagger API 路由
	if swagCfg.Enabled {
		swagRouter := r.PathPrefix("/api").Subrouter()

		swagRouter.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
			doc, err := swag.ReadDoc()
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(doc)
		})

		swagRouter.HandleFunc("/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
			doc, err := swag.ReadDoc()
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/yaml")
			w.WriteHeader(http.StatusOK)
			w.Write(doc)
		})
	}

	// 預設路由，顯示歡迎畫面
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome to PromptScroll!")
	})

	return r
}
