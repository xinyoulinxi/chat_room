
var chatRoom = ""
var username = ""
var password = ""
var userid = ""
initRoomList()

function getInput(title,showCancelButton, callback) {
    Swal.fire({
        title: title,
        input: 'text',
        inputAttributes: {
            autocapitalize: 'off'
        },
        showCancelButton: true,
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        showCancelButton: showCancelButton,
        showLoaderOnConfirm: true,
        preConfirm: (login) => {
            // 可以在这里处理输入的数据，如发送到服务器
            if (login == null || login == "") {
                showToast("用户名不能为空，请重新输入");
                return false;
            }
            return login; // 或者返回一个 promise
        },
        allowOutsideClick: () => false
    }).then((result) => {
        if (result.isConfirmed) {
            callback(result.value)
        }
    });
}

function showToast(text) {
    Toastify({
        text: text,
        duration: 2000, // Toast 持续显示的时间（毫秒）
        close: false, // 是否显示关闭按钮
        gravity: "top", // Toast 出现的位置，可以是 "top" 或 "bottom"
        position: 'center', // Toast 水平方向的位置，可以是 "left", "center", 或 "right"
        backgroundColor: "linear-gradient(to right, #00b09b, #96c93d)", // 背景色
        className: "chat-toast", // 自定义类名，用于添加特定的样式
        onClick: function () { } // 点击 Toast 时执行的函数
    }).showToast();
}
function initRoomList() {
    // 访问服务器，通过room_list接口获取房间列表
    fetch('/room_list')
        .then(response => response.json())
        .then(data => {
            // 更新房间列表
            var roomList = document.getElementById('roomList');
            roomList.innerHTML = '';
            // Create a select element
            var select = document.createElement('select');
            select.classList.add('select-text');
            if(data == null || data.length == 0) {
                return
            }
            chatRoom = data[0]
            select.onchange = function () {
                chatRoom = this.value;
            };
            data.forEach(room => {
                // Create an option for each room
                var option = document.createElement('option');
                option.classList.add('select-text');
                option.value = room;
                option.textContent = room;
                // Add the option to the select element
                select.appendChild(option);
            });
            select.selectedIndex = 0
            // Add the select element to the room list
            roomList.appendChild(select);
        })
        .catch(error => console.error('Error:', error));
}

function doRegister() {
    // 访问服务器，通过register接口注册新用户
    fetch('/register_user', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            username: username,
            password: password
        })
    })
        .then(response => response.json())
        .then(data => {
            console.log(data)
            if (data.errorCode !== 0) {
                showToast(data.message);
            } else {
                userid = data.message
                // 注册成功，跳转到聊天室页面
                goToChatRoom(userid,username,chatRoom);
            }
        })
        .catch(error => console.error('Error:', error));
}

function register() {
    username = document.getElementById('usernameInput').value;
    password = document.getElementById('passwordInput').value;
    if (username === '' || password === '') {
        showToast("用户名和密码不能为空，请重新输入");
        return;
    }
    doRegister(username,password)
}
function enterChatRoom(username, chatRoom) {

}
function login() {
    var username = document.getElementById('usernameInput').value;
    var password = document.getElementById('passwordInput').value;
    console.log("username:"+username+" password:"+password)
    if (username === '' || password === '') {
        showToast("用户名和密码不能为空，请重新输入");
        return;
    }
    if(chatRoom == null || chatRoom === ""){
        chatRoom = "default"
    }
    // 访问服务器，通过login接口验证用户名和密码
    fetch('/login_user', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            username: username,
            password: password
        })
    })
        .then(response => response.json())
        .then(data => {
            console.log(data)
            if (data.errorCode !== 0) {
                showToast(data.message);
            } else {
                userid = data.message
                // 登录成功，跳转到聊天室页面
                goToChatRoom(userid,username,chatRoom);
            }
        })
        .catch(error => console.error('Error:', error));
}

function goToChatRoom(userid,username,chatRoom){
    console.log("userid:"+userid+" chatRoom:"+chatRoom+" username:"+username)
    var url = '/chat_room?userid=' + encodeURIComponent(userid)+"&chatroom="+encodeURIComponent(chatRoom)+"&username="+encodeURIComponent(username);
    window.location.href = url;
}
