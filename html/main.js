var cardColors = [
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
    $(".ui.modal.add-command-modal").modal("show")
});

$("#cards-parent").on("click", ".delete", function () {
    $(this).parents(".ui.card").remove()
});

$(".ui.modal.add-command-modal").modal("setting", "transition", "horizontal flip").modal({
    closable: true,
    onApprove: function () {
        // If title empty, return false
        let parentContent = $(this).children(".content");
        let inputHeader = parentContent.find("input[name ='header']");
        if (inputHeader.val() === "")
        {
            inputHeader.focus();
            return false;
        }

        // Else call method below
        submitModalNewCommand(parentContent);
    },
    onHide: function () {
        $(this).find("form").form("clear");
    },
});

$(".ui.modal.add-command-modal .content form").form({
    on: 'blur',
    inline: false,
    delay: false,
    fields: {
        header: {
            identifier: 'header',
            rules: [{
                type: 'empty',
                prompt: 'Please enter a title'
            }]
        }
    }
});

function submitModalNewCommand(parentContent) {
    // Create new ui card
    let newcard = "\n" +
        "                        <div class=\"ui " + cardColors[Math.floor(Math.random() * cardColors.length)] + " card\">\n" +
        "                            <div class=\"content\">\n" +
        "                                <div class=\"header\">\n" +
        "                                    <div class=\"ui transparent fluid input header\">\n" +
        "                                        <input type=\"text\" placeholder=\"Title ...\" value=\"" + parentContent.find("input[name ='header']").val() + "\">\n" +
        "                                    </div>\n" +
        "                                </div>\n" +
        "                                <div class=\"meta\">\n" +
        "                                    <div class=\"ui transparent fluid input meta\">\n" +
        "                                        <input type=\"text\" placeholder=\"Desription ...\" value=\"" + parentContent.find("input[name ='description']").val() + "\">\n" +
        "                                    </div>\n" +
        "                                </div>\n" +
        "                            </div>\n" +
        "                            <div class=\"extra content\">\n" +
        "                                <div class=\"description\">\n" +
        "                                    <div class=\"ui transparent fluid input\">\n" +
        "                                        <input type=\"text\" placeholder=\"Your command ...\" value=\"" + parentContent.find("input[name ='command']").val() + "\">\n" +
        "                                        <div class=\"ui left pointing label olive\" style=\"display: none\">Keep it short\n" +
        "                                        </div>\n" +
        "                                    </div>\n" +
        "                                </div>\n" +
        "                            </div>\n" +
        "                            <div class=\"extra content\">\n" +
        "                                <div class=\"description\">\n" +
        "                                    <div class=\"ui transparent fluid input\">\n" +
        "                                        <input type=\"text\" placeholder=\"Output text ...\" value=\"" + parentContent.find("input[name ='output']").val() + "\">\n" +
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
        existingDiv.after(newcard)
    } else {
        $("#cards-parent").append(newcard)
    }
}