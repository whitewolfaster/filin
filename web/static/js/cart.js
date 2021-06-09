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
    var empty_cart = document.getElementsByClassName("empty")[0];
    if (!empty_cart.classList.contains("none")) {
        empty_cart.classList.add("none");
    }
    var about_us = document.getElementsByName("about_us")[0];
    var contacts = document.getElementsByName("contacts")[0];
    var delivery = document.getElementsByName("delivery")[0];
    about_us.href = "/#about";
    contacts.href = "/#contacts";
    cards = document.getElementsByClassName("book__card");
    var clear__btns = document.getElementsByClassName("clear__btn");
    var order__btns = document.getElementsByClassName("order__btn");

    if (cards.length < 1) {
        for (let i = 0; i < clear__btns.length; i++) {
            clear__btns[i].classList.add("none");
        }
        for (let i = 0; i < order__btns.length; i++) {
            order__btns[i].classList.add("none");
        }
        empty_cart.classList.remove("none");
    } else {
        for (let i = 0; i < clear__btns.length; i++) {
            clear__btns[i].addEventListener("click", function() {
                var request = new XMLHttpRequest();
                request.open("GET", "https://filin-shop.herokuapp.com/api/public/cleanCart", false);
                request.send();
                window.location = "/cart";
            })
        }
        for (let i = 0; i < order__btns.length; i++) {
            order__btns[i].addEventListener("click", function() {
                window.location = "/order";
            })
        }
        var cart_buttons = document.getElementsByClassName("cart__btn");
        for (let i = 0; i < cart_buttons.length; i++) {
            cart_buttons[i].addEventListener("click", function(event) {
                var cart_info = event.target.previousElementSibling;
                if (isLoggedIn === "false") {
                    cart_info.innerText = "Вы не зарегистрированы!";
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
                request.open("POST", "https://filin-shop.herokuapp.com/api/public/deleteFromCart", false);
                request.setRequestHeader('Content-Type', 'application/json');
                request.send(json);
                window.location = "/cart";
            })
        }
    }
})