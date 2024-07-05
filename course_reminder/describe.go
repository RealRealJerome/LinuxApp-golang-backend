package courseReminder

import (
	"database/sql"
	"encoding/json"
	mysqlUtil "github.com/A5-golang-backend/mysql"
	"github.com/A5-golang-backend/user"
	"log"
	"net/http"
)

type DescribeReply struct {
	Name       string `json:"name"`
	Form       string `json:"form"`
	RemindTime int    `json:"remind_time"`
}

func DescribesCourseReminderByTeacherName(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	if !user.CheckToken(token) {
		// 构建响应
		response := struct {
			Code    int
			Message string
		}{
			Code:    403,
			Message: "登录信息无效或已失效，请重新登录",
		}
		// 4. 序列化返回体为 JSON 格式
		responseJSON, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 返回JSON响应
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseJSON)
		return
	}
	phoneNum := user.GetToken(token)
	// 查询数据
	selectQuery := "SELECT cr.form,cr.time_minutes,ud.name FROM course_reminder cr LEFT JOIN course c ON c.course_reminder_id = cr.id LEFT JOIN course_teacher ct ON c.id = ct.course_id LEFT JOIN user t ON ct.teacher_id = t.id LEFT JOIN user_detail ud ON t.user_detail_id = ud.id WHERE u.phone_num = ?"
	var selectedTimeMinutes sql.NullInt32
	var selectedForm, selectedCourseName sql.NullString
	rows, err := mysqlUtil.DB.Query(selectQuery, phoneNum)
	if err != nil {
		log.Fatal(err)
	}
	describeReplies := make([]*DescribeReply, 0)
	for rows.Next() {
		err := rows.Scan(&selectedForm, &selectedTimeMinutes, &selectedCourseName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		describeReplies = append(describeReplies, &DescribeReply{
			Name:       selectedCourseName.String,
			Form:       selectedForm.String,
			RemindTime: int(selectedTimeMinutes.Int32),
		})
	}
	// 构建响应
	response := struct {
		Code            int
		Message         string
		DescribeReplies []*DescribeReply
	}{
		Code:            200,
		Message:         "查询成功",
		DescribeReplies: describeReplies,
	}
	// 4. 序列化返回体为 JSON 格式
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 返回JSON响应
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}
