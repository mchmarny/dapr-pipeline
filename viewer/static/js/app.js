window.onload = function () {

    console.log("Protocol: " + location.protocol);
    var wsURL = "ws://" + document.location.host + "/ws"
    if (location.protocol == 'https:') {
        wsURL = "wss://" + document.location.host + "/ws"
    }
    console.log("WS URL: " + wsURL);

    var log = document.getElementById("tweets");

    function appendLog(item) {
        var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
        log.appendChild(item);
        if (doScroll) {
            log.scrollTop = log.scrollHeight - log.clientHeight;
        }

    }

    if (log) {

        sock = new WebSocket(wsURL);

        var connDiv = document.getElementById("connection-status");
        connDiv.innerText = "closed";

        sock.onopen = function () {
            console.log("connected to " + wsURL);
            connDiv.innerText = "open";
        };

        sock.onclose = function (e) {
            console.log("connection closed (" + e.code + ")");
            connDiv.innerText = "closed";
        };

        /*
        sample
        {
            "content":"some text",
            "sentiment":1,
            "published":"2020-04-13T00:21:09Z",
            "id":124949278000000000,
            "query":"the term that was used to search",
            "author":"twitterusername",
            "author_pic": "https://adasdsa"
        }
        */

        sock.onmessage = function (e) {
            console.log(e);
            var t = JSON.parse(e.data);
            console.log(t);
            var item = document.createElement("div");
            item.className = "item";
            var tmsg = "<img src='" + t.author_pic + "' class='profile-pic' />" +
                "<div class='item-text'><b><img src='static/img/s" + t.sentiment +
                ".svg' alt='sentiment' class='sentiment' />" + t.author +
                "<a href='https://twitter.com/" + t.author + "/status/" + t.id +
                "' target='_blank'><img src='static/img/tw.svg' class='tweet-link' /></a></b>" +
                "<br /><i>" + t.content + "</i></div>";
            item.innerHTML = tmsg
            appendLog(item);
        };

    } // if log


};