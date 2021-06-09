function IsLoggedIn() {
    var request = new XMLHttpRequest();
    request.open("GET", "https://filin-shop.herokuapp.com/api/public/isLoggedIn", false);
    request.send();
    var resp = request.response;
    var data = JSON.parse(resp);
    return data.isLoggedIn;
}


var isLoggedIn = IsLoggedIn();



document.addEventListener('DOMContentLoaded', function () {
    if (isLoggedIn === "true") {
        var logout__btn = document.getElementsByClassName("logout__btn")[0];
            logout__btn.addEventListener("click", function() {
                var request = new XMLHttpRequest();
                    request.open("GET", "https://filin-shop.herokuapp.com/api/public/deleteSession", false);
                    request.send();
                    window.location= "/";
            })
    }
    var book_of_month_slider = new SimpleAdaptiveSlider('#book_of_month_slider', {
      loop: false,
      autoplay: false,
      interval: 5000,
      swipe: true,
    });
    var lider_slider = new SimpleAdaptiveSlider('#lider_slider', {
      loop: false,
      autoplay: false,
      interval: 5000,
      swipe: true,
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