document.addEventListener("DOMContentLoaded", function() {
    var info = document.getElementsByClassName("info")[0];
    var submit__btn = document.getElementById("submit");
    submit__btn.addEventListener("click", function() {
        info.classList.add("none");
        info.classList.remove("info_red");
        var login_input = document.getElementById("login");
        var login = login_input.value;
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
        } else {
            var data = {
                "login" : login,
                "password" : password,
            }
            var json = JSON.stringify(data);

            var request = new XMLHttpRequest();
            request.open("POST", "https://filin-shop.herokuapp.com/api/private/createAdminSession", false);
            request.setRequestHeader('Content-Type', 'application/json');
            request.send(json);
            if (request.status === 200) {
                window.location = "/admin";
            } else {
                info.innerText = request.responseText;
                info.classList.add("info_red");
                info.classList.remove("none");
            }
        }
    })
});