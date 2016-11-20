function getDate() {
    var months = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"];
    var days = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
    var today = new Date();
    var year = today.getFullYear();
    var month = months[today.getMonth()];
    var day = days[today.getDay()];
    var date = today.getDate();
    var dateString = day + ", " + month + " " + date + ", " + year;
    document.getElementById("date").innerHTML = dateString;
}

function getTime() {
    //var setInterval;
    var hour;
    var period;
    setInterval(function() {
        var today = new Date();
        var hours = today.getHours();
        if (hours >= 12) {
            period = "PM";
        } else {
            period = "AM";
        }
        if (hours > 12) {
            hour = hours - 12;
        } else {
            hour = hours;
        }
        var minutes = today.getMinutes();
        var minute = minutes.toString();
        if (minute.length < 2) {
            minute = "0" + minutes;
        }
        var timeString = hour + ":" + minute + " " + period;
        document.getElementById("time").innerHTML = timeString;
    }, 500);
}

function refreshFromHTML() {
    var sec = 1000;
    var min = sec * 60;
    var hr = min * 60;

    setInterval(function() {
        location.reload(false);
        console.log("Refreshed Screen");
    }, 300000);
}

function getQOTD() {
    var file = "../data/qotd.txt";
    var rawFile = new XMLHttpRequest();
    rawFile.open("GET", file, false);
    rawFile.onreadystatechange = function() {
        if (rawFile.readyState === 4) {
            if (rawFile.status === 200 || rawFile.status === 0) {
                var allText = rawFile.responseText;
                alert(allText);
            }
        }
    }
    rawFile.send(null);
}