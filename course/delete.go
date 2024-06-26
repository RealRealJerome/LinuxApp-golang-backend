package course

import (
	"database/sql"
	"encoding/json"
	mysqlUtil "github.com/A5-golang-backend/mysql"
	"github.com/A5-golang-backend/user"
	"log"
	"net/http"
)

func DeleteCourse(w http.ResponseWriter, r *http.Request) {
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
	// 查询数据
	selectQuery := "SELECT id,course_reminder_id FROM course WHERE name = ?"
	selectedId, selectedReminderId := 0, 0
	err := mysqlUtil.DB.QueryRow(selectQuery, name).Scan(&selectedId, &selectedReminderId)
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
	// 删除数据的 SQL 语句
	query := "DELETE FROM course WHERE name = ?"
	// 执行删除操作
	_, err = mysqlUtil.DB.Exec(query, name)
	if err != nil {
		log.Fatal(err)
	}
	// 删除数据的 SQL 语句
	query = "DELETE FROM course_teacher WHERE course_id = ?"
	// 执行删除操作
	_, err = mysqlUtil.DB.Exec(query, selectedId)
	if err != nil {
		log.Fatal(err)
	}
	// 删除数据的 SQL 语句
	query = "DELETE FROM course_reminder WHERE id = ?"
	// 执行删除操作
	_, err = mysqlUtil.DB.Exec(query, selectedReminderId)
	if err != nil {
		log.Fatal(err)
	}
	// 构建响应
	response := struct {
		Code    int
		Message string
	}{
		Code:    200,
		Message: "课程删除成功",
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
