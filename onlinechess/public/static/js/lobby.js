window.onload = function(){
    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + document.location.host + "/lobbyws");
        conn.onclose = function (evt) {
           console.log("Connection closed")
        };
        conn.onmessage = function (evt) {
            var resp = evt.data;
            console.log(resp)
            if (resp == "Match Found") {
                window.location.href = "/guest/match";
            }
        };
    } else {
        document.innerHTML = "<b>Your browser does not support WebSockets.</b>";
    }
}