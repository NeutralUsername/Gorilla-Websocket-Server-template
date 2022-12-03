var ws = new WebSocket("ws://localhost:8080/ws")
const MSG_DELIMITER = "<;>"


export class Index extends React.Component {
	constructor(props) {
		super(props)
		this.state = {}
        this.messageHandler = this.messageHandler.bind(this)
        ws.onmessage = (event) => {
            this.messageHandler(event.data)
        };
    }   
    messageHandler(message) {
        console.log(message)
        let data = message.split(MSG_DELIMITER);
        if(data[0] == "request_credentials") {
            ws.send("credentials"+MSG_DELIMITER+cookie("id")+ MSG_DELIMITER + cookie("name")+ MSG_DELIMITER +cookie("password"));
        }
        if(data[0] == "user_data") {
            document.cookie = "id=" + data[1];
            document.cookie = "name=" + data[2];
            document.cookie = "password=" + data[3];
            this.setState({
                id: data[1],
                name: data[2],
                password: data[3],
                email :     data[4],
                secretQuestion : data[5],
                secretAnswer : data[6],
                
                power : Number(data[7]),
                lastActive : new Date(data[8]),
                rating : Number(data[9]),
            })
        }
    }
    render() {
        return React.createElement("div", {}, "Hello "+ this.state.name)
    }
}

function start_app() {
    ReactDOM.render(
        React.createElement(Index), document.getElementById('root')
    )
  }
  start_app();


function cookie(name) {
    var value = "; " + document.cookie;
    var parts = value.split("; " + name + "=");
    if (parts.length == 2) return parts.pop().split(";").shift();
  }