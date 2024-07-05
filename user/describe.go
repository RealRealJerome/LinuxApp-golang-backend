package user

import (
	"database/sql"
	"encoding/json"
	mysqlUtil "github.com/A5-golang-backend/mysql"
	"log"
	"net/http"
)

type DescribeReply struct {
	PhoneNum string `json:"phoneNum"`
	Name     string `json:"name"`
	College  string `json:"college"`
	School   string `json:"school"`
}

func Describe(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	if !CheckToken(token) {
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
	phoneNum := GetToken(token)
	// 查询数据
	selectQuery := "SELECT ud.name,ud.college,ud.school from user u left join user_detail ud ON u.user_detail_id = ud.id WHERE u.phone_num = ?"
	var selectedName, selectedCollege, selectedSchool sql.NullString
	err := mysqlUtil.DB.QueryRow(selectQuery, phoneNum).Scan(&selectedName, &selectedCollege, &selectedSchool)
	if err != nil {
		log.Fatal(err)
	}
	response := struct {
		Code          int
		Message       string
		DescribeReply *DescribeReply
	}{
		Code:    200,
		Message: "查询成功",
		DescribeReply: &DescribeReply{
			PhoneNum: phoneNum,
			Name:     selectedName.String,
			College:  selectedCollege.String,
			School:   selectedSchool.String,
		},
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
