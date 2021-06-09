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
            var genre_name = target.getAttribute("data-genre");
            var data = {
                "name": genre_name
            }
            var json = JSON.stringify(data);
    
            var request = new XMLHttpRequest();
            request.open("POST", "https://filin-shop.herokuapp.com/api/private/deleteGenre", false);
            request.setRequestHeader('Content-Type', 'application/json');
            request.send(json);
            if (request.status === 200) {
                window.location = "/admin/genres_list";
            } else {
                info.innerText = request.responseText;
                info.classList.add("info_red");
                info.classList.remove("none");
            }
        })
    }
    var genre_name_input = document.getElementsByName("name")[0];
    var add__btn = document.getElementsByClassName("add__btn")[0];
    add__btn.addEventListener("click", function() {
        info.classList.add("none");
        info.classList.remove("info_red");
        var name = genre_name_input.value;
        if (name === "") {
            info.innerText = "Введите название жанра.";
            info.classList.add("info_red");
            info.classList.remove("none");
        } else {
            var data = {
                "name": name
            }
            var json = JSON.stringify(data);

            var request = new XMLHttpRequest();
            request.open("POST", "https://filin-shop.herokuapp.com/api/private/createGenre", false);
            request.setRequestHeader('Content-Type', 'application/json');
            request.send(json);
            if (request.status === 201) {
                window.location = "/admin/genres_list";
            } else {
                info.innerText = request.responseText;
                info.classList.add("info_red");
                info.classList.remove("none");
            }
        }
    })

});