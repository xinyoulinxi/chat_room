
var chatRoom = ""
var username = ""
var password = ""
init()

function init() {
    initRoomList()
}
function getInput(title, showCancelButton, callback) {
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

function initRoomList() {
    // 访问服务器，通过room_list接口获取房间列表
    fetch('/room_list')
        .then(response => response.json())
        .then(data => {
            console.log(data)
            var roomListData = data.data
            // 更新房间列表
            var roomList = document.getElementById('roomList');
            roomList.innerHTML = '';
            // Create a select element
            var select = document.createElement('select');
            select.classList.add('select-text');
            if (roomListData == null || roomListData.length === 0) {
                return
            }
            chatRoom = localStorage.getItem("chatRoom");
            // 自动登录
            // if(getCookie("userId") !== "" && getCookie("userName") !== ""){
            //     var userid = getCookie("userId")
            //     var username = getCookie("userName")
            //     goToChatRoom(userid, username, chatRoom)
            //     return
            // }
            select.onchange = function () {
                chatRoom = this.value;
            };
            var index = 0
            roomListData.forEach(room => {
                // Create an option for each room
                var option = document.createElement('option');
                option.classList.add('select-text');
                option.value = room;
                option.textContent = room;
                // Add the option to the select element
                select.appendChild(option);
                if (chatRoom === room) {
                    select.selectedIndex = index
                }
                index++
            });
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
                var userid = data.message
                // 注册成功，跳转到聊天室页面
                goToChatRoom(userid, username, chatRoom);
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
    doRegister(username, password)
}

function login() {
    var username = document.getElementById('usernameInput').value;
    var password = document.getElementById('passwordInput').value;
    console.log("username:" + username + " password:" + password)
    if (username === '' || password === '') {
        showToast("用户名和密码不能为空，请重新输入");
        return;
    }
    if (chatRoom == null || chatRoom === "") {
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
                localStorage.setItem("username", username);
                localStorage.setItem("password", password);
                var userid = data.message
                // 登录成功，跳转到聊天室页面
                goToChatRoom(userid, username, chatRoom);
            }
        })
        .catch(error => console.error('Error:', error));
}
