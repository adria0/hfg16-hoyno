(function() {

    if (window["WebSocket"]) {
        var loc = window.location, wsUrl;
        if (loc.protocol === "https:") {
                 wsUrl = "wss:";

        } else {
                new_uri = "ws:";

        }
        wsUrl += "//" + loc.host;
        wsUrl += loc.pathname + "ws";
        conn = new WebSocket(wsUrl);
        conn.onclose = function(evt) {
            appendLog("<div><b>Connection closed.</b></div>")
        }
        conn.onmessage = function(evt) {
            appendLog(evt.data)
        }
    } else {
        appendLog("<div><b>Your browser does not support WebSockets.</b></div>")
    }

    function submit() {
        msg = $("#input").val()

        if (!conn) {
            return false;
        }
        conn.send(msg);
        $("#input").val("")     
    }

    $( "#input" ).keypress(function( event ) {
      if ( event.which == 13 ) {
         event.preventDefault();
         submit()
      }
    });

    $("#button").click(function() {
        submit()
        return false
    });

    function appendLog(msg) {
        output.innerHTML = '<p><span>' +  msg + '</span></p>' + output.innerHTML;
    }

    $("#configure").click(function() {
        window.open("/config");
    });


})();
