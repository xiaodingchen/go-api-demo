package controllers

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"html/template"
	"log"
	"net/http"
	"test.local/internal/service"
	"test.local/pkg/utils"
)

type WebSocket struct {
}

var upgrader = websocket.Upgrader{}

func NewWs() *WebSocket {
	return &WebSocket{}
}

func (ws *WebSocket) Echo(ctx *gin.Context) {
	c, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		utils.Log().Error("controllers.WebSocket.index err", zap.Error(err))
		return
	}
	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			utils.Log().Error("controllers.WebSocket.index err", zap.Error(err))
			break
		}
		buf := bytes.NewBuffer(message)
		buf.WriteString(":wwwwssss")
		err = c.WriteMessage(mt, buf.Bytes())
		if err != nil {
			utils.Log().Error("controllers.WebSocket.index err", zap.Error(err))
			break
		}
	}
}

func (ws *WebSocket) Index(ctx *gin.Context) {
	homeTemplate.Execute(ctx.Writer, "ws://"+ctx.Request.Host+"/ws/echo")
}

func (ws *WebSocket) Chat(ctx *gin.Context) {
	http.ServeFile(ctx.Writer, ctx.Request, utils.TemplateFile("websocket/chat.html"))
}

func (ws *WebSocket) ChatWs(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		utils.Log().Error("controllers.WebSocket.ChatMsg err", zap.Error(err))
		return
	}
	log.Println("websocket conn")
	client := service.NewClient()
	client.Start(conn)
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))
