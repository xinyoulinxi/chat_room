var userId = getCookie("userId")
var userName = getCookie("userName")
var chatRoom = localStorage.getItem("chatRoom");
var socket = null;// WebSocket
var messageDisplay = document.getElementById('messageDisplay');// 消息显示区域
var messageInput = document.getElementById('messageInput');// 消息输入框
let pageVisibility = true// 页面是否可见
let unreadMessageCount = 0// 未读消息计数
let settingMenuShow = false// 设置菜单是否显示
let userListMenuShow = false// 用户列表菜单是否显示
var createChatRoomBtn = document.getElementById('create-room');// 创建聊天室按钮
var unLoginBtn = document.getElementById('un-login');// 退出登录按钮
var profileBtn = document.getElementById('profile-btn');// 个人中心按钮
var helpBtn = document.getElementById('help');// 帮助按钮
var messageIndex = 0;// 消息索引
var hasMoreMessage = true;// 是否还有更多消息
var hasLoadedHistory = false; // 是否已经加载过历史消息

// 初始化
init()

function init() {
    initPopupMenu()
    initUserListMenu()
    initListener()
    // Focus on the message input field
    messageInput.focus();
    connectToChatRoom(chatRoom);
}

function goProfile() {
    window.location.href = "/profile";
}

function initListener() {
    messageInput.addEventListener('paste', function (event) {
        // 检查粘贴的内容是否为图片
        if (event.clipboardData.items && event.clipboardData.items[0].type.startsWith('image/')) {
            const file = event.clipboardData.items[0].getAsFile();
            let fileName = Math.random().toString(36).substring(2, 10) + "_" + file.name
            const reader = new FileReader();
            // 读取文件内容
            reader.onload = function (e) {
                Swal.fire({
                    title: '是否发送图片?',
                    html: `<img src="${e.target.result}" style="width: 100%" />`,
                    type: 'question',
                    showCancelButton: true,
                    confirmButtonColor: '#95EC69',
                    // cancelButtonColor: '#d33',
                    confirmButtonText: '发送',
                    cancelButtonText: '取消',
                }).then(function (res) {
                    console.log(res)
                    if (res.isConfirmed) {
                        uploadFile(fileName, e.target.result, function (filePath, fileType) {
                            // 创建一个包含用户名、内容、图片和文件名的消息对象
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
                            messageInput.value = ''; // 重置文件输入以便下次使用
                        })
                    }
                })

            };
            reader.readAsDataURL(file);
        }
    })
    // 监听发送消息按钮
    messageInput.addEventListener('keydown', function (event) {
        if (!event.isComposing && event.key === 'Enter') {
            event.preventDefault(); // Prevent form submission
            sendMessage();
        }
    });
    // 监听页面显示状态，离开页面时增加消息记数且阻止消息滚动，重新进入页面后重置记数并自动滚动到底部最新消息
    document.addEventListener('visibilitychange', function () {
        switch (document.visibilityState) {
            case 'hidden':
                console.log("离开时间点：" + new Date());
                pageVisibility = false;
                break;
            case 'visible':
                console.log("重新进入时间点：" + new Date());
                pageVisibility = true;
                unreadMessageCount = 0;
                document.title = `Chat Room - ${chatRoom}`;
                scrollToBottom(0)
                break;
            default:
                break;
        }
    })
    messageDisplay.addEventListener('scroll', function () {
        if (messageDisplay.scrollTop === 0 && hasLoadedHistory) {
            getHistoryMessages(false)
        }
    });
    messageDisplay.addEventListener('touchmove', function (event) {
        // var currentTouchY = event.touches[0].clientY;
        if (messageDisplay.scrollTop === 0 && hasLoadedHistory) {
            getHistoryMessages(false)
        }
    });
}

function initUserListMenu() {
    userListMenu = document.getElementById("drop-user-list")
    userCountBtn = document.getElementById("userCount")
    userCountBtn.addEventListener('click', function () {
        console.log("click user count", userListMenuShow)
        if (userListMenuShow) {
            userListMenu.style.display = 'none';
            userListMenuShow = !userListMenuShow
            return
        }
        userListMenuShow = !userListMenuShow
        userListMenu.style.display = 'block';
        Popper.createPopper(userCountBtn, userListMenu, {
            placement: 'bottom-start',
        });
    })
    document.addEventListener('click', function (event) {
        if (event.target !== userCountBtn) {
            userListMenu.style.display = 'none';
        }
    });
}

function initPopupMenu() {
    createChatRoomBtn.addEventListener('click', createRoom);
    unLoginBtn.addEventListener('click', unLogin);
    profileBtn.addEventListener('click', goProfile);
    document.getElementById('update-avatar-btn').addEventListener('click', updateAvatar);
    var button = document.getElementById('menu-button');
    var menu = document.getElementById('drop-menu');
    button.addEventListener('click', function () {
        if (settingMenuShow) {
            menu.style.display = 'none';
            settingMenuShow = !settingMenuShow
            return
        }
        settingMenuShow = !settingMenuShow
        menu.style.display = 'block';
        Popper.createPopper(button, menu, {
            placement: 'bottom-end',
        });
    });

    document.addEventListener('click', function (event) {
        if (event.target !== button) {
            menu.style.display = 'none';
        }
    });
}

/**
 * 发送浏览器通知
 * @param title 标题
 * @param message 内容
 */
function browserNotify(title, message) {
    // if (window.Notification && Notification.permission !== "denied") {
    //     try {
    //         Notification.requestPermission(function () {
    //             new Notification(title, {body: message});
    //         }).catch(e=>console.log(e));
    //     }catch (e){
    //         console.log(e)
    //     }
    // }
}

function getHistoryMessages(isFirstFetch = true) {
    console.log("getHistoryMessages", messageIndex, isFirstFetch)
    fetch('/history_messages', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            userId: userId,
            chatRoom: chatRoom,
            index: messageIndex,
            count: 100
        })
    })
        .then(response => response.json())
        .then(data => {
            if (data.errorCode === 8) {
                hasMoreMessage = false
                showToast("暂无更多消息")
            } else if (data.errorCode !== 0) {
                showToast(data.message);
            } else {
                console.log("getHistoryMessages success size:", data.data.length)
                messageIndex++
                if (isFirstFetch) {
                    data.data.forEach(message => {
                        handleMessage(message)
                    })
                    messageDisplay.scrollTop = messageDisplay.scrollHeight;
                    hasLoadedHistory = true
                } else {
                    showToast("已加载更多消息")
                    var scrollTop = messageDisplay.scrollHeight;
                    for (var i = data.data.length - 1; i >= 0; i--) {
                        handleHistoryMessage(data.data[i])
                    }
                    messageDisplay.scrollTop = messageDisplay.scrollHeight - scrollTop
                }
            }
        })
        .catch(error => console.error('Error:', error));
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
        connectToChatRoom(this.value);
    };
    if (chatRoom == null || chatRoom === "" || isChatRoomExist(chatRoom, chatRoomList) === false) {
        showToast("当前聊天室不存在，已自动切换到第一个聊天室")
        connectToChatRoom(chatRoomList[0]);
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

var avatar_fetching = new Map();
var avatar_map = new Map()

function getAvatarUrl(userName, avatarElement) {
    if (avatar_map.has(userName)) {
        // console.log("use cache", userName, avatar_map.get(userName))
        avatarElement.src = avatar_map.get(userName)
        return
    }
    if (avatar_fetching.has(userName)) {
        avatar_fetching.get(userName).then(url => {
            // console.log("get avatar finish, use cache", userName, url)
            avatarElement.src = url;
        });
        return;
    }
    let fetchPromise = fetch('/get_avatar?userName=' + userName)
        .then(response => response.json())
        .then(data => {
            if (data.errorCode === 0) {
                // console.log("fetch avatar by net success", userName, data.data)
                avatar_map.set(userName, data.data)
                avatarElement.src = data.data
                return data.data;
            }
        });
    avatar_fetching.set(userName, fetchPromise);
    fetchPromise.finally(() => {
        avatar_fetching.delete(userName);
    });
}

function updateAvatar() {
    getInput("输入头像地址", true, function (avatarUrl) {
        if (avatarUrl == null || avatarUrl === "") {
            return;
        }
        fetch('/update_avatar?id=' + userId + '&url=' + avatarUrl)
            .then(response => response.json())
            .then(data => {
                if (data.errorCode === 0) {
                    showToast("头像设置成功");
                    document.getElementById("avatar").src = avatarUrl
                } else {
                    showToast(data.message);
                }
            });
    });
}

function displayNoticeElement(message, isAddToBottom) {
    const messageElement = document.createElement('div');
    messageElement.textContent = message.content
    messageElement.classList.add('message-notice-text')
    addMessageToDisplay(messageElement, isAddToBottom)
}

function displayProfileElement(message, element, isAddToBottom) {
    // 创建最外层容器
    const container = document.createElement('div');
    container.classList.add('message-container');

    // 创建并设置图片元素
    const avatar = document.createElement('img');
    avatar.src = "data/default_avatar.webp" // message.avatarUrl;
    avatar.alt = 'S';
    avatar.className = 'avatar';
    getAvatarUrl(message.userName, avatar)
    avatar.onclick = function () {
        showBigImage(this.src)
    }
    // 创建文本容器
    const textContainer = document.createElement('div');
    textContainer.className = 'text-container';

    // 创建并设置名字元素
    const nameDiv = document.createElement('div');
    nameDiv.className = 'user-name';
    nameDiv.textContent = message.userName;

    // 创建并设置内容元素
    const contentDiv = document.createElement('div');
    element.classList.add('message-bubble')
    var isSelf = message.userId === userId || message.userName === userName
    if (isSelf) {
        container.classList.add("my-message-container")
        container.align = "right"
        element.classList.add('my-message')
    } else {
        // container.align = "left"
        container.classList.add("other-message-container")
        element.classList.add('other-message')
    }
    contentDiv.textContent = message.content;

    if (!isSelf) {
        container.appendChild(avatar);
    }
    textContainer.appendChild(nameDiv);
    textContainer.appendChild(element);
    // 将文本容器加入到最外层容器中
    container.appendChild(textContainer);
    if (isSelf) {
        container.appendChild(avatar);
    }
    addMessageToDisplay(container, isAddToBottom)
}

function displayDownloadElement(message, isAddToBottom) {
    const downLoadElement = document.createElement('a');
    downLoadElement.download = message.content;
    downLoadElement.href = message.file;
    downLoadElement.textContent = "[下载]";
    downLoadElement.classList.add('message-download');
    if (message.userName === userName || message.userId === userId) {
        downLoadElement.classList.add('message-download-me');
    } else {
        downLoadElement.classList.add('message-download-other');
    }
    addMessageToDisplay(downLoadElement, isAddToBottom)
}

function getFileMessageElement(message) {
    const fileElement = document.createElement('a');
    // 让文件在新窗口打开
    fileElement.target = "_blank";
    // 点击downLoadElement下载文件
    fileElement.href = message.file; // Set src to the image field of the message
    fileElement.textContent = message.content;
    fileElement.classList.add('message-bubble');
    fileElement.classList.add('message-bubble-text');
    return fileElement
}

function displayMessage(message, isAddToBottom = true) {
    // console.log("display message", message)
    if (isAddToBottom) {
        insertSendTime(message, isAddToBottom)
    }
    switch (message.type) {
        case "notice":
            displayNoticeElement(message, isAddToBottom)
            if (!pageVisibility) {
                browserNotify(`[${chatRoom}]房间通知`, message.content)
            }
            break;
        case "image":
            displayProfileElement(message, getImageMessageElement(message), isAddToBottom)
            if (!pageVisibility) {
                browserNotify(`[${chatRoom}]${message.userName}`, "发送了一张图片")
            }
            break;
        case "file":
            if (!isAddToBottom) {
                displayDownloadElement(message, isAddToBottom)
            }
            displayProfileElement(message, getFileMessageElement(message), isAddToBottom)
            if (isAddToBottom) {
                displayDownloadElement(message, isAddToBottom)
            }
            if (!pageVisibility) {
                browserNotify(`[${chatRoom}]${message.userName}`, "上传了一个文件")
            }
            break;
        case "text":
            displayProfileElement(message, getNormalMessage(message), isAddToBottom)
            if (!pageVisibility) {
                browserNotify(`[${chatRoom}]${message.userName}`, message.content)
            }
            break;
        default:
            console.log("unknown message type", message)
    }
    if (!isAddToBottom) {
        insertSendTime(message, isAddToBottom)
    }
}

function handleHistoryMessage(message) {
    if (message.type === "text" || message.type === "image" || message.type === "file" || message.type === "notice") {
        displayMessage(message, false)
    }
}

function handleMessage(message) {
    if (message.type === "text" || message.type === "image" || message.type === "file" || message.type === "notice") {
        displayMessage(message)
        if (pageVisibility) {
            if (message.userId === userId || message.userName === userName) {
                scrollToBottom(0)
            } else {
                scrollToBottom()
            }
        } else {
            unreadMessageCount++;
            if (unreadMessageCount >= 99) {
                document.title = `Chat Room - ${chatRoom} (99+)`;
            } else {
                document.title = `Chat Room - ${chatRoom} (${unreadMessageCount})`;
            }
        }
    } else if (message.type === "userCount") {
        document.getElementById("userCount").textContent = "在线用户数：" + message.data
    } else if (message.type === "roomList") {
        handleRoomList(message.data)
    } else if (message.type === "userList") {
        console.log("userList", message.data)
        handleUserList(message.data)
    }
}

function handleUserList(userList) {
    const userElement = document.getElementById("drop-user-list")
    userElement.innerHTML = ""
    for (let i = 0; i < userList.length; i++) {
        const user = userList[i]
        const userDiv = document.createElement("button")
        userDiv.classList.add("user-item")
        userDiv.textContent = user
        userElement.appendChild(userDiv)
        // userElement.appendChild(document.createElement("br"))
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
    let yesterday = new Date();
    yesterday.setDate(yesterday.getDate() - 1); // Set to "yesterday"
    let givenDate = new Date(time);
    if (givenDate.getFullYear() === yesterday.getFullYear() &&
        givenDate.getMonth() === yesterday.getMonth() &&
        givenDate.getDate() === yesterday.getDate()) {
        return "昨天 " + time.substring(11, 16);
    }
    return time.substring(0, 16);
}

var lastTime = ""; // "2024-11-12 00:00:00"
// 判断当前的time是否需要显示时间，如果需要则返回time，否则返回""，并更新lastTime
function getTimeInterval(time, isNotice) {
    var currentDate = new Date(time);
    if (lastTime === "") {
        lastTime = time;
        return getOkTimeText(currentDate, time)
    }
    var lastDate = new Date(lastTime);
    var interval = currentDate - lastDate;
    if (interval > 1000 * 60 * 3) {
        // console.log("interval", interval, lastTime, time)
        lastTime = time;
        return getOkTimeText(currentDate, time);
    }
    // console.log("interval", interval, lastTime, time)
    if (isNotice) {
        return getOkTimeText(currentDate, time);
    }
    return "";
}

function addMessageToDisplay(item, isAddToBottom) {
    if (isAddToBottom) {
        messageDisplay.appendChild(item);
    } else {
        messageDisplay.insertBefore(item, messageDisplay.firstChild);
    }
}

function insertSendTime(message, isAddToBottom) {
    const time = getTimeInterval(message.sendTime, message.type === "notice");
    if (time === "") {
        return
    }
    const timeElement = document.createElement('div');
    timeElement.classList.add("message-time-text")
    timeElement.textContent = time;
    timeElement.classList.add('send-time');
    addMessageToDisplay(timeElement, isAddToBottom);
}

function initSocket() {
    closeSocket()
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
        // console.log("socket open", event)
        document.title = `Chat Room - ${chatRoom}`;
    };
    // 自动重连
    socket.onclose = function (event) {
        if (socket.readyState === WebSocket.CLOSED) {
            // console.log("socket close", event)
            setTimeout(function () {
                showToast("连接断开，尝试重新连接...", {duration: 3000});
                connectToChatRoom(chatRoom, true);
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
                    connectToChatRoom(roomName);
                } else {
                    showToast(data.message);
                }
            });
    });

}

function unLogin() {
    document.cookie = "userId=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
    document.cookie = "userName=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
    goLoginPage()
}

function goLoginPage() {
    window.location.href = "/login";
}

function tryConnectToChatRoom(callback) {
    fetch('/connect', {
        credentials: 'include'  // Include cookies in the request
    })
        .then(response => response.json())
        .then(data => {
            console.log(data);
            if (data.errorCode !== 0) {
                if (data.errorCode === 7) {// 重复登录
                    showToast("账号已超过三个地方登录，请退出其他页面")
                    sleep(1500).then(() => {
                        goLoginPage()
                    })
                } else {
                    showToast(data.message);
                }
            } else {
                showToast("连接成功，欢迎你，" + userName)
                callback()
            }
        });
}

function connectToChatRoomInternal(room) {
    // Check if the userId or room number is empty
    chatRoom = room;
    localStorage.setItem("chatRoom", chatRoom);
    if (userId == null || userId === "" || chatRoom == null || chatRoom === "") {
        showToast('请登录并选择聊天室后再进入');
        sleep(500).then(() => {
                goLoginPage()
            }
        )
        return;
    }
    // clear the message display
    messageDisplay.innerHTML = "";
    // Get the chat history
    getHistoryMessages();
    // Get the room list
    getRoomList();
    // Initialize the WebSocket connection
    initSocket();
}

function connectToChatRoom(room, isRetry = false) {
    messageIndex = 0
    hasMoreMessage = true
    hasLoadedHistory = false
    settingMenuShow = false
    userListMenuShow = false
    // closeSocket()
    tryConnectToChatRoom(function () {
        connectToChatRoomInternal(room)
    })
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

function getImageMessageElement(message) {
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
    $(document).ready(function () {
        // Your code here
        imageElement.addEventListener('click', function () {
            showBigImage(this.src)
        });
    });
    return imageElement
}

function showBigImage(src) {
    $.fancybox.open({
        src: src,
        type: 'image'
    });
}

// Function to display a message
function getNormalMessage(message) {
    const messageElement = document.createElement('div');
    const messageText = document.createElement('pre');
    const content = extractText(message.content)
    for (let part of content) {
        switch (part.type) {
            case "text":
                const textEl = document.createElement('pre');
                textEl.textContent = part.text
                messageElement.appendChild(textEl)
                break
            case "link":
                const linkEl = document.createElement('a');
                console.log("link", part)
                linkEl.innerText = part.url
                linkEl.href = part.url
                linkEl.target = "_blank"
                messageElement.appendChild(linkEl)
                break
        }
    }
    return messageElement
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

function extractText(text) {
    const results = [];
    let regx = /https?:\/\/[^\s\u4e00-\u9fa5]+/g
    let currentPosition = 0;
    while (currentPosition < text.length) {
        const match = regx.exec(text);
        if (match) {
            // 链接前段文字
            let content = text.substring(currentPosition, match.index)
            if (content !== "") {
                results.push({type: 'text', text: content});
            }

            // 判断url是否存在非法字符
            let link = match[0]
            const urlMatch = (/[^a-zA-Z0-9-_.~%&=+\/]*$/g).exec(link)
            if (urlMatch) {
                link = link.substring(0, /[^a-zA-Z0-9-_.~%&=+\/]*$/g.exec(link).index)
            }
            results.push({type: 'link', url: link});
            currentPosition += match.index + link.length
        } else {
            results.push({type: 'text', text: text.substring(currentPosition)});
            break;
        }
    }
    return results;
}