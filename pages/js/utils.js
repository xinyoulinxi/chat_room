function getCookie(name) {
    var cookieValue = "";
    if (document.cookie && document.cookie !== '') {
        var cookies = document.cookie.split(';');
        for (var i = 0; i < cookies.length; i++) {
            var cookie = cookies[i].trim();
            if (cookie.substring(0, name.length + 1) === (name + '=')) {
                cookieValue = decodeURIComponent(cookie.substring(name.length + 1));
                break;
            }
        }
    }
    return cookieValue;
}

function addCookie(name, value) {
    var cookieString = name + '=' + encodeURIComponent(value);
    document.cookie = cookieString;
}

function showToast(text, duration = 2000) {
    Toastify({
        text: text,
        duration: duration, // Toast 持续显示的时间（毫秒）
        close: false, // 是否显示关闭按钮
        gravity: "top", // Toast 出现的位置，可以是 "top" 或 "bottom"
        position: 'center', // Toast 水平方向的位置，可以是 "left", "center", 或 "right"
        backgroundColor: "linear-gradient(to right, #00b09b, #96c93d)", // 背景色
        className: "chat-toast", // 自定义类名，用于添加特定的样式
        onClick: function () {
        } // 点击 Toast 时执行的函数
    }).showToast();
}

function removeCookie(name) {
    document.cookie = name + '=; expires=Thu, 01 Jan 1970 00:00:01 GMT;';
}

function sleep(time) {
    return new Promise((resolve) => setTimeout(resolve, time));
}


function goToChatRoom(userid, username, chatRoom) {
    // 设置userId 和 username到cookie中
    addCookie("userId", userid)
    addCookie("userName", username)
    localStorage.setItem("chatRoom", chatRoom);
    window.location.href = '/chat_room'
}
