var ws = new WebSocket("ws://"+ window.location.hostname +":8080/ws")
const MSG_DELIMITER = "<;>"

class Index extends React.Component {
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
            let userData = JSON.parse(data[1])
            document.cookie = "id="+userData.id
            document.cookie = "name="+userData.name
            document.cookie = "password="+userData.password
            this.setState(userData)
        }
    }
    render() {
        if(this.state.id)
            return React.createElement("div", {}, "Hello "+ this.state.name)
        else return React.createElement("div", {}, "authenticating...")
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