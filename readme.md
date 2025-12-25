## 🛠️ 踩坑日记与技术挑战 (Troubleshooting & Technical Challenges)

### 1. 模板递归解析失效 | Template Recursive Loading

---

* **现象描述 (The Symptom)**
在使用 `r.LoadHTMLGlob("templates/**/*")` 时，Gin 仅能加载子目录下的模板，无法识别并响应 `templates` 根目录下的文件（如 `index.html`），渲染时触发 `html/template: undefined` 异常。
* **核心原因 (The Root Cause)**
Go 内置的 `filepath.Glob` 对通配符 `**` 的递归支持存在局限性。它往往要求路径中**必须**包含至少一级子目录，从而忽略了直接位于根目录下的 HTML 文件。
* **解决方案 (The Solution)**
通过自定义 `filepath.Walk` 函数手动遍历目录，彻底解决递归解析不全的问题。

---

### 2. 雪花算法 ID 精度丢失 | Snowflake ID Precision Loss

---

* **问题描述 (The Symptom)**
后端返回的 `int64` ID（如 `...9100`），在前端 JS 控制台打印却变成了 `...9104`。这种“自动进位”导致前端拿着错误的 ID 请求 API，后端返回 **404 Not Found**。
* **原因分析 (The Root Cause)**
* **后端 (Go):** `int64` 范围达 （约 19 位数字）。
* **前端 (JavaScript):** 遵循 **IEEE 754** 标准，其“安全整数”上限仅为 （约 16 位数字）。
* **冲突:** 19 位的雪花 ID 超出了 JS 的处理极限，解析时会发生舍入误差。


* **解决方案 (The Solution)**
在 Go Model 定义中，利用结构体标签 `json:",string"` 强制将大整数序列化为**字符串**传输。

```go
type Comment struct {
    // 强制转为字符串，避免 JS 解析时精度丢失
    ID        int64     `json:"id,string"`         
    ArticleID int64     `json:"article_id,string"`  
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
}

```

