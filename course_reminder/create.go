package courseReminder

import (
	"database/sql"
	"encoding/json"
	mysqlUtil "github.com/A5-golang-backend/mysql"
	"github.com/A5-golang-backend/user"
	"log"
	"net/http"
)

func CreateCourseReminder(w http.ResponseWriter, r *http.Request) {
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
	var selectedId, selectedTeacherId sql.NullInt32
	row := mysqlUtil.DB.QueryRow(selectQuery, name)
	err = row.Scan(&selectedId, &selectedTeacherId)
	if err == sql.ErrNoRows || !selectedId.Valid {
		// 插入数据的 SQL 语句
		query := "INSERT INTO course_reminder(form,time_minutes,teacher_id) VALUES (?, ?, ?)"
		// 执行插入操作
		result, err := mysqlUtil.DB.Exec(query, reminder.Form, reminder.RemindTime, selectedTeacherId)
		if err != nil {
			log.Fatal(err)
		}
		lastInsertID, err := result.LastInsertId()
		// 更新数据的 SQL 语句
		query = "UPDATE course SET course_reminder_id = ?"
		// 执行更新操作
		result, err = mysqlUtil.DB.Exec(query, lastInsertID)
		if err != nil {
			log.Fatal(err)
		}
		// 构建响应
		response := struct {
			Code    int
			Message string
		}{
			Code:    200,
			Message: "课程提醒创建成功",
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
	// 构建响应
	response := struct {
		Code    int
		Message string
	}{
		Code:    400,
		Message: "每个课程只可设置一个课程提醒",
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
