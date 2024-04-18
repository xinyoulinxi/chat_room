// Establish a WebSocket connection
var params = new URLSearchParams(window.location.search);
var userId = params.get("userid")
var userName = params.get("username")
var chatRoom = params.get("chatroom")
var socket = null;
var messageDisplay = document.getElementById('messageDisplay'); // 全局缓存
var messageInput = document.getElementById('messageInput'); // 全局缓存

init()

function init() {
    messageInput.addEventListener('keydown', function (event) {
        if (!event.isComposing && event.key === 'Enter') {
            event.preventDefault(); // Prevent form submission
            sendMessage();
        }
    });
    // Focus on the message input field
    messageInput.focus();
    connectToChatRoom();
}

function getHistoryMessages() {
    fetch('/history_messages', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            userId: userId,
            chatRoom: chatRoom
        })
    })
        .then(response => response.json())
        .then(data => {
            if (data.errorCode !== 0) {
                showToast(data.message);
            } else {
                console.log("get history message size", data.data.length)
                data.data.forEach(message => {
                    handleMessage(message)
                })
                messageDisplay.scrollTop = messageDisplay.scrollHeight;
                // Initialize the WebSocket connection
                initSocket();
            }
        })
        .catch(error => console.error('Error:', error));
}

function getInput(title, showCancelButton, callback) {
    Swal.fire({
        title: title,
        input: 'text',
        inputAttributes: {
            autocapitalize: 'off'
        },
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        showCancelButton: showCancelButton,
        showLoaderOnConfirm: true,
        preConfirm: (login) => {
            // 可以在这里处理输入的数据，如发送到服务器
            if (login == null || login === "") {
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

function getRoomList() {
    fetch('/room_list?id=' + userId)
        .then(response => response.json())
        .then(data => {
            if (data.errorCode !== 0) {
                showToast(data.message);
            } else {
                handleRoomList(data.data)
            }
        });
}

function isChatRoomExist(roomName, roomList) {
    for (var i = 0; i < roomList.length; i++) {
        if (roomList[i] === roomName)
            return true;
    }
    return false;
}

function handleRoomList(chatRoomList) {
    if (chatRoomList == null || chatRoomList.length === 0) {
        showToast("当前没有聊天室，请先创建一个")
        return;
    }
    var roomList = document.getElementById('roomList');
    roomList.innerHTML = ""; // Clear the room list
    var select = roomList.querySelector('select') || document.createElement('select');
    select.innerHTML = ""; // 仅清空select内部，避免重复创建select元素
    select.classList.add('select-text-selected');
    select.onchange = function () {
        chatRoom = this.value;
        connectToChatRoom();
    };
    if (chatRoom == null || chatRoom === "" || isChatRoomExist(chatRoom, chatRoomList) === false) {
        showToast("当前聊天室不存在，已自动切换到第一个聊天室")
        chatRoom = chatRoomList[0];
        connectToChatRoom();
    }

    // 将select设置为当前chatroom的值
    for (var i = 0; i < chatRoomList.length; i++) {
        var room = chatRoomList[i];
        // Create an option for each room
        var option = document.createElement('option');
        option.classList.add('select-text');
        option.value = room;
        option.textContent = room;
        // Add the option to the select element
        select.appendChild(option);
        if (chatRoom === room) {
            option.classList.add("select-text-selected")
        }
    }
    select.selectedIndex = chatRoomList.indexOf(chatRoom);
    // Add the select element to the room list
    roomList.appendChild(select);
}

function closeSocket() {
    if (socket != null) {
        socket.close();
    }
}

function displayFileMessage(message) {
    const fileElement = document.createElement('a');
    const downLoadElement = document.createElement('a');
    // 让文件在新窗口打开
    fileElement.target = "_blank";
    // 点击downLoadElement下载文件
    downLoadElement.download = message.content;
    downLoadElement.href = message.file;
    downLoadElement.textContent = "[下载]";
    downLoadElement.classList.add('message-download');
    fileElement.href = message.file; // Set src to the image field of the message
    fileElement.textContent = message.content;
    fileElement.classList.add('message-bubble');
    fileElement.classList.add('message-bubble-text');
    if (message.userName === userName || message.userId === userId) {
        downLoadElement.classList.add('message-download-me');
        fileElement.classList.add('my-message');
    } else {
        downLoadElement.classList.add('message-download-other');
        const userName = document.createElement('div');
        userName.className = 'user-name';
        userName.textContent = message.userName;
        messageDisplay.appendChild(userName);
        fileElement.classList.add('other-message');
    }
    messageDisplay.appendChild(fileElement);
    messageDisplay.appendChild(downLoadElement);
}

function handleMessage(message) {
    insertSendTime(message)
    console.log("handleMessage", message)
    if (message.type === "image") {
        displayImageMessage(message);
    } else if (message.type === "file") {
        displayFileMessage(message);
    } else if (message.type === "text" || message.type === "") {
        displayNormalMessage(message);
    }else if(message.type === "userCount"){
        console.log("userCount",message.data)
        document.getElementById("userCount").textContent = message.data.UserCount
    }else if(message.type === "roomList"){
        console.log("roomList",message.data)
        handleRoomList(message.chatRoomList)
    }
    if (message.userId === userId || message.userName === userName) {
        scrollToBottom(0)
    } else {
        scrollToBottom()
    }
}

function scrollToBottom(offset = 500) {
    if (offset <= 0) {
        messageDisplay.scrollTop = messageDisplay.scrollHeight;
    } else if (messageDisplay.scrollHeight - messageDisplay.scrollTop - messageDisplay.clientHeight <= offset) {// 检查滚动条是否距离底部一个指定的偏移量
        // 如果滚动条距离底部小于等于指定偏移量，则将滚动条滚动到底部
        messageDisplay.scrollTop = messageDisplay.scrollHeight;
    }
}

function getOkTimeText(currentDate, time) {
    if (currentDate.getDate() === new Date(time).getDate()) {
        return time.substring(11, 16);
    }
    // 昨天显示昨天加时间
    if (currentDate.getDate() === new Date().getDate() - 1) {
        return "昨天 " + time.substring(11, 16);
    }
    return time.substring(0, 16);
}

var lastTime = ""; // "2024-11-12 00:00:00"
// 判断当前的time是否需要显示时间，如果需要则返回time，否则返回""，并更新lastTime
function getTimeInterval(time) {
    var currentDate = new Date();
    if (lastTime === "" || lastTime.search("年") !== -1) {
        lastTime = time;
        return getOkTimeText(currentDate, time)
    }
    var lastDate = new Date(lastTime);
    var interval = currentDate - lastDate;
    if (interval > 1000 * 60 * 5) {
        lastTime = time;
        // 如果是今天的消息，不显示年月日，只显示时分
        return getOkTimeText(currentDate, time);
    }

    return "";
}


function insertSendTime(message) {
    const time = getTimeInterval(message.sendTime);
    if (time === "") {
        return;
    }
    const timeElement = document.createElement('div');
    timeElement.classList.add("message-time-text")
    timeElement.textContent = time;
    timeElement.classList.add('send-time');
    messageDisplay.appendChild(timeElement);

}

function initSocket() {
    closeSocket();
    lastTime = ""
    socket = new WebSocket('ws://' + window.location.host + '/ws?id=' + userId + '&chatroom=' + chatRoom);
    // Event listener for receiving messages from the server
    socket.onmessage = function (event) {
        var messages = JSON.parse(event.data); // Parse the JSON data into an array
        messages.forEach(function (message) { // Iterate over each message in the array
            // console.log(message)
            handleMessage(message);
        });
        // Scroll to the bottom of the message display
        // messageDisplay.scrollTop = messageDisplay.scrollHeight;
    };
    socket.onopen = function (event) {
        console.log("socket open", event)
        document.title = `Chat Room - ${chatRoom}`;
    };
    // 自动重连
    socket.onclose = function (event) {
        if (socket.readyState === WebSocket.CLOSED) {
            console.log("socket close")
            setTimeout(function () {
                showToast("连接断开，尝试重新连接...", {duration: 3000});
                connectToChatRoom();
            }, 5000);
        }
    };
}

function createRoom() {
    getInput("输入需要创建的聊天室名称", true, function (roomName) {
        if (roomName == null || roomName === "") {
            return;
        }
        fetch('/create_room?id=' + userId + '&roomName=' + roomName)
            .then(response => response.json())
            .then(data => {
                if (data.errorCode === 0) {
                    showToast("聊天室创建成功");
                    chatRoom = roomName;
                    connectToChatRoom();
                } else {
                    showToast(data.message);
                }
            });
    });

}

function goLoginPage() {
    window.location.href = "/login";
}

function connectToChatRoom() {
    // Check if the userId or room number is empty
    if (userId == null || userId === "" || chatRoom == null || chatRoom === "") {
        showToast('请登录并选择聊天室后再进入');
        return;
    }
    // clear the message display
    messageDisplay.innerHTML = "";
    // Get the chat history
    getHistoryMessages();
    // Get the room list
    getRoomList();
}

// Function to send a message
function sendMessage() {
    const message = messageInput.value;
    if (message == null || message === "") {
        showToast('请输入消息内容');
        return;
    }
    // Create a message object with username and content
    const messageObj = {
        userName: userName,
        type: "text",
        userId: userId,
        roomName: chatRoom,
        content: message,
        image: null
    };
    // Send the message to the server
    socket.send(JSON.stringify(messageObj));

    // Clear the input field
    messageInput.value = '';
}

function uploadFile(fileName, data, callback) {
    fetch('/upload_file', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            userId: userId,
            fileName: fileName,
            data: data
        })
    })
        .then(response => response.json())
        .then(data => {
            if (data.errorCode !== 0) {
                showToast(data.message);
            } else {
                console.log("upload file success")
                callback(data.data.filePath, data.data.fileType)
            }
        })
        .catch(error => console.error('Error:', error));
}

function displayImageMessage(message) {
    const imageElement = document.createElement('img');
    imageElement.src = message.image; // Set src to the image field of the message
    imageElement.dataset.retry = 0
    imageElement.loading = 'lazy'; // 设置图片懒加载
    imageElement.classList.add('message-bubble');
    imageElement.classList.add('message-bubble-img');
    imageElement.onerror = function (e) {
        let src = e.target.src
        let retry = parseInt(e.target.dataset.retry) || 0
        if (retry < 3) {
            e.target.dataset.retry = retry + 1
            console.log("图片加载失败 重试", retry, src)
            setTimeout(function () {
                e.target.src = src
            }, 3000 * (retry + 1))
        }
    }
    imageElement.onload = function (e) {
        scrollToBottom(500)
    }
    if (message.userId === userId || message.userName === userName) {
        imageElement.classList.add('my-message');
    } else {
        const userName = document.createElement('div');
        userName.className = 'user-name';
        userName.textContent = message.userName;
        messageDisplay.appendChild(userName);
        imageElement.classList.add('other-message');
    }
    messageDisplay.appendChild(imageElement);
    $(document).ready(function () {
        // Your code here
        imageElement.addEventListener('click', function () {
            $.fancybox.open({
                src: this.src,
                type: 'image'
            });
        });
    });
}

// Function to display a message
function displayNormalMessage(message) {
    const messageElement = document.createElement('div');
    // messageElement.textContent = message.content; // Add username before message
    const messageText = document.createElement('pre');
    messageText.textContent = message.content
    messageElement.appendChild(messageText)

    messageElement.classList.add('message-bubble');
    if (message.userId === userId || message.userName === userName) {
        messageElement.classList.add('my-message');
        messageElement.align = "right"
    } else {
        const userName = document.createElement('div');
        userName.className = 'user-name';
        userName.textContent = message.userName;
        userName.align = "left"
        messageDisplay.appendChild(userName);
        messageElement.classList.add('other-message');
    }
    messageDisplay.appendChild(messageElement);
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

function sendFile() {
    const fileInput = document.getElementById('image-select');
    // 确保有文件被选中
    if (fileInput.files.length > 0) {
        const file = fileInput.files[0]; // 获取第一个文件
        const fileName = file.name; // 获取文件名
        console.log(fileName, file.size)
        if (file.size > 1024 * 1024 * 20) {
            showToast("文件大小不能超过20MB");
            return;
        }
        const reader = new FileReader();
        reader.onload = function (e) {
            // 限制文件大小
            if (file.size > 1024 * 1024 * 20) {
                showToast("文件大小不能超过20MB");
                return;
            }
            uploadFile(fileName, e.target.result, function (filePath, fileType) {
                // 创建一个包含用户名、内容、图片和文件名的消息对象
                console.log("upload file", filePath, fileType)
                const messageObj = {
                    type: fileType,
                    userId: userId,
                    content: fileName,
                    userName: userName,
                    roomName: chatRoom,
                };
                if (fileType === 'image')
                    messageObj.image = filePath;
                else
                    messageObj.file = filePath;
                socket.send(JSON.stringify(messageObj));
                showToast("文件发送成功");
                fileInput.value = ''; // 重置文件输入以便下次使用
            })
        };
        reader.readAsDataURL(file);
    } else {
        showToast("请选择一个文件");
    }
}

