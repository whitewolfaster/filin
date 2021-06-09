document.addEventListener('DOMContentLoaded', function () {
    var logout__btn = document.getElementsByClassName("logout__btn")[0];
    logout__btn.addEventListener("click", function() {
        var request = new XMLHttpRequest();
            request.open("GET", "https://filin-shop.herokuapp.com/api/private/deleteAdminSession", false);
            request.send();
            window.location= "/admin/login";
    })
})