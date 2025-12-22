📝 踩坑日记 (Troubleshooting)
1. 模板递归加载失效问题
现象描述：使用 r.LoadHTMLGlob("templates/**/*") 时，Gin 无法识别并响应 templates 根目录下的文件（如 index.html），控制台报 undefined 错误。

核心原因：Go 语言内置的 filepath.Glob 在处理 **/* 模式时存在局限性。它会强制要求匹配至少一级子目录，从而忽略了位于搜索根目录下的孤立文件。

解决方案：通过 filepath.Walk 自定义递归解析函数。
