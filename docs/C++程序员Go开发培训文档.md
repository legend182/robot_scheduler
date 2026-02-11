# Go Web开发文档

## Gin + GORM 实战指南

---

## 目录

1. [Go语言快速入门（针对C++程序员）](#1-go语言快速入门针对c程序员)
2. [Gin Web框架详解](#2-gin-web框架详解)
3. [GORM ORM框架详解](#3-gorm-orm框架详解)
4. [Web开发流程和最佳实践](#4-web开发流程和最佳实践)
5. [常见问题和注意事项](#5-常见问题和注意事项)
6. [实战案例](#6-实战案例)

---

## 1. Go语言快速入门（针对C++程序员）

### 1.1 与C++的主要区别

| 特性   | C++            | Go             |
| ---- | -------------- | -------------- |
| 内存管理 | 手动（new/delete） | 自动GC           |
| 指针   | 支持，可空指针        | 支持，但更安全        |
| 继承   | 支持（类继承）        | 不支持（组合代替）      |
| 泛型   | 模板             | Go 1.18+支持     |
| 异常   | try/catch      | error返回值       |
| 并发   | std::thread    | goroutine（更轻量） |

### 1.2 关键语法速览

#### 变量声明

```go
// C++风格（不推荐）
var x int = 10
var s string = "hello"

// Go风格（推荐）
x := 10              // 类型推断
s := "hello"
var x int            // 零值初始化
```

#### 指针（比C++更安全）

```go
// Go指针不能进行算术运算
var x int = 10
var p *int = &x      // 获取地址
*p = 20              // 解引用

// 没有空指针异常，nil检查是必须的
var p *int = nil
if p != nil {
    *p = 10
}
```

#### 错误处理（Go的核心特性）

```go
// C++: try-catch
try {
    result = doSomething();
} catch (Exception& e) {
    // handle error
}

// Go: 错误作为返回值
result, err := doSomething()
if err != nil {
    // handle error
    return err
}
```

#### 结构体和方法

```go
// 定义结构体（类似C++的struct，但没有class）
type User struct {
    ID   int
    Name string
}

// 方法（不是成员函数，是接收者）
func (u *User) GetName() string {
    return u.Name
}

// 调用
user := &User{ID: 1, Name: "Alice"}
name := user.GetName()
```

#### 接口（Go的核心抽象）

```go
// 接口定义（类似C++的纯虚类）
type Writer interface {
    Write([]byte) (int, error)
}

// 实现接口（隐式实现，不需要显式声明）
type FileWriter struct {}
func (f *FileWriter) Write(data []byte) (int, error) {
    // implementation
    return len(data), nil
}

// 使用
var w Writer = &FileWriter{}
w.Write([]byte("hello"))
```

#### 并发（goroutine vs std::thread）

```go
// C++: std::thread
std::thread t([](){
    // do work
});
t.join();

// Go: goroutine（更轻量，开销更小）
go func() {
    // do work
}()

// 同步（channel vs mutex）
ch := make(chan int)
go func() {
    ch <- 42  // 发送
}()
value := <-ch  // 接收
```

### 1.3 包管理

```go
// 导入包
import (
    "fmt"                    // 标准库
    "github.com/gin-gonic/gin"  // 第三方库
    "robot_scheduler/internal/model"  // 本地包
)

// 包名约定：小写，简短
package handler

// 导出：首字母大写=public，小写=private
type UserHandler struct {  // 可导出
    userService *service.UserService  // 可导出
    config      *config.Config        // 可导出
}

type internalState int  // 不可导出
```

---

## 2. Gin Web框架详解

### 2.1 Gin简介

Gin是Go语言最流行的Web框架之一，类似C++中的：

- **Crow**（轻量级）
- **Drogon**（高性能）
- **Beast**（底层）

**特点**：

- 高性能（基于httprouter）
- 中间件支持
- JSON验证
- 路由组
- 错误管理

### 2.2 基本使用

#### 2.2.1 创建Gin应用

```go
package main

import "github.com/gin-gonic/gin"

func main() {
    // 创建Gin引擎
    r := gin.Default()  // 包含Logger和Recovery中间件

    // 或
    r := gin.New()      // 不包含默认中间件

    // 定义路由
    r.GET("/hello", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "hello",
        })
    })

    // 启动服务器
    r.Run(":8080")
}
```

**对比C++（使用Crow）**：

```cpp
// C++
crow::SimpleApp app;
CROW_ROUTE(app, "/hello")([](){
    return "hello";
});
app.port(8080).run();
```

#### 2.2.2 路由定义

```go
// GET请求
r.GET("/users/:id", getUserHandler)

// POST请求
r.POST("/users", createUserHandler)

// PUT请求
r.PUT("/users/:id", updateUserHandler)

// DELETE请求
r.DELETE("/users/:id", deleteUserHandler)

// 路径参数
r.GET("/users/:id/posts/:postId", func(c *gin.Context) {
    id := c.Param("id")           // 获取路径参数
    postId := c.Param("postId")
})

// 查询参数
r.GET("/search", func(c *gin.Context) {
    keyword := c.Query("keyword")        // ?keyword=test
    page := c.DefaultQuery("page", "1")  // 带默认值
})
```

**对比C++（使用Drogon）**：

```cpp
// C++
app.registerHandler("/users/{id}",
    [](const HttpRequestPtr& req, std::function<void(const HttpResponsePtr&)> callback, int id) {
        // handle
    });
```

#### 2.2.3 请求处理

```go
// Handler函数签名
func handler(c *gin.Context) {
    // c.Request: *http.Request（原始请求）
    // c.Writer: http.ResponseWriter（响应写入器）

    // 获取请求体（JSON）
    var req UserCreateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // 获取表单数据
    name := c.PostForm("name")

    // 获取文件
    file, err := c.FormFile("file")

    // 返回JSON响应
    c.JSON(200, gin.H{
        "code": 0,
        "data": result,
    })

    // 返回XML
    c.XML(200, data)

    // 返回HTML
    c.HTML(200, "index.html", gin.H{"title": "Home"})

    // 重定向
    c.Redirect(302, "/login")

    // 设置Header
    c.Header("X-Custom-Header", "value")

    // 设置Cookie
    c.SetCookie("token", "abc123", 3600, "/", "localhost", false, true)
}
```

### 2.3 中间件（Middleware）

中间件类似C++中的**拦截器**或**过滤器**，在请求处理前后执行。

#### 2.3.1 中间件工作原理

```
请求 → 中间件1 → 中间件2 → Handler → 中间件2 → 中间件1 → 响应
```

```go
// 中间件函数签名
func middleware(c *gin.Context) {
    // 请求前处理
    start := time.Now()

    // 调用下一个处理器
    c.Next()

    // 请求后处理
    latency := time.Since(start)
    logger.Info("request completed", zap.Duration("latency", latency))
}

// 提前终止（不调用c.Next()）
func authMiddleware(c *gin.Context) {
    token := c.GetHeader("Authorization")
    if token == "" {
        c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
        return  // 不调用c.Next()，终止请求链
    }
    c.Next()
}
```

#### 2.3.2 内置中间件

```go
// Logger中间件（记录请求日志）
r.Use(gin.Logger())

// Recovery中间件（捕获panic）
r.Use(gin.Recovery())

// CORS中间件（跨域支持）
r.Use(func(c *gin.Context) {
    c.Header("Access-Control-Allow-Origin", "*")
    c.Next()
})
```

#### 2.3.3 自定义中间件示例

```go
// 认证中间件
func JWTAuth(jwtSecret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 从Header获取token
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(401, gin.H{"error": "missing token"})
            return
        }

        // 验证token
        token := strings.TrimPrefix(authHeader, "Bearer ")
        claims, err := validateToken(token, jwtSecret)
        if err != nil {
            c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
            return
        }

        // 将用户信息存入Context（类似C++的request attributes）
        c.Set("user_id", claims.UserID)
        c.Set("user_name", claims.UserName)

        c.Next()
    }
}

// 使用中间件
r.Use(JWTAuth("secret-key"))
// 或特定路由组
authGroup := r.Group("/api")
authGroup.Use(JWTAuth("secret-key"))
```

#### 2.3.4 中间件链

```go
// 全局中间件
r.Use(middleware1)
r.Use(middleware2)

// 路由组中间件
api := r.Group("/api")
api.Use(authMiddleware)
{
    api.GET("/users", getUserHandler)  // 会经过middleware1, middleware2, authMiddleware
}

// 单个路由中间件
r.GET("/public", middleware3, publicHandler)
```

### 2.4 路由组（Route Groups）

路由组用于组织相关路由，类似C++中的**命名空间**。

```go
// 创建路由组
api := r.Group("/api/v1")
{
    // /api/v1/users
    api.GET("/users", listUsers)
    api.POST("/users", createUser)

    // 嵌套路由组
    admin := api.Group("/admin")
    admin.Use(adminMiddleware)  // 仅admin组使用
    {
        // /api/v1/admin/users
        admin.GET("/users", adminListUsers)
    }
}
```

### 2.5 参数绑定和验证

#### 2.5.1 JSON绑定

```go
type UserCreateRequest struct {
    UserName string `json:"user_name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Age      int    `json:"age" binding:"min=18,max=100"`
}

func createUser(c *gin.Context) {
    var req UserCreateRequest
    // ShouldBindJSON: 绑定失败返回error，不中断程序
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // BindJSON: 绑定失败自动返回400，中断程序
    // c.BindJSON(&req)  // 不推荐，错误处理不灵活
}
```

#### 2.5.2 查询参数绑定

```go
type Pagination struct {
    Page     int `form:"page" binding:"min=1"`
    PageSize int `form:"page_size" binding:"min=1,max=100"`
}

func listUsers(c *gin.Context) {
    var pagination Pagination
    if err := c.ShouldBindQuery(&pagination); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    // 使用 pagination.Page, pagination.PageSize
}
```

#### 2.5.3 路径参数

```go
r.GET("/users/:id", func(c *gin.Context) {
    id := c.Param("id")  // 字符串类型

    // 转换为int
    userID, err := strconv.Atoi(id)
    if err != nil {
        c.JSON(400, gin.H{"error": "invalid id"})
        return
    }
})
```

### 2.6 Context（上下文）

`gin.Context`是请求的上下文对象，类似C++中的**HttpRequest**对象。

```go
// 存储数据（请求生命周期内有效）
c.Set("user_id", 123)
c.Set("request_id", uuid.New())

// 获取数据
userID, exists := c.Get("user_id")
if exists {
    id := userID.(int)  // 类型断言
}

// 获取原始请求对象
req := c.Request
method := req.Method
url := req.URL

// 获取客户端IP
ip := c.ClientIP()

// 获取请求头
userAgent := c.GetHeader("User-Agent")
```

### 2.7 错误处理

```go
// 统一错误响应
func ErrorResponse(c *gin.Context, code int, message string) {
    c.JSON(code, gin.H{
        "code":    code,
        "message": message,
        "data":    nil,
    })
}

// 使用
func handler(c *gin.Context) {
    if err := doSomething(); err != nil {
        ErrorResponse(c, 500, err.Error())
        return
    }
}

// 中止请求（类似C++的return early）
c.Abort()                    // 中止，不发送响应
c.AbortWithStatus(404)       // 中止并发送状态码
c.AbortWithStatusJSON(400, gin.H{"error": "bad request"})  // 中止并发送JSON
```

---

## 3. GORM ORM框架详解

### 3.1 GORM简介

GORM是Go语言最流行的ORM框架，类似C++中的：

- **ODB**（对象关系映射）
- **SOCI**（数据库抽象层）

**特点**：

- 自动迁移
- 关联关系
- 钩子（Hooks）
- 预加载（Preload）
- 事务支持

### 3.2 数据库连接

```go
import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

// PostgreSQL连接
dsn := "host=localhost user=postgres password=postgres dbname=mydb port=5432 sslmode=disable"
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

// SQLite连接
db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

// 连接池配置
sqlDB, _ := db.DB()
sqlDB.SetMaxIdleConns(10)           // 最大空闲连接
sqlDB.SetMaxOpenConns(100)          // 最大打开连接
sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大生存时间
```

### 3.3 定义模型（Entity）

```go
import "gorm.io/gorm"

// 基础模型（包含ID, CreatedAt, UpdatedAt, DeletedAt）
type User struct {
    gorm.Model              // 嵌入gorm.Model
    UserName  string        `gorm:"type:varchar(100);not null;uniqueIndex"`
    Email     string        `gorm:"type:varchar(255);uniqueIndex"`
    Age       int           `gorm:"default:0"`
    IsActive  bool          `gorm:"default:true"`
    Profile   UserProfile   `gorm:"foreignKey:UserID"`  // 一对一
    Posts     []Post        `gorm:"foreignKey:UserID"`   // 一对多
}

// 自定义表名
func (User) TableName() string {
    return "user_info"
}

// 一对一关联
type UserProfile struct {
    gorm.Model
    UserID   uint
    User     User    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
    Bio      string
    Avatar   string
}

// 一对多关联
type Post struct {
    gorm.Model
    UserID  uint
    User    User    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
    Title   string
    Content string
}
```

**GORM标签说明**：

- `type:varchar(100)`: 数据库类型
- `not null`: 非空约束
- `uniqueIndex`: 唯一索引
- `default:0`: 默认值
- `foreignKey:UserID`: 外键
- `constraint:OnDelete:CASCADE`: 级联删除

### 3.4 CRUD操作

#### 3.4.1 Create（创建）

```go
// 创建单个记录
user := User{UserName: "alice", Email: "alice@example.com"}
result := db.Create(&user)
if result.Error != nil {
    // 处理错误
}
// user.ID 自动填充

// 批量创建
users := []User{
    {UserName: "bob", Email: "bob@example.com"},
    {UserName: "charlie", Email: "charlie@example.com"},
}
db.Create(&users)

// 指定字段创建（忽略其他字段）
db.Select("UserName", "Email").Create(&user)

// 忽略字段创建
db.Omit("Age").Create(&user)
```

#### 3.4.2 Read（查询）

```go
// 查询单个记录（主键）
var user User
db.First(&user, 1)  // SELECT * FROM users WHERE id = 1 LIMIT 1

// 查询单个记录（条件）
db.Where("user_name = ?", "alice").First(&user)

// 查询所有记录
var users []User
db.Find(&users)

// 条件查询
db.Where("age > ?", 18).Find(&users)
db.Where("user_name IN ?", []string{"alice", "bob"}).Find(&users)

// 排序
db.Order("created_at DESC").Find(&users)

// 限制数量
db.Limit(10).Find(&users)

// 分页
db.Offset(10).Limit(10).Find(&users)

// 计数
var count int64
db.Model(&User{}).Where("age > ?", 18).Count(&count)

// 检查记录是否存在
var user User
result := db.Where("email = ?", "alice@example.com").First(&user)
if errors.Is(result.Error, gorm.ErrRecordNotFound) {
    // 记录不存在
}

// 选择特定字段
db.Select("user_name", "email").Find(&users)

// 原生SQL
db.Raw("SELECT * FROM users WHERE age > ?", 18).Scan(&users)
```

#### 3.4.3 Update（更新）

```go
// 更新单个字段
db.Model(&user).Update("age", 25)

// 更新多个字段
db.Model(&user).Updates(User{Age: 25, IsActive: true})

// 使用map更新（可以更新零值）
db.Model(&user).Updates(map[string]interface{}{
    "age": 25,
    "is_active": false,
})

// 条件更新
db.Model(&User{}).Where("age < ?", 18).Update("is_active", false)

// 保存整个对象（更新所有字段）
db.Save(&user)

// 更新指定字段
db.Model(&user).Select("age").Updates(User{Age: 25, Email: "new@example.com"})
// 只更新age，忽略email
```

#### 3.4.4 Delete（删除）

```go
// 软删除（GORM默认，设置DeletedAt字段）
db.Delete(&user)

// 硬删除（物理删除）
db.Unscoped().Delete(&user)

// 条件删除
db.Where("age < ?", 18).Delete(&User{})

// 批量删除
db.Delete(&User{}, []int{1, 2, 3})
```

### 3.5 关联查询（Relations）

#### 3.5.1 预加载（Preload）

```go
// 预加载一对一关联
var user User
db.Preload("Profile").First(&user, 1)
// 执行两条SQL：
// SELECT * FROM users WHERE id = 1
// SELECT * FROM user_profiles WHERE user_id = 1

// 预加载一对多关联
db.Preload("Posts").First(&user, 1)

// 嵌套预加载
db.Preload("Posts.Comments").First(&user, 1)

// 条件预加载
db.Preload("Posts", "status = ?", "published").First(&user, 1)
```

#### 3.5.2 关联操作

```go
// 创建关联
user := User{UserName: "alice"}
profile := UserProfile{Bio: "Developer"}
db.Create(&user)
db.Model(&user).Association("Profile").Append(&profile)

// 替换关联
db.Model(&user).Association("Posts").Replace(&posts)

// 删除关联
db.Model(&user).Association("Posts").Delete(&post)

// 清空关联
db.Model(&user).Association("Posts").Clear()

// 计数
count := db.Model(&user).Association("Posts").Count()
```

### 3.6 事务（Transactions）

```go
// 自动事务
err := db.Transaction(func(tx *gorm.DB) error {
    // 创建用户
    if err := tx.Create(&user).Error; err != nil {
        return err  // 回滚
    }

    // 创建用户资料
    if err := tx.Create(&profile).Error; err != nil {
        return err  // 回滚
    }

    return nil  // 提交
})

// 手动事务
tx := db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

if err := tx.Create(&user).Error; err != nil {
    tx.Rollback()
    return err
}

if err := tx.Create(&profile).Error; err != nil {
    tx.Rollback()
    return err
}

tx.Commit()
```

### 3.7 钩子（Hooks）

GORM提供生命周期钩子，类似C++的**事件回调**。

```go
type User struct {
    gorm.Model
    UserName string
    Password string
}

// BeforeCreate: 创建前
func (u *User) BeforeCreate(tx *gorm.DB) error {
    // 加密密码
    hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    u.Password = string(hashed)
    return nil
}

// AfterCreate: 创建后
func (u *User) AfterCreate(tx *gorm.DB) error {
    // 发送欢迎邮件
    sendWelcomeEmail(u.Email)
    return nil
}

// BeforeUpdate: 更新前
func (u *User) BeforeUpdate(tx *gorm.DB) error {
    // 记录更新时间
    u.UpdatedAt = time.Now()
    return nil
}
```

### 3.8 作用域（Scopes）

作用域用于封装常用查询逻辑。

```go
// 定义作用域
func ActiveUsers(db *gorm.DB) *gorm.DB {
    return db.Where("is_active = ?", true)
}

func OlderThan(age int) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("age > ?", age)
    }
}

// 使用作用域
db.Scopes(ActiveUsers, OlderThan(18)).Find(&users)
```

### 3.9 原生SQL

```go
// 执行原生SQL
db.Exec("UPDATE users SET age = ? WHERE id = ?", 25, 1)

// 查询
type Result struct {
    UserName string
    PostCount int
}
var results []Result
db.Raw(`
    SELECT u.user_name, COUNT(p.id) as post_count
    FROM users u
    LEFT JOIN posts p ON u.id = p.user_id
    GROUP BY u.id
`).Scan(&results)
```

---

## 4. Web开发流程和最佳实践

### 4.1 标准开发流程

#### 4.1.1 项目结构

```
project/
├── cmd/              # 应用入口
│   └── main.go
├── internal/         # 内部代码（不对外暴露）
│   ├── api/          # API层
│   │   ├── handler/  # 请求处理器
│   │   ├── middleware/
│   │   └── router/
│   ├── service/      # 业务逻辑层
│   ├── dao/          # 数据访问层
│   │   └── interfaces/  # DAO接口
│   ├── model/        # 数据模型
│   │   ├── entity/   # 数据库实体
│   │   └── dto/      # 数据传输对象
│   └── config/       # 配置
├── configs/          # 配置文件
├── docs/             # 文档
└── go.mod            # 依赖管理
```

#### 4.1.2 开发步骤

**步骤1: 定义数据模型（Entity）**

```go
// internal/model/entity/user.go
package entity

import "gorm.io/gorm"

type User struct {
    gorm.Model
    UserName string `gorm:"type:varchar(100);not null;uniqueIndex"`
    Email    string `gorm:"type:varchar(255);uniqueIndex"`
}
```

**步骤2: 定义DTO（数据传输对象）**

```go
// internal/model/dto/user.go
package dto

type UserCreateRequest struct {
    UserName string `json:"user_name" binding:"required,min=3"`
    Email    string `json:"email" binding:"required,email"`
}

type UserResponse struct {
    ID       uint   `json:"id"`
    UserName string `json:"user_name"`
    Email    string `json:"email"`
}
```

**步骤3: 定义DAO接口和实现**

```go
// internal/dao/interfaces/user.go
package interfaces

import (
    "context"
    "robot_scheduler/internal/model/entity"
)

type UserDAO interface {
    Create(ctx context.Context, user *entity.User) error
    FindByID(ctx context.Context, id uint) (*entity.User, error)
    FindAll(ctx context.Context) ([]*entity.User, error)
}

// internal/dao/user.go
package impl

import (
    "context"
    dao "robot_scheduler/internal/dao/interfaces"
    "robot_scheduler/internal/model/entity"
    "gorm.io/gorm"
)

type UserDAOImpl struct {
    db *gorm.DB
}

func NewUserDAO(db *gorm.DB) dao.UserDAO {
    return &UserDAOImpl{db: db}
}

func (d *UserDAOImpl) Create(ctx context.Context, user *entity.User) error {
    return d.db.WithContext(ctx).Create(user).Error
}

func (d *UserDAOImpl) FindByID(ctx context.Context, id uint) (*entity.User, error) {
    var user entity.User
    err := d.db.WithContext(ctx).First(&user, id).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}
```

**步骤4: 实现Service层**

```go
// internal/service/user.go
package service

import (
    "context"
    dao "robot_scheduler/internal/dao/interfaces"
    "robot_scheduler/internal/model/dto"
    "robot_scheduler/internal/model/entity"
)

type UserService struct {
    userDAO dao.UserDAO
}

func NewUserService(userDAO dao.UserDAO) *UserService {
    return &UserService{userDAO: userDAO}
}

func (s *UserService) CreateUser(ctx context.Context, req *dto.UserCreateRequest) (*entity.User, error) {
    user := &entity.User{
        UserName: req.UserName,
        Email:    req.Email,
    }
    if err := s.userDAO.Create(ctx, user); err != nil {
        return nil, err
    }
    return user, nil
}
```

**步骤5: 实现Handler层**

```go
// internal/api/handler/user.go
package handler

import (
    "robot_scheduler/internal/model/dto"
    "robot_scheduler/internal/service"
    "github.com/gin-gonic/gin"
)

type UserHandler struct {
    userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
    return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
    var req dto.UserCreateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        BadRequest(c, "无效的请求参数: "+err.Error())
        return
    }

    user, err := h.userService.CreateUser(c.Request.Context(), &req)
    if err != nil {
        InternalServerError(c, "创建用户失败: "+err.Error())
        return
    }

    Success(c, user)
}
```

**步骤6: 注册路由**

```go
// internal/api/router/router.go
package router

import (
    "robot_scheduler/internal/api/handler"
    impl "robot_scheduler/internal/dao"
    "robot_scheduler/internal/database"
    "robot_scheduler/internal/service"
    "github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
    db := database.DB

    // 初始化各层
    userDAO := impl.NewUserDAO(db)
    userService := service.NewUserService(userDAO)
    userHandler := handler.NewUserHandler(userService)

    // 注册路由
    api := r.Group("/api/v1")
    {
        api.POST("/users", userHandler.CreateUser)
        api.GET("/users/:id", userHandler.GetUser)
    }
}
```

### 4.2 依赖注入模式

```go
// 构造函数注入（推荐）
func NewUserService(userDAO dao.UserDAO) *UserService {
    return &UserService{userDAO: userDAO}
}

// 使用
userDAO := impl.NewUserDAO(db)
userService := service.NewUserService(userDAO)
userHandler := handler.NewUserHandler(userService)
```

### 4.3 错误处理最佳实践

```go
// 定义业务错误
var (
    ErrUserNotFound = errors.New("用户不存在")
    ErrInvalidPassword = errors.New("密码错误")
)

// Service层返回业务错误
func (s *UserService) GetUser(ctx context.Context, id uint) (*entity.User, error) {
    user, err := s.userDAO.FindByID(ctx, id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrUserNotFound
        }
        return nil, err
    }
    return user, nil
}

// Handler层转换错误
func (h *UserHandler) GetUser(c *gin.Context) {
    user, err := h.userService.GetUser(c.Request.Context(), id)
    if err != nil {
        switch err {
        case service.ErrUserNotFound:
            NotFound(c, "用户不存在")
        default:
            InternalServerError(c, "获取用户失败")
        }
        return
    }
    Success(c, user)
}
```

### 4.4 Context使用

```go
// Context用于传递请求上下文信息（超时、取消、值传递）
func (s *UserService) GetUser(ctx context.Context, id uint) (*entity.User, error) {
    // 检查超时
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
    }

    // 传递Context给DAO层
    return s.userDAO.FindByID(ctx, id)
}

// 设置超时
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

user, err := userService.GetUser(ctx, 1)
```

---

## 5. 常见问题和注意事项

### 5.1 Gin相关

#### 问题1: 中间件执行顺序

```go
// 中间件按注册顺序执行
r.Use(middleware1)  // 1. 最先执行
r.Use(middleware2)  // 2. 其次执行
r.GET("/hello", handler)  // 3. 最后执行

// 执行顺序：middleware1 → middleware2 → handler → middleware2 → middleware1
```

#### 问题2: 路由冲突

```go
// ❌ 错误：路由冲突
r.GET("/users/:id", handler1)
r.GET("/users/new", handler2)  // 冲突！"new"会被当作id

// ✅ 正确：具体路由在前
r.GET("/users/new", handler2)
r.GET("/users/:id", handler1)
```

#### 问题3: 响应多次写入

```go
// ❌ 错误：多次写入响应
c.JSON(200, data1)
c.JSON(200, data2)  // 错误！

// ✅ 正确：只写入一次
c.JSON(200, data)
```

### 5.2 GORM相关

#### 问题1: 零值更新

```go
// ❌ 问题：零值不会被更新
user := User{Age: 0, IsActive: false}
db.Model(&user).Updates(user)  // Age和IsActive不会被更新

// ✅ 解决：使用map或Select
db.Model(&user).Updates(map[string]interface{}{
    "age": 0,
    "is_active": false,
})
```

#### 问题2: N+1查询问题

```go
// ❌ 问题：N+1查询
var users []User
db.Find(&users)  // 1次查询
for _, user := range users {
    db.Model(&user).Association("Posts").Find(&user.Posts)  // N次查询
}

// ✅ 解决：使用Preload
var users []User
db.Preload("Posts").Find(&users)  // 2次查询
```

#### 问题3: 事务中的错误处理

```go
// ❌ 错误：事务中panic不会回滚
db.Transaction(func(tx *gorm.DB) error {
    tx.Create(&user)
    panic("error")  // 不会回滚！
    return nil
})

// ✅ 正确：返回error
db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&user).Error; err != nil {
        return err  // 会回滚
    }
    return nil
})
```

### 5.3 性能优化

#### 优化1: 数据库连接池

```go
sqlDB, _ := db.DB()
sqlDB.SetMaxIdleConns(10)           // 根据实际情况调整
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

#### 优化2: 批量操作

```go
// ❌ 慢：循环插入
for _, user := range users {
    db.Create(&user)
}

// ✅ 快：批量插入
db.Create(&users)
```

#### 优化3: 索引优化

```go
// 为常用查询字段添加索引
type User struct {
    gorm.Model
    Email string `gorm:"index"`                    // 单列索引
    Name  string `gorm:"index:idx_name_email"`    // 复合索引
}
```

### 5.4 安全注意事项

#### 1. SQL注入防护

```go
// ✅ GORM自动防护SQL注入
db.Where("user_name = ?", userName).First(&user)

// ❌ 不要使用字符串拼接
db.Where("user_name = '" + userName + "'")  // 危险！
```

#### 2. 密码加密

```go
// ❌ 错误：明文存储
user.Password = password

// ✅ 正确：使用bcrypt
hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
user.Password = string(hashed)
```

#### 3. 参数验证

```go
// ✅ 使用binding标签验证
type UserCreateRequest struct {
    UserName string `json:"user_name" binding:"required,min=3,max=20"`
    Email    string `json:"email" binding:"required,email"`
}
```

---

## 6. 实战案例

### 6.1 完整的CRUD示例

假设我们要实现一个**文章管理**功能：

#### 步骤1: 定义Entity

```go
// internal/model/entity/article.go
package entity

import "gorm.io/gorm"

type Article struct {
    gorm.Model
    Title   string `gorm:"type:varchar(255);not null"`
    Content string `gorm:"type:text"`
    UserID  uint
    User    User   `gorm:"foreignKey:UserID"`
}
```

#### 步骤2: 定义DTO

```go
// internal/model/dto/article.go
package dto

type ArticleCreateRequest struct {
    Title   string `json:"title" binding:"required,min=5,max=255"`
    Content string `json:"content" binding:"required"`
}

type ArticleResponse struct {
    ID        uint   `json:"id"`
    Title     string `json:"title"`
    Content   string `json:"content"`
    UserID    uint   `json:"user_id"`
    CreatedAt string `json:"created_at"`
}
```

#### 步骤3: 实现DAO

```go
// internal/dao/interfaces/article.go
package interfaces

import (
    "context"
    "robot_scheduler/internal/model/entity"
)

type ArticleDAO interface {
    Create(ctx context.Context, article *entity.Article) error
    FindByID(ctx context.Context, id uint) (*entity.Article, error)
    FindAll(ctx context.Context, offset, limit int) ([]*entity.Article, int64, error)
    Update(ctx context.Context, article *entity.Article) error
    Delete(ctx context.Context, id uint) error
}

// internal/dao/article.go
package impl

import (
    "context"
    dao "robot_scheduler/internal/dao/interfaces"
    "robot_scheduler/internal/model/entity"
    "gorm.io/gorm"
)

type ArticleDAOImpl struct {
    db *gorm.DB
}

func NewArticleDAO(db *gorm.DB) dao.ArticleDAO {
    return &ArticleDAOImpl{db: db}
}

func (d *ArticleDAOImpl) Create(ctx context.Context, article *entity.Article) error {
    return d.db.WithContext(ctx).Create(article).Error
}

func (d *ArticleDAOImpl) FindByID(ctx context.Context, id uint) (*entity.Article, error) {
    var article entity.Article
    err := d.db.WithContext(ctx).Preload("User").First(&article, id).Error
    if err != nil {
        return nil, err
    }
    return &article, nil
}

func (d *ArticleDAOImpl) FindAll(ctx context.Context, offset, limit int) ([]*entity.Article, int64, error) {
    var articles []*entity.Article
    var total int64

    query := d.db.WithContext(ctx).Model(&entity.Article{})
    query.Count(&total)

    err := query.Offset(offset).Limit(limit).Preload("User").Find(&articles).Error
    return articles, total, err
}

func (d *ArticleDAOImpl) Update(ctx context.Context, article *entity.Article) error {
    return d.db.WithContext(ctx).Save(article).Error
}

func (d *ArticleDAOImpl) Delete(ctx context.Context, id uint) error {
    return d.db.WithContext(ctx).Delete(&entity.Article{}, id).Error
}
```

#### 步骤4: 实现Service

```go
// internal/service/article.go
package service

import (
    "context"
    dao "robot_scheduler/internal/dao/interfaces"
    "robot_scheduler/internal/model/dto"
    "robot_scheduler/internal/model/entity"
    "errors"
)

type ArticleService struct {
    articleDAO dao.ArticleDAO
}

func NewArticleService(articleDAO dao.ArticleDAO) *ArticleService {
    return &ArticleService{articleDAO: articleDAO}
}

func (s *ArticleService) CreateArticle(ctx context.Context, userID uint, req *dto.ArticleCreateRequest) (*entity.Article, error) {
    article := &entity.Article{
        Title:   req.Title,
        Content: req.Content,
        UserID:  userID,
    }
    if err := s.articleDAO.Create(ctx, article); err != nil {
        return nil, err
    }
    return article, nil
}

func (s *ArticleService) GetArticle(ctx context.Context, id uint) (*entity.Article, error) {
    article, err := s.articleDAO.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    if article == nil {
        return nil, errors.New("文章不存在")
    }
    return article, nil
}

func (s *ArticleService) ListArticles(ctx context.Context, page, pageSize int) ([]*entity.Article, int64, error) {
    offset := (page - 1) * pageSize
    articles, total, err := s.articleDAO.FindAll(ctx, offset, pageSize)
    return articles, total, err
}

func (s *ArticleService) UpdateArticle(ctx context.Context, id uint, req *dto.ArticleUpdateRequest) error {
    article, err := s.articleDAO.FindByID(ctx, id)
    if err != nil {
        return err
    }
    if article == nil {
        return errors.New("文章不存在")
    }

    article.Title = req.Title
    article.Content = req.Content
    return s.articleDAO.Update(ctx, article)
}

func (s *ArticleService) DeleteArticle(ctx context.Context, id uint) error {
    return s.articleDAO.Delete(ctx, id)
}
```

#### 步骤5: 实现Handler

```go
// internal/api/handler/article.go
package handler

import (
    "strconv"
    "robot_scheduler/internal/model/dto"
    "robot_scheduler/internal/service"
    "github.com/gin-gonic/gin"
)

type ArticleHandler struct {
    articleService *service.ArticleService
}

func NewArticleHandler(articleService *service.ArticleService) *ArticleHandler {
    return &ArticleHandler{articleService: articleService}
}

func (h *ArticleHandler) CreateArticle(c *gin.Context) {
    // 获取用户ID（从中间件设置）
    userID, _ := c.Get("user_id")

    var req dto.ArticleCreateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        BadRequest(c, "无效的请求参数: "+err.Error())
        return
    }

    article, err := h.articleService.CreateArticle(c.Request.Context(), userID.(uint), &req)
    if err != nil {
        InternalServerError(c, "创建文章失败: "+err.Error())
        return
    }

    Success(c, article)
}

func (h *ArticleHandler) GetArticle(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        BadRequest(c, "无效的文章ID")
        return
    }

    article, err := h.articleService.GetArticle(c.Request.Context(), uint(id))
    if err != nil {
        NotFound(c, "文章不存在")
        return
    }

    Success(c, article)
}

func (h *ArticleHandler) ListArticles(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

    articles, total, err := h.articleService.ListArticles(c.Request.Context(), page, pageSize)
    if err != nil {
        InternalServerError(c, "获取文章列表失败: "+err.Error())
        return
    }

    Success(c, gin.H{
        "items": articles,
        "total": total,
        "page":  page,
        "page_size": pageSize,
    })
}

func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        BadRequest(c, "无效的文章ID")
        return
    }

    var req dto.ArticleUpdateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        BadRequest(c, "无效的请求参数: "+err.Error())
        return
    }

    if err := h.articleService.UpdateArticle(c.Request.Context(), uint(id), &req); err != nil {
        InternalServerError(c, "更新文章失败: "+err.Error())
        return
    }

    Success(c, gin.H{"message": "更新成功"})
}

func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        BadRequest(c, "无效的文章ID")
        return
    }

    if err := h.articleService.DeleteArticle(c.Request.Context(), uint(id)); err != nil {
        InternalServerError(c, "删除文章失败: "+err.Error())
        return
    }

    Success(c, gin.H{"message": "删除成功"})
}
```

#### 步骤6: 注册路由

```go
// internal/api/router/router.go
func SetupRouter(r *gin.Engine) {
    db := database.DB

    // 初始化
    articleDAO := impl.NewArticleDAO(db)
    articleService := service.NewArticleService(articleDAO)
    articleHandler := handler.NewArticleHandler(articleService)

    // 路由组（需要认证）
    api := r.Group("/api/v1")
    authenticated := api.Group("")
    authenticated.Use(middleware.JWTAuth("secret"))
    {
        articles := authenticated.Group("/articles")
        {
            articles.POST("", articleHandler.CreateArticle)
            articles.GET("/:id", articleHandler.GetArticle)
            articles.GET("", articleHandler.ListArticles)
            articles.PUT("/:id", articleHandler.UpdateArticle)
            articles.DELETE("/:id", articleHandler.DeleteArticle)
        }
    }
}
```

---

## 7. 总结

### 7.1 关键要点

1. **分层架构**: Handler → Service → DAO → Database
2. **依赖注入**: 通过构造函数注入，使用接口定义
3. **错误处理**: 错误作为返回值，统一错误响应格式
4. **Context传递**: 用于超时控制和值传递
5. **中间件**: 用于横切关注点（认证、日志、CORS等）

### 7.2 与C++的对比

| 概念    | C++         | Go        |
| ----- | ----------- | --------- |
| Web框架 | Crow/Drogon | Gin       |
| ORM   | ODB/SOCI    | GORM      |
| 错误处理  | 异常          | 返回值       |
| 并发    | std::thread | goroutine |
| 内存管理  | 手动          | GC        |

### 7.3 学习资源

- **Go官方文档**: https://go.dev/doc/
- **Gin文档**: https://gin-gonic.com/docs/
- **GORM文档**: https://gorm.io/docs/

---

**文档版本**: v1.0  
**最后更新**: 2025-01-XX  
**适用对象**: C++程序员转Go Web开发
