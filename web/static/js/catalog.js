function IsLoggedIn() {
    var request = new XMLHttpRequest();
    request.open("GET", "https://filin-shop.herokuapp.com/api/public/isLoggedIn", false);
    request.send();
    var resp = request.response;
    var data = JSON.parse(resp);
    return data.isLoggedIn;
}


var isLoggedIn = IsLoggedIn();

var cards;

document.addEventListener('DOMContentLoaded', function() {
    if (isLoggedIn === "true") {
        var logout__btn = document.getElementsByClassName("logout__btn")[0];
            logout__btn.addEventListener("click", function() {
                var request = new XMLHttpRequest();
                    request.open("GET", "https://filin-shop.herokuapp.com/api/public/deleteSession", false);
                    request.send();
                    console.log(request.status);
                    console.log(request.statusText);
                    window.location= "/";
            })
    }
    var about_us = document.getElementsByName("about_us")[0];
    var contacts = document.getElementsByName("contacts")[0];
    var delivery = document.getElementsByName("delivery")[0];
    about_us.href = "/#about";
    contacts.href = "/#contacts";
    cards = document.getElementsByClassName("book__card");

    var button__show = document.getElementsByClassName("button__show")[0];
    var button__reset = document.getElementsByClassName("button__reset")[0];
    var button__search = document.getElementsByClassName("button__search")[0];

    button__show.addEventListener("click", function() {
        var info__elem = document.getElementsByClassName("info")[0];
        info__elem.classList.add("none");
        for (let i = 0; i < cards.length; i++) {
            if (cards[i].classList.contains("none")) {
                cards[i].classList.remove("none");
            } else {
                continue;
            }
        }
        var filters = "";
        var genre__inputs = document.getElementsByName("genre");

        for (let i = 0; i < genre__inputs.length; i++) {
            if (genre__inputs[i].checked) {
                filters += genre__inputs[i].value + " ";
            }
        }

        filters = filters.trimEnd();

        for (let i = 0; i < cards.length; i++) {
            var card_data_genre = cards[i].getAttribute("data-genres");
            if (card_data_genre.includes(filters)) {
                continue;
            } else {
                cards[i].classList.add("none");
            }
        }

        var count = 0;
        for (let i = 0; i < cards.length; i++) {
            if (cards[i].classList.contains("none")) {
                count++;
            } else {
                continue;
            }
        }
        if (count == cards.length) {
            info__elem.classList.remove("none");
            info__elem.innerText = "По Вашему запросу книг не найдено...";
        }
    });


    button__reset.addEventListener("click", function() {
        document.getElementsByClassName("info")[0].classList.add("none");
        var genre__inputs = document.getElementsByName("genre");

        for (let i = 0; i < genre__inputs.length; i++) {
            if (genre__inputs[i].checked) {
                genre__inputs[i].checked = false;
            }
        }

        for (let i = 0; i < cards.length; i++) {
            if (cards[i].classList.contains("none")) {
                cards[i].classList.remove("none");
            } else {
                continue;
            }
        }
    });


    button__search.addEventListener("click", function() {
        var info__elem = document.getElementsByClassName("info")[0];
        info__elem.classList.add("none");
        for (let i = 0; i < cards.length; i++) {
            if (cards[i].classList.contains("none")) {
                cards[i].classList.remove("none");
            } else {
                continue;
            }
        }

        var search__input = document.getElementsByName("search")[0];
        var search_string = search__input.value.toLowerCase();
        for (let i = 0; i < cards.length; i++) {
            var book__name = cards[i].getElementsByClassName("book__name")[0].innerText.toLowerCase();
            if (book__name.includes(search_string)) {
                continue
            } else {
                cards[i].classList.add("none");
            }
        }


        var count = 0;
        for (let i = 0; i < cards.length; i++) {
            if (cards[i].classList.contains("none")) {
                count++;
            } else {
                continue;
            }
        }
        if (count == cards.length) {
            info__elem.classList.remove("none");
            info__elem.innerText = "По Вашему запросу книг не найдено...";
        }
    });

    var cart_buttons = document.getElementsByClassName("cart__btn");
    for (let i = 0; i < cart_buttons.length; i++) {
        cart_buttons[i].addEventListener("click", function(event) {
            var cart_info = event.target.previousElementSibling;
            if (isLoggedIn === "false") {
                cart_info.innerText = "Вы не зарегистрированы/авторизованы!";
                cart_info.classList.add("cart_info_red");
                cart_info.classList.remove("none");
                setTimeout(
                    () => {
                        cart_info.classList.add("none");
                        cart_info.classList.remove("cart_info_red");
                        cart_info.innerText = "";
                    },
                    3 * 1000
                );
                return;
            }
            var bookID = cart_buttons[i].getAttribute("data-book-id");

            var data = {
                "book_id" : bookID,
            }
            var json = JSON.stringify(data);

            var request = new XMLHttpRequest();
            request.open("POST", "https://filin-shop.herokuapp.com/api/public/addToCart", false);
            request.setRequestHeader('Content-Type', 'application/json');
            request.send(json);
            if (request.status === 201) {
                cart_info.innerText = "Добавлено";
                cart_info.classList.add("cart_info_green");
                cart_info.classList.remove("none");
                setTimeout(
                    () => {
                        cart_info.classList.add("none");
                        cart_info.classList.remove("cart_info_green");
                        cart_info.innerText = "";
                    },
                    3 * 1000
                );
            } else {
                cart_info.innerText = "Уже в корзине";
                cart_info.classList.add("cart_info_green");
                cart_info.classList.remove("none");
                setTimeout(
                    () => {
                        cart_info.classList.add("none");
                        cart_info.classList.remove("cart_info_green");
                        cart_info.innerText = "";
                    },
                    3 * 1000
                );
            }
        })
    }

});