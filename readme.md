🛠️踩坑日记:
1. 
问题：使用 r.LoadHTMLGlob("templates/**/*") 无法加载 templates 根目录下的文件（如 index.html）。 
原因：Go 的 filepath.Glob 在使用 **/* 模式时，会强制匹配子目录，导致根目录文件被忽略。 
解决方法：在main.go 文件里面自己写一个函数递归遍历目的目录下的所有.html文件