document.addEventListener("DOMContentLoaded", function() {
    var logout__btn = document.getElementsByClassName("logout__btn")[0];
    logout__btn.addEventListener("click", function() {
        var request = new XMLHttpRequest();
            request.open("GET", "https://filin-shop.herokuapp.com/api/private/deleteAdminSession", false);
            request.send();
            window.location= "/admin/login";
    })
    var submit__btn = document.getElementById("submit");
    var id = submit__btn.getAttribute("data-id");
    var info = document.getElementsByClassName("info")[0];
    submit__btn.addEventListener("click", function() {
        info.classList.add("none");
        info.classList.remove("info_red");
        var name = document.getElementsByName("name")[0];
        var author = document.getElementsByName("author")[0];
        var year = document.getElementsByName("year")[0];
        var genre = "";
        var pub_house = document.getElementsByName("pub_house")[0];
        var description = document.getElementsByName("description")[0];
        var price = document.getElementsByName("price")[0];
        var genre__inputs = document.getElementsByName("genre");
        for (let i = 0; i < genre__inputs.length; i++) {
            if (genre__inputs[i].checked) {
                genre += genre__inputs[i].value + " ";
            }
        }
        genre = genre.trimEnd();

        var book = {
            id: id,
            name: name.value,
            author: author.value,
            year: parseInt(year.value, 10),
            genre: genre,
            pub_house: pub_house.value,
            description: description.value,
            price: parseInt(price.value, 10),
        };

        var file__form = document.getElementById("file_form");
        var data = new FormData(file__form);

        var json = JSON.stringify(book);

        data.append("json", json);

        var request = new XMLHttpRequest();
        request.open("POST", "https://filin-shop.herokuapp.com/api/private/updateBook", false);
        request.send(data);
        if (request.status === 200) {
            window.location = "/admin/book_list";
        } else {
            info.innerText = request.responseText;
            info.classList.add("info_red");
            info.classList.remove("none");
        }
    })
});