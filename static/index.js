var global_offset = 0;

$(document).ready(function(){

    loadPosts(global_offset);

    $("#moreposts").click(function(){
        loadPosts(global_offset);
    });


});

function loadPosts(offset){
     $.ajax({
         type: "POST",
         url: "/potentiallist",
         data: JSON.stringify({offset: offset}),
         success: function(a){
         console.log(a)
        console.log(a, typeof a)
        a = JSON.parse(a)
        console.log(a, typeof a)
         global_offset += a.length;

         if(a.length < 25){
            $("#moreposts").slideUp();
            $("#nomore").slideUp();
         }
        for(var i = 0; i < a.length; i++){

            var item = "<div class='comment-box col-xs-10 col-xs-offset-1 col-md-8 col-md-offset-2' data-id="+a[i].facebookid+">\
                <dl>\
                    <dt>"+ (a[i].owner.id != "" ? "<a href='http://facebook.com/"+a[i].owner.id+"' target='_blank'>": "") + a[i].owner.name +(a[i].owner.id != ""? "</a>":"") +" | "+ GetFormattedDate(a[i].timestamp) + "</dt>\
                    <dd>"+ a[i].message +"</dd>\
                </dl>\
                <div class='text-right col-md-12'>\
                    <a type'button' class='btn btn-primary' href='http://facebook.com/"+a[i].facebookid+"' target='_blank'>Link</a>\
                    <button type'button' class='btn approved btn-danger'>Confirmed Harassment</button>\
                    <button type='button' class='btn report btn-success'>Not Harassment</button>\
                </div>\
            </div>";
            $("#msgcontainer").append(item);
        };

        $(".approved").click(function(){
            var parent = $(this).parent().parent();
            $.post("/report", JSON.stringify({id:  parent.attr("data-id"), status_name: "confirmed"}), function(){
                parent.hide();
            }, "application/json");

            parent.remove();
        });


        $(".report").click(function(){
            var parent = $(this).parent().parent();
            $.post("/report", JSON.stringify({id:  parent.attr("data-id"), status_name: "discarted"}), function(){
                parent.hide();
            }, "application/json");

            parent.remove();
        });
    },
    datatype: "json"});

}

function GetFormattedDate(time) {
    var todayTime = new Date(time);
    var month = (todayTime.getMonth() + 1);
    var day = (todayTime.getDate());
    var year = (todayTime.getFullYear());
    var hour = todayTime.getHours();
    var minute = todayTime.getMinutes();
    return month + "/" + day + "/" + year + " - " + (hour < 10? "0" + hour : hour) + ":" + (minute < 10? "0" + minute : minute);
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
