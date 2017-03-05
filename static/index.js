$(document).ready(function(){
    for(var i = 0; i < a.length; i++){
        function GetFormattedDate(todayTime) {
            var month = (todayTime.getMonth() + 1);
            var day = (todayTime.getDate());
            var year = (todayTime.getFullYear());
            var hour = todayTime.getHours();
            var minute = todayTime.getMinutes();
            return month + "/" + day + "/" + year + " - " + hour + ":" + minute;
        }
        var item = "<div class='comment-box col-md-6 col-xs-offset-3' data-id='1'>\
            <dl>\
                <dt>"+ a[i].username +" | "+ GetFormattedDate(a[i].timestamp) + "</dt>\
                <dd>"+ a[i].message +"</dd>\
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
    username: "Anonymous1",
    timestamp: new Date(),
    message: "Random1"
    },
    {
    username: "Anonymous2",
    timestamp: new Date(),
    message: "Random2"
    },
    {
    username: "Anonymous3",
    timestamp: new Date(),
    message: "Random3"
    }
];
