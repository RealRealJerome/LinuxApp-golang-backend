package main

import (
	"fmt"
	"github.com/A5-golang-backend/sms"
	"github.com/A5-golang-backend/user"
	"github.com/gorilla/mux"
	"net/http"
)

func SetRoute(router *mux.Router) {
	// 设置路由和处理函数
	router.HandleFunc("/user/{phoneNum}/signIn:password", user.HandleSignInPW).Methods("POST")
	router.HandleFunc("/user/{phoneNum}/signIn:smsCode", user.HandleSignInSMS).Methods("POST")
	router.HandleFunc("/user/{phoneNum}/chPw", user.HandleChangePW).Methods("PUT")
	router.HandleFunc("/user/{phoneNum}/signUp:teacher", user.HandleSignUpTeacher).Methods("POST")
	router.HandleFunc("/user/{phoneNum}/smsCode", sms.GetSmsCode).Methods("GET")
	router.HandleFunc("/user/{phoneNum}/smsCode", sms.VerifySmsCode).Methods("POST")
}
func main() {
	// 创建新的路由器
	router := mux.NewRouter()
	SetRoute(router)

	// 指定监听端口
	port := ":6666"
	fmt.Printf("Starting server on port %s...\n", port)

	// 启动HTTP服务器，监听指定端口
	err := http.ListenAndServe(port, router)
	if err != nil {
		fmt.Printf("Server failed to start: %s\n", err)
	}

}
