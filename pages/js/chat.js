
// Establish a WebSocket connection
var params = new URLSearchParams(window.location.search);
var username = params.get("username")
var chatroom = params.get("chatroom")
var socket = null;

console.log("username: " + username, "chatroom: " + chatroom);
if (username == null || username == "") {
    username = getInput("请输入用户名: ",false, function (name) {
        username = name
        init()
    });
}

init()
function init() {
    document.getElementById('messageInput').addEventListener('keydown', function (event) {
        if (event.key === 'Enter') {
            event.preventDefault(); // Prevent form submission
            sendMessage();
        }
    });
    connectToChatRoom(username, chatroom);
}

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
            console.log("User's name: ", login);
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

function handleRoomList(message) {
    var roomList = document.getElementById('roomList');
    roomList.innerHTML = ""; // Clear the room list
    // clear the message display
    const messageDisplay = document.getElementById('messageDisplay');
    messageDisplay.innerHTML = "";
    // Create a select element
    var select = document.createElement('select');

    select.classList.add('select-text-selected');
    select.onchange = function () {
        chatroom = this.value;
        initSocket();
    };
    if (chatroom == null || chatroom == "") {
        chatroom = message.chatRoomList[0];
    }

    // 将select设置为当前chatroom的值
    for (var i = 0; i < message.chatRoomList.length; i++) {
        var room = message.chatRoomList[i];
        // Create an option for each room
        var option = document.createElement('option');
        option.classList.add('select-text');
        option.value = room;
        option.textContent = room;
        // Add the option to the select element
        select.appendChild(option);
        if (chatroom == room) {
            option.classList.add("select-text-selected")
        }
    }
    select.selectedIndex = message.chatRoomList.indexOf(chatroom);
    // Add the select element to the room list
    roomList.appendChild(select);
}

function closeSocket() {
    if (socket != null) {
        socket.close();
    }
}

function handleMessage(message) {
    insertSendTime(message)
    if (message.type == "image") {
        displayImageMessage(message);
    } else if (message.type == "text" || message.type == "") {
        displayNormalMessage(message);
    } else {
        handleRoomList(message)
    }
}

function getOkTimeText(currentDate, time) {
    if (currentDate.getDate() == new Date(time).getDate()) {
        console.log("today")
        return time.substring(11, 16);
    }
    // 昨天显示昨天加时间
    if (currentDate.getDate() == new Date().getDate() - 1) {
        return "昨天 " + time.substring(11, 16);
    }
    return time.substring(0, 16);
}
var lastTime = ""; // "2024-11-12 00:00:00"
// 判断当前的time是否需要显示时间，如果需要则返回time，否则返回""，并更新lastTime
function getTimeInterval(time) {
    var currentDate = new Date();
    if (lastTime == "" || lastTime.search("年") != -1) {
        lastTime = time;
        return getOkTimeText(currentDate, time)
    }
    var lastDate = new Date(lastTime);
    var interval = currentDate - lastDate;
    console.log("interval: " + interval + " currentDate: " + currentDate + " lastDate: " + lastDate + " lastTime：" + lastTime + " time: " + time)
    if (interval > 1000 * 60 * 5) {
        console.log(" in interval: " + interval + " currentDate: " + currentDate + " lastDate: " + lastDate + " lastTime：" + lastTime + " time: " + time)
        lastTime = time;
        // 如果是今天的消息，不显示年月日，只显示时分
        console.log("currentDate.getDate(): " + currentDate.getDate() + " new Date(time).getDate(): " + new Date(time).getDate())
        return getOkTimeText(currentDate, time);
    }

    return "";
}



function insertSendTime(message) {
    const time = getTimeInterval(message.sendTime);
    if (time == "") {
        return;
    }
    const messageDisplay = document.getElementById('messageDisplay');
    const timeElement = document.createElement('div');
    timeElement.classList.add("message-time-text")
    timeElement.textContent = time;
    timeElement.classList.add('send-time');
    messageDisplay.appendChild(timeElement);

}

function initSocket() {
    closeSocket();
    lastTime = ""
    socket = new WebSocket('ws://' + window.location.host + '/ws?id=' + username + '&chatroom=' + chatroom);
    // Event listener for receiving messages from the server
    socket.onmessage = function (event) {
        console.log(event.data);
        var messages = JSON.parse(event.data); // Parse the JSON data into an array
        messages.forEach(function (message) { // Iterate over each message in the array
            handleMessage(message);
        });
        // Scroll to the bottom of the message display
        messageDisplay.scrollTop = messageDisplay.scrollHeight;
    };
    socket.onopen = function (event) {
        console.log('WebSocket is connected.');
    };
}

function createRoom() {
    getInput("输入需要创建的聊天室名称",true, function (roomName) {
        if (roomName == null || roomName == "") {
            return;
        }
        fetch('/create_room?roomName=' + roomName)
            .then(response => response.json())
            .then(data => {
                console.log(data);
                if (data.errorCode == 0) {
                    showToast("聊天室创建成功");
                    chatroom = roomName;
                    initSocket();
                } else {
                    showToast(data.message);
                }
            });
    });

}

function connectToChatRoom() {
    // Check if the username or room number is empty
    if (username == null || username == "") {
        showToast('请输入用户名');
        return;
    }

    // Focus on the message input field
    document.getElementById('messageInput').focus();

    // Initialize the WebSocket connection
    initSocket();
}

// Function to send a message
function sendMessage() {
    const messageInput = document.getElementById('messageInput');
    const message = messageInput.value;
    if (message == null || message == "") {
        showToast('请输入消息内容');
        return;
    }
    // Create a message object with username and content
    const messageObj = {
        userId: username,
        content: message,
        image: null
    };

    console.log(messageObj);
    // Send the message to the server
    socket.send(JSON.stringify(messageObj));

    // Clear the input field
    messageInput.value = '';
}

function displayImageMessage(message) {
    const messageDisplay = document.getElementById('messageDisplay');
    const imageElement = document.createElement('img');
    imageElement.src = message.image; // Set src to the image field of the message
    imageElement.classList.add('message-bubble');
    imageElement.classList.add('message-bubble-img');
    if (message.userId === username) {
        imageElement.classList.add('my-message');
    } else {
        const userName = document.createElement('div');
        userName.className = 'user-name';
        userName.textContent = message.userId;
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
    const messageDisplay = document.getElementById('messageDisplay');
    const messageElement = document.createElement('div');
    messageElement.textContent = message.content; // Add username before message
    messageElement.classList.add('message-bubble');
    if (message.userId === username) {
        messageElement.classList.add('my-message');
        messageElement.align = "right"
    } else {
        const userName = document.createElement('div');
        userName.className = 'user-name';
        userName.textContent = message.userId;
        userName.align = "left"
        messageDisplay.appendChild(userName);
        messageElement.classList.add('other-message');
    }
    messageDisplay.appendChild(messageElement);
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

function sendImage() {
    const imageInput = document.getElementById('image-select');
    const reader = new FileReader();
    reader.onload = function (e) {
        // Create a message object with username, content, and image
        const messageObj = {
            userId: username,
            room: chatroom,
            content: "",
            image: e.target.result // Base64-encoded image data
        };
        socket.send(JSON.stringify(messageObj));
        imageInput.value = '';
    };
    reader.readAsDataURL(imageInput.files[0]);
    showToast("图片发送成功");
}
