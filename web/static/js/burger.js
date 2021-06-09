var burger = document.getElementsByClassName("burger_button");
var header_navigation = document.getElementsByClassName("naw__row");
function test() {
    burger[0].classList.toggle("active");
    header_navigation[0].classList.toggle("active");
}
burger[0].addEventListener("click", test);