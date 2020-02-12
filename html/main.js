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
let modalForm = $(".ui.modal.add-command-modal #form-modal");

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
        modalUi.modal("hide");
    }
}

modalUi.modal("setting", "transition", "horizontal flip").modal({
    closable: true,
    onApprove() {
        return false;
    },
    onShow() {
        modalForm.form("clear");
    }
});

function validateModal(form) {
    // If title empty, return false
    let inputTitle = form.find("input[name ='title']");
    if (inputTitle.val() === "") {
        form.focus();
        return false;
    }

    // Else call method below
    submitModalNewCommand(inputTitle.val(),
        form.find("input[name ='description']").val(),
        form.find("input[name ='command']").val(),
        form.find("input[name ='output']").val()
    );
}

$("#form-modal").submit(function (e) {
    validateModal($(this));

    return false;
});

modalUi.find(".ok").click(function () {
    validateModal(modalForm);
});

modalForm.form({
    on: "blur",
    inline: false,
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

