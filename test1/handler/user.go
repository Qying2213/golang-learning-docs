package handler

import (
	"strconv"
	"test1/model"
	"test1/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserHandler 用户处理器
type UserHandler struct {
	DB *gorm.DB
}

// NewUserHandler 创建用户处理器
func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{DB: db}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=2,max=20"`
	Password string `json:"password" binding:"required,min=6,max=20"`
	Email    string `json:"email" binding:"required,email"`
}

// Register 用户注册
// TODO: 任务4 - 实现注册功能
// 要求：
//  1. 绑定并验证参数
//  2. 检查用户名是否已存在
//  3. 密码加密（使用 bcrypt）
//  4. 创建用户
//  5. 返回用户信息（不含密码）
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.CodeParamError, "参数错误："+err.Error())
		return
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

// Login 用户登录
// TODO: 任务4 - 实现登录功能
// 要求：
//  1. 绑定并验证参数
//  2. 查询用户
//  3. 验证密码（使用 bcrypt.CompareHashAndPassword）
//  4. 生成 JWT Token
//  5. 返回 Token 和用户信息
func (h *UserHandler) Login(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.CodeParamError, "参数错误")
		return
	}
	var existUser model.User
	err := h.DB.Where("username = ?", req.Username).First(&existUser).Error
	if err != nil {
		utils.Error(c, utils.CodeUserExist, "用户名已存在")
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
	}
	h.DB.Create(&user)

	utils.Success(c, user) // 改成你的实现
}

// GetProfile 获取当前用户信息
// TODO: 任务4 - 实现获取用户信息
// 要求：
//  1. 从 Context 获取 userID（JWT 中间件已存入）
//  2. 查询用户信息
//  3. 返回用户信息
func (h *UserHandler) GetProfile(c *gin.Context) {
	// TODO: 从 Context 获取 userID
	// 提示：userID := c.GetUint("userID")
	userID := c.GetUint("userID")
	var user model.User
	err := h.DB.Where("id = ?", userID).First(&user).Error
	if err != nil {
		utils.Error(c, utils.CodeNotFound, "用户不存在")
		return
	}
	utils.Success(c, user) // 改成你的实现
}

// UpdateProfileRequest 更新用户信息请求
type UpdateProfileRequest struct {
	Email string `json:"email" binding:"omitempty,email"`
}

// UpdateProfile 更新用户信息
// TODO: 任务4 - 实现更新用户信息
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// TODO: 实现
	userID:=c.GetUint("userID")
	var req UpdatePasswordRequest
	if err:=c.ShouldBindJSON(&req);err!=nil{
		utils.Error(c,utils.CodeParamError,"更新失败")
		return
	}

	var user model.User
	h.DB.Where("id = ?",userID).First(&user)
	utils.Success(c, user)
}

// UpdatePasswordRequest 修改密码请求
type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=20"`
}

// UpdatePassword 修改密码
// TODO: 任务4 - 实现修改密码
// 要求：
//  1. 验证旧密码
//  2. 加密新密码
//  3. 更新密码
func (h *UserHandler) UpdatePassword(c *gin.Context) {
	// TODO: 实现
	userID:=c.GetUint("userID")
	var req UpdatePasswordRequest
	if err:=c.ShouldBindJSON(&req);err!=nil{
		utils.Error(c,utils.CodeParamError,"更新失败")
		return
	}
	var user model.User
	if err:=h.DB.Where("id = ?",userID).First(&user).Error;err!=nil{
		utils.Error(c,utils.CodeNotFound,"用户不存在")
		return
	}
	if err:=bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(req.OldPassword));err!=nil{
		utils.Error(c,utils.CodePasswordError,"旧密码错误")
		return
	}
	hashedPassword,_:=bcrypt.GenerateFromPassword([]byte(req.NewPassword),bcrypt.DefaultCost)
	h.DB.Model(&user).Update("password",string(hashedPassword))
	utils.SuccessWithMessage(c, "密码修改成功", nil)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
    // 1. 获取分页参数
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

    // 2. 查询用户列表
    var users []model.User
    h.DB.Offset((page - 1) * limit).Limit(limit).Find(&users)

    // 3. 查询总数
    var total int64
    h.DB.Model(&model.User{}).Count(&total)

    // 4. 返回响应
    utils.Success(c, gin.H{
        "list":  users,
        "total": total,
        "page":  page,
        "limit": limit,
    })
}
