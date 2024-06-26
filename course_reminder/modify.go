package courseReminder

import (
	"database/sql"
	"encoding/json"
	mysqlUtil "github.com/A5-golang-backend/mysql"
	"github.com/A5-golang-backend/user"
	"log"
	"net/http"
)

func ModifyCourseReminder(w http.ResponseWriter, r *http.Request) {
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
	vars := r.URL.Query()
	name := vars.Get("name")
	var reminder ReminderReq
	err := json.NewDecoder(r.Body).Decode(&reminder)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	// 查询数据
	selectQuery := "SELECT cr.id,u.id FROM course c LEFT JOIN course_teacher ct ON c.id = ct.course_id LEFT JOIN user u ON ct.teacher_id = u.id LEFT JOIN course_reminder cr ON c.course_reminder_id = cr.id  WHERE c.name = ?"
	selectedId, selectedTeacherId := 0, 0
	err = mysqlUtil.DB.QueryRow(selectQuery, name).Scan(&selectedId, &selectedTeacherId)
	if err == sql.ErrNoRows || selectedId == 0 {
		// 构建响应
		response := struct {
			Code    int
			Message string
		}{
			Code:    400,
			Message: "所选课程不存在或已被删除",
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
	// 更新数据的 SQL 语句
	query := "UPDATE course_reminder cr JOIN course c ON cr.id = c.course_reminder_id SET cr.form = ?, cr.time_minutes = ?  WHERE c.name = ?"
	// 执行更新操作
	_, err = mysqlUtil.DB.Exec(query, reminder.Form, reminder.RemindTime, name)
	if err != nil {
		log.Fatal(err)
	}
	// 构建响应
	response := struct {
		Code    int
		Message string
	}{
		Code:    200,
		Message: "更新课程提醒成功",
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
