function EmailValidate(email) {
    
    if(/^([a-z0-9_-]+\.)*[a-z0-9_-]+@[a-z0-9_-]+(\.[a-z0-9_-]+)*\.[a-z]{2,6}$/.test(email) === false) {
        return false;
    }

    if(email.length < 6 || email.length > 30) {
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
    var about_us = document.getElementsByName("about_us")[0];
    var contacts = document.getElementsByName("contacts")[0];
    var delivery = document.getElementsByName("delivery")[0];
    about_us.href = "/#about";
    contacts.href = "/#contacts";
    var info = document.getElementsByClassName("info")[0];
    var submit__btn = document.getElementById("submit");
    submit__btn.addEventListener("click", function() {
        info.classList.add("none");
        info.classList.remove("info_red");
        var email_input = document.getElementById("email");
        var email = email_input.value;
        var password_input = document.getElementById("password");
        var password = password_input.value;
        if (email === "") {
            info.innerText = "Email - это обязательное поле.";
            info.classList.add("info_red");
            info.classList.remove("none");
        } else if (password === "") {
            info.innerText = "Пароль - это обязательное поле.";
            info.classList.add("info_red");
            info.classList.remove("none");
        } else if (!EmailValidate(email)) {
            info.innerText = "Некорректный e-mail";
            info.classList.add("info_red");
            info.classList.remove("none");
        } else if (!PasswordValidate(password)) {
            info.innerText = "В пароле должен быть от 8 до 30 символов";
            info.classList.add("info_red");
            info.classList.remove("none");
        } else {
            var data = {
                "email" : email,
                "password" : password,
            }
            var json = JSON.stringify(data);

            var request = new XMLHttpRequest();
            request.open("POST", "https://filin-shop.herokuapp.com/api/public/createSession", false);
            request.setRequestHeader('Content-Type', 'application/json');
            request.send(json);
            if (request.status === 200) {
                window.location = "/";
            } else {
                info.innerText = request.responseText;
                info.classList.add("info_red");
                info.classList.remove("none");
            }
        }
    })
});