$(document).ready(function(){

    $.get("/potentiallist", function(a){
        console.log(a, typeof a)
        a = JSON.parse(a)
        console.log(a, typeof a)
        for(var i = 0; i < a.length; i++){

            var item = "<div class='comment-box col-xs-10 col-xs-offset-1 col-md-8 col-md-offset-2'>\
                <dl>\
                    <dt>"+ a[i].username +" | "+ GetFormattedDate(a[i].timestamp) + "</dt>\
                    <dd>"+ a[i].message +"</dd>\
                </dl>\
                <div class='text-right col-md-12'>\
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

});
function GetFormattedDate(time) {
    var todayTime = new Date(time);
    var month = (todayTime.getMonth() + 1);
    var day = (todayTime.getDate());
    var year = (todayTime.getFullYear());
    var hour = todayTime.getHours();
    var minute = todayTime.getMinutes();
    return month + "/" + day + "/" + year + " - " + hour + ":" + minute;
}
/*
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
];*/
