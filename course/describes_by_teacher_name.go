package course

import (
	"database/sql"
	"encoding/json"
	mysqlUtil "github.com/A5-golang-backend/mysql"
	"github.com/A5-golang-backend/user"
	"log"
	"net/http"
)

type DescribeReply struct {
	Name      string                `json:"name"`
	Time      Time                  `json:"time"`
	Teacher   string                `json:"teacher"`
	Classroom string                `json:"classroom"`
	Credit    float64               `json:"credit"`
	Reminder  DescribeReminderReply `json:"reminder"`
}
type DescribeReminderReply struct {
	HasReminder bool   `json:"has_reminder"`
	Form        string `json:"form"`
	RemindTime  int    `json:"remind_time"`
}
type DescribesReminderReplyNode struct {
	Name       string `json:"name"`
	Form       string `json:"form"`
	RemindTime string `json:"remind_time"`
}
type DescribesReq struct {
	Week Week   `json:"week"`
	Name string `json:"name"`
}
type Week struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// DescribeCourse 根据名字查询某个课程详细信息
func DescribeCourse(w http.ResponseWriter, r *http.Request) {
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
	selectQuery := "SELECT c.name,c.time_weeks,c.time_days,c.day_time,c.credit,c.classroom,ud.name,cr.form,cr.time_minutes FROM course c LEFT JOIN course_teacher ct ON c.id = ct.course_id LEFT JOIN user u ON ct.teacher_id = u.id LEFT JOIN user_detail ud ON u.user_detail_id = ud.id LEFT JOIN course_reminder cr on c.course_reminder_id = cr.id and cr.teacher_id = u.id WHERE c.name = ?"
	selectedName, selectedTimeWeeks, selectedTimeDays, selectedDayTime, selectedCredit, selectedClassroom := "", "", "", "", 0.0, ""
	var selectedTeacherName, selectedForm sql.NullString
	var selectedTimeMinutes sql.NullInt32
	err := mysqlUtil.DB.QueryRow(selectQuery, name).Scan(&selectedName, &selectedTimeWeeks, &selectedTimeDays, &selectedDayTime, &selectedCredit, &selectedClassroom, &selectedTeacherName, &selectedForm, &selectedTimeMinutes)
	if err != nil {
		log.Fatal(err)
	}
	hasReminder := true
	if !selectedForm.Valid {
		hasReminder = false
	}
	// 构建响应
	response := struct {
		Code          int
		Message       string
		DescribeReply DescribeReply
	}{
		Code:    200,
		Message: "查询成功",
		DescribeReply: DescribeReply{
			Name: selectedName,
			Time: Time{
				Weeks:   Str2IntArr(selectedTimeWeeks),
				Days:    Str2IntArr(selectedTimeDays),
				DayTime: Str2IntArr(selectedDayTime),
			},
			Teacher:   selectedTeacherName.String,
			Classroom: selectedClassroom,
			Credit:    selectedCredit,
			Reminder: DescribeReminderReply{
				HasReminder: hasReminder,
				Form:        selectedForm.String,
				RemindTime:  int(selectedTimeMinutes.Int32),
			},
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
func DescribeCourses(w http.ResponseWriter, r *http.Request) {
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
	var body DescribesReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	weeks := make([]int, 0)
	for i := body.Week.Start; i < body.Week.End; i++ {
		weeks = append(weeks, i)
	}
	// 查询数据
	selectQuery := "SELECT c.name,c.time_weeks,c.time_days,c.day_time,c.credit,c.classroom,ud.name,cr.form,cr.time_minutes FROM course c LEFT JOIN course_teacher ct ON c.id = ct.course_id LEFT JOIN user u ON ct.teacher_id = u.id LEFT JOIN user_detail ud ON u.user_detail_id = ud.id LEFT JOIN course_reminder cr on c.course_reminder_id = cr.id and cr.teacher_id = u.id where c.name = ?"
	rows, err := mysqlUtil.DB.Query(selectQuery, body.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	describeReplies := make([]*DescribeReply, 0)
	for rows.Next() {
		selectedName, selectedTimeWeeks, selectedTimeDays, selectedDayTime, selectedCredit, selectedClassroom := "", "", "", "", 0.0, ""
		var selectedTeacherName, selectedForm sql.NullString
		var selectedTimeMinutes sql.NullInt32
		err := rows.Scan(&selectedName, &selectedTimeWeeks, &selectedTimeDays, &selectedDayTime, &selectedCredit, &selectedClassroom, &selectedTeacherName, &selectedForm, &selectedTimeMinutes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !JudgeWeekContains(weeks, selectedTimeWeeks) {
			continue
		}
		hasReminder := true
		if !selectedForm.Valid {
			hasReminder = false
		}
		// 构建响应
		response := &DescribeReply{
			Name: selectedName,
			Time: Time{
				Weeks:   Str2IntArr(selectedTimeWeeks),
				Days:    Str2IntArr(selectedTimeDays),
				DayTime: Str2IntArr(selectedDayTime),
			},
			Teacher:   selectedTeacherName.String,
			Classroom: selectedClassroom,
			Credit:    selectedCredit,
			Reminder: DescribeReminderReply{
				HasReminder: hasReminder,
				Form:        selectedForm.String,
				RemindTime:  int(selectedTimeMinutes.Int32),
			},
		}
		describeReplies = append(describeReplies, response)
	}
	// 构建响应
	response := struct {
		Code           int
		Message        string
		DescribesReply []*DescribeReply
	}{
		Code:           200,
		Message:        "查询成功",
		DescribesReply: describeReplies,
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
