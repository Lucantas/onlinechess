var current = {piece:""};

window.onload = function () {
    document.querySelectorAll("#chessboard [data-id]").forEach(function(item){
        item.onclick = function(){
        console.log("click")
            if (!current.piece.length){
                current.piece = item.dataset.piece;
                current.id = item.dataset.id;
                current.html = item.innerHTML;
            } else {
                movePiece(this);
            }
        }
    })
    var conn;

    function movePiece(item){
        if (item.innerHTML == undefined){
            item = document.querySelector(`[data-id=${item.next}]`)
        }
        var element = document.querySelector(`[data-id=${current.id}]`);
        element.removeAttribute("data-piece");
        element.innerHTML = "";
        item.setAttribute("data-piece", current.piece);
        item.innerHTML = current.html;
        current.next = item.dataset.id;
        sendPieceInfo(current);
        current = {piece:""};
        console.log("move piece")
    }

    function sendPieceInfo(piece){
        console.log("send piece", piece)
        conn.send(JSON.stringify(piece));
    }
    if (window["WebSocket"]) {
        var matchId = "nothing"
        conn = new WebSocket("ws://" + document.location.host + "/ws?match=" + matchId);
        conn.onclose = function (evt) {
           console.log("Connection closed")
        };
        conn.onmessage = function (evt) {
            var piece = JSON.parse(evt.data);
            console.log(piece.id, piece.piece)
            if (document.querySelector(`[data-id=${piece.next}][data-piece=${piece.piece}]`) != null){
                return;
            }
            console.log("err")
            current = piece;
            movePiece(piece)
        };
    } else {
        var item = document.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
        appendLog(item);
    }
};