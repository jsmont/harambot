$(document).ready(function(){
    for(var i = 0; i < a.length; i++){
        var item = "<div class='comment-box col-md-6 col-xs-offset-3' data-id='1'>\
            <dl>\
                <dt>Nom | Id | HH:MI</dt>\
                <dd>comentari</dd>\
            </dl>\
            <div class='buttons col-md-offset-9'>\
                <button type'button' class='btn approved btn-success'>Approved</button>\
                <button type='button' class='btn report btn-danger'>Report</button>\
            </div>\
        </div>";
        $("#msgcontainer").append(item);
    };

    $(".approved").click(function(){
        $(this).parent().parent().hide();
    });

    $(".report").click(function(){
        $(this).parent().parent().hide();
    });

});

a = [
    {
    username: "Anonymous",
    timestamp: new Date(),
    message: "Random"
    },
    {
    username: "Anonymous",
    timestamp: new Date(),
    message: "Random"
    },
    {
    username: "Anonymous",
    timestamp: new Date(),
    message: "Random"
    }
];
