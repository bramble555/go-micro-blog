/**
 * Article management functions
 */

async function deleteArticle(articleId) {
    const token = localStorage.getItem('token');
    const csrfTokenTag = document.querySelector('meta[name="csrf-token"]');
    const csrfToken = csrfTokenTag ? csrfTokenTag.content : '';

    try {
        const response = await fetch(`/api/admin/articles/${articleId}/delete`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': token ? 'Bearer ' + token : '',
                'X-CSRF-Token': csrfToken
            }
        });

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            throw new Error(errorData.msg || '删除失败');
        }
        
        return await response.json();
    } catch (error) {
        console.error('删除文章错误:', error);
        throw error;
    }
}
