# API 请求 Header 格式说明

## 1. 登录接口（无需认证）

### 请求 Header
```
Content-Type: application/json
```

### 示例：登录请求
```http
POST /api/v1/auth/login HTTP/1.1
Host: localhost:8080
Content-Type: application/json

{
  "userName": "superAdmin",
  "password": "superAdmin"
}
```

### 响应示例
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "userName": "superAdmin",
      "role": "administrator",
      "isLocked": 0,
      "createTime": "2026-02-09T11:00:00Z",
      "updateTime": "2026-02-09T11:00:00Z"
    }
  }
}
```

---

## 2. 需要认证的接口

### 请求 Header（必需）
```
Content-Type: application/json
Authorization: Bearer <JWT_TOKEN>
```

**重要说明：**
- `Authorization` header 格式必须是：`Bearer <token>`（注意 Bearer 和 token 之间有一个空格）
- token 是登录接口返回的 JWT token
- token 有效期为配置文件中设置的 `jwt_expire_hours`（默认24小时）

### 示例：获取用户列表
```http
GET /api/v1/users HTTP/1.1
Host: localhost:8080
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VyX25hbWUiOiJzdXBlckFkbWluIiwicm9sZSI6ImFkbWluaXN0cmF0b3IiLCJleHAiOjE3MDc0ODQ4MDB9.xxxxx
```

### 示例：创建用户
```http
POST /api/v1/users HTTP/1.1
Host: localhost:8080
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "userName": "testuser",
  "password": "password123",
  "role": "user"
}
```

---

## 3. 使用不同客户端的示例

### JavaScript (Fetch API)
```javascript
// 登录
const loginResponse = await fetch('http://localhost:8080/api/v1/auth/login', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    userName: 'superAdmin',
    password: 'superAdmin'
  })
});

const loginData = await loginResponse.json();
const token = loginData.data.token;

// 使用token访问需要认证的接口
const usersResponse = await fetch('http://localhost:8080/api/v1/users', {
  method: 'GET',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  }
});
```

### JavaScript (Axios)
```javascript
import axios from 'axios';

// 创建axios实例
const api = axios.create({
  baseURL: 'http://localhost:8080/api/v1',
  headers: {
    'Content-Type': 'application/json'
  }
});

// 登录
const loginResponse = await api.post('/auth/login', {
  userName: 'superAdmin',
  password: 'superAdmin'
});

const token = loginResponse.data.data.token;

// 设置默认Authorization header
api.defaults.headers.common['Authorization'] = `Bearer ${token}`;

// 后续请求会自动携带token
const usersResponse = await api.get('/users');
```

### cURL
```bash
# 登录
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "userName": "superAdmin",
    "password": "superAdmin"
  }'

# 使用token访问需要认证的接口
curl -X GET http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Python (requests)
```python
import requests

base_url = "http://localhost:8080/api/v1"

# 登录
login_response = requests.post(
    f"{base_url}/auth/login",
    json={
        "userName": "superAdmin",
        "password": "superAdmin"
    },
    headers={"Content-Type": "application/json"}
)

token = login_response.json()["data"]["token"]

# 使用token访问需要认证的接口
headers = {
    "Content-Type": "application/json",
    "Authorization": f"Bearer {token}"
}

users_response = requests.get(
    f"{base_url}/users",
    headers=headers
)
```

### Postman
1. **登录请求：**
   - Method: `POST`
   - URL: `http://localhost:8080/api/v1/auth/login`
   - Headers:
     - `Content-Type: application/json`
   - Body (raw JSON):
     ```json
     {
       "userName": "superAdmin",
       "password": "superAdmin"
     }
     ```

2. **需要认证的请求：**
   - Method: `GET` / `POST` / `PUT` / `DELETE`
   - URL: `http://localhost:8080/api/v1/...`
   - Headers:
     - `Content-Type: application/json`
     - `Authorization: Bearer <从登录响应中复制的token>`

---

## 4. 错误响应

### 缺少 Authorization Header
```json
{
  "code": 401,
  "message": "未授权，请先登录",
  "data": null
}
```

### Token 格式错误
```json
{
  "code": 401,
  "message": "无效的授权头格式",
  "data": null
}
```

### Token 无效或过期
```json
{
  "code": 401,
  "message": "无效的token或token已过期",
  "data": null
}
```

---

## 5. 注意事项

1. **Token 格式**：Authorization header 必须是 `Bearer <token>` 格式，Bearer 和 token 之间必须有空格
2. **Token 有效期**：默认24小时，可在配置文件中修改 `auth.jwt_expire_hours`
3. **Content-Type**：所有包含请求体的接口都需要设置 `Content-Type: application/json`
4. **CORS**：服务端已配置CORS，支持跨域请求
5. **退出登录**：退出接口需要认证，但只是简单返回成功，客户端需要自行删除本地存储的token

---

## 6. 需要认证的接口列表

以下接口都需要在 Header 中携带 `Authorization: Bearer <token>`：

- `POST /api/v1/auth/logout` - 退出登录
- `POST /api/v1/users` - 创建用户
- `GET /api/v1/users` - 获取用户列表
- `GET /api/v1/users/{id}` - 获取用户详情
- `PUT /api/v1/users/{id}` - 更新用户
- `DELETE /api/v1/users/{id}` - 删除用户
- 所有 `/api/v1/pcd-files/*` 接口
- 所有 `/api/v1/semantic-maps/*` 接口
- 所有 `/api/v1/tasks/*` 接口
- 所有 `/api/v1/devices/*` 接口
- 所有 `/api/v1/operations/*` 接口
