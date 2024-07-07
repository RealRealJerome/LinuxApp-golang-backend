package course

import (
	"encoding/json"
	mysqlUtil "github.com/A5-golang-backend/mysql"
	"github.com/A5-golang-backend/user"
	"log"
	"net/http"
)

// 定义请求体结构
type CreateReq struct {
	Name      string  `json:"name"`
	Time      Time    `json:"time"`
	Classroom string  `json:"classroom"`
	Credit    float64 `json:"credit"`
}

type Time struct {
	Weeks   []int `json:"weeks"`
	Days    []int `json:"days"`
	DayTime []int `json:"dayTime"`
}

func CreateCourse(w http.ResponseWriter, r *http.Request) {
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
	// 解析请求体
	var body CreateReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	phoneNum := user.GetToken(token)
	// 插入数据的 SQL 语句
	query := "INSERT INTO course (name, time_weeks,time_days,day_time,credit,classroom) VALUES (?, ?, ?, ?, ?, ?)"
	// 执行插入操作
	result, err := mysqlUtil.DB.Exec(query, body.Name, IntArr2Str(body.Time.Weeks), IntArr2Str(body.Time.Days), IntArr2Str(body.Time.DayTime), body.Credit, body.Classroom)
	if err != nil {
		log.Fatal(err)
	}
	// 获取插入的记录的 ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	// 关联teacher
	// 查询数据
	selectQuery := "SELECT id FROM user WHERE phone_num = ?"
	selectedId := 0
	err = mysqlUtil.DB.QueryRow(selectQuery, phoneNum).Scan(&selectedId)
	// 插入数据的 SQL 语句
	query = "INSERT INTO course_teacher(course_id,teacher_id) VALUES (?, ?)"
	// 执行插入操作
	result, err = mysqlUtil.DB.Exec(query, lastInsertID, selectedId)
	if err != nil {
		log.Fatal(err)
	}
	// 构建响应
	response := struct {
		Code    int
		Message string
	}{
		Code:    200,
		Message: "课程创建成功",
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
