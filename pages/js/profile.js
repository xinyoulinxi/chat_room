function saveProfile() {
    var originalName = document.getElementById('usernameInput').value;
    var originalPassword = document.getElementById('userPasswordInput').value;
    // var newName = document.getElementById('newUsernameInput').value;
    var newPassword = document.getElementById('newUserPasswordInput').value;
    var newPasswordConfirm = document.getElementById('newUserPasswordInputConfirm').value;
    if (originalName === '' || originalPassword === '') {
        showToast("用户名和密码不能为空，请重新输入");
        return;
    }
    if (
        // newName === '' ||
        newPassword === ''
    ) {
        showToast("新的密码不能为空，请重新输入");
        return;
    }
    if (
        // newName === originalName &&
        newPassword === originalPassword) {
        showToast("密码未改变，请重新输入");
        return;
    }

    if (newPassword.length < 6) {
        showToast("密码长度不能小于6位，请重新输入");
        return;
    }

    if (newPassword !== newPasswordConfirm) {
        showToast("两次输入的密码不一致，请重新输入");
        return;
    }

    // console.log("originalName:" + originalName + " originalPassword:" + originalPassword)
    // console.log("newName:" + newName + " newPassword:" + newPassword)
    userId = getCookie("userId");
    // 访问服务器，通过login接口验证用户名和密码
    fetch('/update_profile', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            id: userId,
            username: originalName,
            password: originalPassword,
            // newUsername: newName,
            newPassword: newPassword
        })
    })
        .then(response => response.json())
        .then(data => {
            console.log(data)
            if (data.errorCode !== 0) {
                showToast(data.message);
            } else {
                // localStorage.setItem("username", newName);
                localStorage.setItem("password", newPassword);
                showToast("更新成功，将自动跳转到聊天室")
                // 更新成功，跳转到聊天室页面
                chatRoom = localStorage.getItem("chatRoom");
                sleep(1000).then(() => {
                    goToChatRoom(userId, originalName, chatRoom);
                })
            }
        })
        .catch(error => console.error('Error:', error));
}

