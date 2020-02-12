const cardColors = [
    "red",
    "orange",
    "yellow",
    "olive",
    "green",
    "teal",
    "blue",
    "violet",
    "purple",
    "pink",
    "brown",
    "grey",
    "black"
];

let cardsParent = $("#cards-parent");
let modalUi = $(".ui.modal.add-command-modal");

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

$(".sortable").sortable({
    placeholder: "card"
});

$(".add-command").click(function () {
    $(".ui.modal.add-command-modal").modal("show");
});

cardsParent.on("click", ".delete", function () {
    $(this).parents(".ui.card").remove();
});

function submitModalNewCommand(title, description, command, output) {
    // Create new ui card
    let newcard = "\n" +
        "                        <div class=\"ui " + cardColors[Math.floor(Math.random() * cardColors.length)] + " card\">\n" +
        "                            <div class=\"content\">\n" +
        "                                <div class=\"header\">\n" +
        "                                    <div class=\"ui transparent fluid input header\">\n" +
        "                                        <input type=\"text\" placeholder=\"Title ...\" value=\"" + title + "\">\n" +
        "                                    </div>\n" +
        "                                </div>\n" +
        "                                <div class=\"meta\">\n" +
        "                                    <div class=\"ui transparent fluid input meta\">\n" +
        "                                        <input type=\"text\" placeholder=\"Desription ...\" value=\"" + description + "\">\n" +
        "                                    </div>\n" +
        "                                </div>\n" +
        "                            </div>\n" +
        "                            <div class=\"extra content\">\n" +
        "                                <div class=\"description\">\n" +
        "                                    <div class=\"ui transparent fluid input\">\n" +
        "                                        <input type=\"text\" placeholder=\"Your command ...\" value=\"" + command + "\">\n" +
        "                                        <div class=\"ui left pointing label olive\" style=\"display: none\">Keep it short\n" +
        "                                        </div>\n" +
        "                                    </div>\n" +
        "                                </div>\n" +
        "                            </div>\n" +
        "                            <div class=\"extra content\">\n" +
        "                                <div class=\"description\">\n" +
        "                                    <div class=\"ui transparent fluid input\">\n" +
        "                                        <input type=\"text\" placeholder=\"Output text ...\" value=\"" + output + "\">\n" +
        "                                    </div>\n" +
        "                                </div>\n" +
        "                            </div>\n" +
        "                            <div class=\"extra content\">\n" +
        "                                <div class=\"ui two buttons\">\n" +
        "                                    <div class=\"ui basic grey button disable\">Disable</div>\n" +
        "                                    <div class=\"ui basic red button delete\">Delete</div>\n" +
        "                                </div>\n" +
        "                            </div>\n" +
        "                        </div>";

    let existingDiv = $("#cards-parent > div").last();
    if (existingDiv.length) {
        existingDiv.after(newcard);
    } else {
        cardsParent.append(newcard);
    }

    if (modalUi.modal("is active")) {
        modalUi.modal("hide")
    }
}

modalUi.modal("setting", "transition", "horizontal flip").modal({
    closable: true,
    onHide() {
        $(this).find("form").form("clear");
    },
});

$("#form-modal").submit(function (e) {
    e.preventDefault();
    console.log("in submit");
    // If title empty, return false
    let inputTitle = $(this).find("input[name ='title']");
    if (inputTitle.val() === "") {
        $(this).focus();
        return false;
    }

    // Else call method below
    submitModalNewCommand(inputTitle.val(),
        $(this).find("input[name ='description']").val(),
        $(this).find("input[name ='command']").val(),
        $(this).find("input[name ='output']").val()
    );
});

$(".ui.modal.add-command-modal .content form").form({
    on: "blur",
    inline: false,
    delay: false,
    fields: {
        title: {
            identifier: "title",
            rules: [{
                type: "empty",
                prompt: "Please enter a title"
            }]
        }
    }
});

