function PasswordValidate(password) {
    if(password.length < 8 || password.length > 30) {
         return false;
        }

 return true;
}


document.addEventListener("DOMContentLoaded", function() {
    var logout__btn = document.getElementsByClassName("logout__btn")[0];
    logout__btn.addEventListener("click", function() {
        var request = new XMLHttpRequest();
            request.open("GET", "https://filin-shop.herokuapp.com/api/private/deleteAdminSession", false);
            request.send();
            window.location= "/admin/login";
    })
    var info = document.getElementsByClassName("info")[0];
    var submit__btn = document.getElementById("submit");
    submit__btn.addEventListener("click", function() {
        info.classList.add("none");
        info.classList.remove("info_red");
        var old_input = document.getElementById("old");
        var old = old_input.value;
        var new_input = document.getElementById("new");
        var newpswd = new_input.value;
        if (old === "") {
            info.innerText = "Старый пароль - это обязательное поле.";
            info.classList.add("info_red");
            info.classList.remove("none");
        } else if (newpswd === "") {
            info.innerText = "Новый пароль - это обязательное поле.";
            info.classList.add("info_red");
            info.classList.remove("none");
        } else if (!PasswordValidate(newpswd)) {
            info.innerText = "В новом пароле должно быть от 8 до 30 символов";
            info.classList.add("info_red");
            info.classList.remove("none");
        } else {
            var data = {
                "old_password" : old,
                "new_password" : newpswd,
            }
            var json = JSON.stringify(data);

            var request = new XMLHttpRequest();
            request.open("POST", "https://filin-shop.herokuapp.com/api/private/changeAdminPassword", false);
            request.setRequestHeader('Content-Type', 'application/json');
            request.send(json);
            if (request.status === 200) {
                info.innerText = "Пароль успешно изменен.";
                info.classList.add("info_green");
                info.classList.remove("none");
            } else {
                info.innerText = request.responseText;
                info.classList.add("info_red");
                info.classList.remove("none");
            }
        }
    })
});