package course

import (
	"database/sql"
	"encoding/json"
	mysqlUtil "github.com/A5-golang-backend/mysql"
	"github.com/A5-golang-backend/user"
	"log"
	"net/http"
)

func ModifyCourse(w http.ResponseWriter, r *http.Request) {
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
	oldName := vars.Get("old_name")
	// 解析请求体
	var body CreateReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// 查询数据
	selectQuery := "SELECT id FROM course WHERE name = ?"
	selectedId := 0
	err := mysqlUtil.DB.QueryRow(selectQuery, oldName).Scan(&selectedId)
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
	query := "UPDATE course SET name = ?, time_weeks = ?,time_days = ?,day_time = ?,credit = ? WHERE name = ?"
	// 执行更新操作
	_, err = mysqlUtil.DB.Exec(query, body.Name, IntArr2Str(body.Time.Weeks), IntArr2Str(body.Time.Days), IntArr2Str(body.Time.DayTime), body.Credit, oldName)
	if err != nil {
		log.Fatal(err)
	}
	// 构建响应
	response := struct {
		Code    int
		Message string
	}{
		Code:    200,
		Message: "课程编辑成功",
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
