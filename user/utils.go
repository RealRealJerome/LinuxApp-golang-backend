package user

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
	"unicode"
)

type TeacherInfo struct {
	Name      string `json:"name"`
	College   string `json:"college"`
	School    string `json:"school"`
	TeacherId string `json:"teacher_id"`
	Password  string `json:"password"`
}
type OldNewPW struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
type Password struct {
	Password string `json:"password"`
}
type SMSCode struct {
	SMSCode string `json:"smsCode"`
}
type UserDetail struct {
	phoneNum string
	time     time.Time
}

var TokenMap = make(map[string]UserDetail)
var expiredTime = 600 * time.Second

// Sha256 加密 -- 随时间而变
func Sha256Time(input string) string {
	// 计算 SHA-256 哈希值
	hash := sha256.New()
	input += time.Now().String()
	hash.Write([]byte(input))
	hashBytes := hash.Sum(nil)

	// 将哈希值转换为 32 位的字母和数字混合字符串
	return hex.EncodeToString(hashBytes)[:32]
}

// 每次结果一样
func Sha256(input string) string {
	// 计算 SHA-256 哈希值
	hash := sha256.New()
	hash.Write([]byte(input))
	hashBytes := hash.Sum(nil)

	// 将哈希值转换为 32 位的字母和数字混合字符串
	return hex.EncodeToString(hashBytes)[:32]
}

// generateToken 模拟生成访问令牌的函数
func generateToken(phoneNum string) string {
	token := Sha256Time(phoneNum)
	TokenMap[token] = UserDetail{
		phoneNum: phoneNum,
		time:     time.Now(),
	}
	return token
}

// checkToken 校验token是否有效
func checkToken(token string) bool {
	val, ok := TokenMap[token]
	if !ok {
		return false
	}
	if time.Now().Sub(val.time) > expiredTime {
		delete(TokenMap, token)
		return false
	}
	return true
}

// ValidatePassword 校验密码是否至少8位字符，包含数字和字母
func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasLetter := false
	hasNumber := false

	for _, char := range password {
		if unicode.IsLetter(char) {
			hasLetter = true
		}
		if unicode.IsDigit(char) {
			hasNumber = true
		}
		if hasLetter && hasNumber {
			return true
		}
	}

	return false
}
