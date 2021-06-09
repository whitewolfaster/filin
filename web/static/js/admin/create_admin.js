function LoginValidate(login) {
    if(login.length < 4 || login.length > 20) {
         return false;
        }
 return true;
}

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
        var login_input = document.getElementById("login");
        var login = login_input.value;
        var key_input = document.getElementById("key");
        var key = key_input.value;
        var password_input = document.getElementById("password");
        var password = password_input.value;
        if (login === "") {
            info.innerText = "Email - это обязательное поле.";
            info.classList.add("info_red");
            info.classList.remove("none");
        } else if (password === "") {
            info.innerText = "Пароль - это обязательное поле.";
            info.classList.add("info_red");
            info.classList.remove("none");
        } else if (key === "") {
            info.innerText = "Ключ от сайта - это обязательное поле.";
            info.classList.add("info_red");
            info.classList.remove("none");
        } else if (!LoginValidate(login)) {
            info.innerText = "Некорректный логин.";
            info.classList.add("info_red");
            info.classList.remove("none");
        } else if (!PasswordValidate(password)) {
            info.innerText = "В пароле должно быть от 8 до 30 символов";
            info.classList.add("info_red");
            info.classList.remove("none");
        } else {
            var data = {
                "primary_key": key,
                "login" : login,
                "password" : password,
            }
            var json = JSON.stringify(data);

            var request = new XMLHttpRequest();
            request.open("POST", "https://filin-shop.herokuapp.com/api/private/createAdmin", false);
            request.setRequestHeader('Content-Type', 'application/json');
            request.send(json);
            if (request.status === 201) {
                window.location = "/admin";
            } else {
                info.innerText = request.responseText;
                info.classList.add("info_red");
                info.classList.remove("none");
            }
        }
    })
});