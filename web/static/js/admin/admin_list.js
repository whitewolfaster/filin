document.addEventListener("DOMContentLoaded", function() {
    var logout__btn = document.getElementsByClassName("logout__btn")[0];
    logout__btn.addEventListener("click", function() {
        var request = new XMLHttpRequest();
            request.open("GET", "https://filin-shop.herokuapp.com/api/private/deleteAdminSession", false);
            request.send();
            window.location= "/admin/login";
    })
    var info = document.getElementsByClassName("info")[0];
    var delete__btn = document.getElementsByClassName("delete__btn");
    for (let i = 0; i < delete__btn.length; i++) {
        delete__btn[i].addEventListener("click", function(event) {
            info.classList.add("none");
            info.classList.remove("info_red");
            var target = event.target;
            var id = target.getAttribute("data-id");
            var key = target.previousElementSibling.value;
            if (key === "") {
                info.innerText = "Введите ключ от сайта.";
                info.classList.add("info_red");
                info.classList.remove("none");
            } else {
                var data = {
                    "primary_key" : key,
                    "admin_id" : id,
                }
                var json = JSON.stringify(data);
    
                var request = new XMLHttpRequest();
                request.open("POST", "https://filin-shop.herokuapp.com/api/private/deleteAdmin", false);
                request.setRequestHeader('Content-Type', 'application/json');
                request.send(json);
                if (request.status === 200) {
                    window.location = "/admin/admin_list";
                } else {
                    info.innerText = request.responseText;
                    info.classList.add("info_red");
                    info.classList.remove("none");
                }
            }
        })
    }
});