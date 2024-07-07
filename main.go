package main

import (
	"fmt"
	"github.com/A5-golang-backend/course"
	courseReminder "github.com/A5-golang-backend/course_reminder"
	"github.com/A5-golang-backend/sms"
	sparkModel "github.com/A5-golang-backend/spark_model"
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
	router.HandleFunc("/user:detail", user.Describe).Methods("GET")
	router.HandleFunc("/user/{phoneNum}/smsCode", sms.GetSmsCode).Methods("GET")
	router.HandleFunc("/user/{phoneNum}/smsCode", sms.VerifySmsCode).Methods("POST")
	router.HandleFunc("/course:create", course.CreateCourse).Methods("POST")
	router.HandleFunc("/course:edit", course.ModifyCourse).Methods("PUT")
	router.HandleFunc("/course:delete", course.DeleteCourse).Methods("DELETE")
	router.HandleFunc("/course:detail", course.DescribeCourse).Methods("GET")
	router.HandleFunc("/course:list", course.DescribeCourses).Methods("GET")
	router.HandleFunc("/courseReminder:create", courseReminder.CreateCourseReminder).Methods("POST")
	router.HandleFunc("/courseReminder:modify", courseReminder.ModifyCourseReminder).Methods("PUT")
	router.HandleFunc("/courseReminder:delete", courseReminder.DeleteCourseReminder).Methods("DELETE")
	router.HandleFunc("/courseReminder:list", courseReminder.DescribesCourseReminderByTeacherName).Methods("GET")
	router.HandleFunc("/assistant:spark", sparkModel.Describe).Methods("POST")
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
