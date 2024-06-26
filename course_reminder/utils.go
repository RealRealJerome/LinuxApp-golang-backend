package courseReminder

type ReminderReq struct {
	Form       string `json:"form"`
	RemindTime int    `json:"remind_time"`
}
