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
    var about_us = document.getElementsByName("about_us")[0];
    var contacts = document.getElementsByName("contacts")[0];
    var delivery = document.getElementsByName("delivery")[0];
    about_us.href = "/#about";
    contacts.href = "/#contacts";
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

    var quantity_inputs = document.getElementsByName("quantity");
    var ordersumm = 0;
    var summ_container = document.getElementsByClassName("summ")[0];
    for (let i = 0; i < quantity_inputs.length; i++) {
        var quantity = quantity_inputs[i].value;
        var price = quantity_inputs[i].getAttribute("data-price");
        ordersumm += parseInt(quantity) * parseInt(price);
    }
    summ_container.innerText = ordersumm;
    for (let i = 0; i < quantity_inputs.length; i++) {
        quantity_inputs[i].addEventListener("change", function(event) {
            var ordersumm = 0;
            var summ_container = document.getElementsByClassName("summ")[0];
            for (let i = 0; i < quantity_inputs.length; i++) {
                var quantity = quantity_inputs[i].value;
                var price = quantity_inputs[i].getAttribute("data-price");
                ordersumm += parseInt(quantity) * parseInt(price);
            }
            summ_container.innerText = ordersumm;
        })
    }

    var info = document.getElementsByClassName("info")[0];
    var submit__btn = document.getElementsByClassName("submit__btn")[0];
    submit__btn.addEventListener("click", function() {
        info.classList.add("none");
        info.classList.remove("info_red");
        var firstname = document.getElementById("firstname").value;
        var lastname = document.getElementById("lastname").value;
        var patronymic = document.getElementById("patronymic").value;
        var phone = document.getElementById("phone").value;
        var address = document.getElementById("address").value;
        var city = document.getElementById("city").value;
        var postindex = document.getElementById("postindex").value;
        if (firstname === "" || lastname === "" || patronymic === "" || phone === "" || address === "" || city === "" || postindex === "") {
            info.innerText = "Заполните все поля.";
            info.classList.add("info_red");
            info.classList.remove("none");
        } else {
            var products = [];
            for (let i = 0; i < quantity_inputs.length; i++) {
                var quantity = quantity_inputs[i].value;
                var price = quantity_inputs[i].getAttribute("data-price");
                var id = quantity_inputs[i].getAttribute("data-id");
                products.push({
                    "id": id,
                    "quantity" : parseInt(quantity)
                })
            }
            var  data = {
                "first_name": firstname,
                "last_name": lastname,
                "patronymic": patronymic,
                "phone": phone,
                "city": city,
                "address": address,
                "post_index": postindex,
                "products": products
            };
            console.log(data);

            var json = JSON.stringify(data);

            var request = new XMLHttpRequest();
            request.open("POST", "https://filin-shop.herokuapp.com/api/public/createOrder", false);
            request.setRequestHeader('Content-Type', 'application/json');
            request.send(json);
            if (request.status === 200) {
                var form__card = document.getElementsByClassName("form__card")[0];
                form__card.innerHTML = "";
                form__card.innerText = "Заказ успешно оформлен, информация о заказе отправлены Вам на почту, указанную при регистрации. Менеджер свяжется с Вами для подтверждения данных.";
            } else {
                info.innerText = request.responseText;
                info.classList.add("info_red");
                info.classList.remove("none");
            }
        }
    })
})