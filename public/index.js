var ws = new WebSocket("ws://"+ window.location.hostname +":8080/ws")

class Index extends React.Component {
	constructor(props) {
		super(props)
		this.state = {}
        ws.onmessage = (event) => {
            if(event.data == "ping")
                ws.send("pong")
        };
    }   
    render() {
        return React.createElement("div", {}, "hello world")
    }
}

function start_app() {
    ReactDOM.render(
        React.createElement(Index), document.getElementById('root')
    )
  }
start_app();