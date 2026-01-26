package config

// 配置（实际项目应该从配置文件或环境变量读取）
var (
	// 数据库配置
	DBHost     = "localhost"
	DBPort     = "3306"
	DBUser     = "root"
	DBPassword = "123456" // 改成你的密码
	DBName     = "test1"

	// JWT 配置
	JWTSecret = "your-secret-key-change-it"

	// 服务器配置
	ServerPort = ":8080"
)

// GetDSN 获取数据库连接字符串
func GetDSN() string {
	return DBUser + ":" + DBPassword + "@tcp(" + DBHost + ":" + DBPort + ")/" + DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
}
