$(".ui.card").hover(
    function () {
        $(this).addClass("raised");
    }, function () {
        $(this).removeClass("raised");
    }
);

$(".ui.transparent.fluid.input").focusin(
    function () {
        $(this).children("div").show(100);
    }
).focusout(
    function () {
        $(this).children("div").hide(100);
    }
);

