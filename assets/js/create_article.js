const API_BASE_URL = "/api/admin/articles";

// 1. 发送文章数据到后端
async function createArticle(data) {
    const token = localStorage.getItem('token');
    const headers = { 'Content-Type': 'application/json' };
    if (token) {
        headers['Authorization'] = 'Bearer ' + token;
    }

    const response = await fetch(API_BASE_URL, {
        method: 'POST',
        headers: headers,
        body: JSON.stringify(data),
    });

    if (response.ok) {
        const result = await response.json();
        return result;  // 返回成功的响应
    } else {
        const error = await response.json();
        alert('Failed to create article: ' + error.message);
        throw new Error('Failed to create article');
    }
}

// 2. 处理表单提交
document.getElementById("createArticleForm").addEventListener("submit", async function (event) {
    event.preventDefault();  // 防止默认的表单提交行为

    const title = document.getElementById("title").value;
    const content = document.getElementById("content").value;

    const articleData = { title, content };

    try {
        const response = await createArticle(articleData);
        alert('发布成功！');
        // 修改这里：跳转到首页
        window.location.href = '/';
    } catch (error) {
        console.error("Error creating article:", error);
    }
});
