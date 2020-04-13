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

        sock.onopen = function () {
            console.log("connected to " + wsURL);
        };

        sock.onclose = function (e) {
            console.log("connection closed (" + e.code + ")");
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
            var tmsg = "<img src='" + t.author_pic + "'/><div class='item-text'><b>" +
                t.author + ":</b> at " + t.published + "<br /><i>" + t.content + "</i>" +
                "<img src='static/img/s" + t.sentiment + ".png' class='sentiment'/></div>";
            item.innerHTML = tmsg
            appendLog(item);
        };

    } // if log


};