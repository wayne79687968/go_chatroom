var backend_url = "localhost:8080";
var conn = new WebSocket('ws://' + backend_url + '/ws');
var chatbox = document.getElementById('chatbox');
var messageForm = document.getElementById('messageForm');

function getUsername() {
    fetch('http://' + backend_url + '/api/getUsername')
        .then(response => response.json())
        .then(data => {
            console.log(data)
            document.getElementById('name').value = data.username;
        });
}

conn.onmessage = function(e) {
    var data = JSON.parse(e.data);
    console.log(data);
    var msg = document.createElement('div');
    if ("getUsername" == data.action) {
        document.getElementById('name').value = data.sender;
        return;
    } else if ("to" == data.action) {
        msg.innerText = " ( To " + data.recipient + " ) " + data.sender + ': ' + data.body;
        msg.className = 'message to-command';
        chatbox.appendChild(msg);
        chatbox.scrollTop = chatbox.scrollHeight;
        return;
    } else if ("chname" == data.action) {
        msg.innerText = data.body;
        msg.className = 'message chname-command';
        if ("" != data.newname) {
            document.getElementById('name').value = data.newname;
        }
    } else {
        msg.innerText = data.sender + ': ' + data.body;
        msg.className = 'message';
    }
    if (document.getElementById('name').value == data.sender) {
        msg.classList.add("my-message");
    } else {
        msg.classList.add("other-message");
    }
    chatbox.appendChild(msg);
    chatbox.scrollTop = chatbox.scrollHeight;
};

messageForm.addEventListener('submit', function(e) {
    e.preventDefault();
    var input = document.getElementById('message').value;
    var currentName = document.getElementById('name').value;

    if (input.startsWith('/chname ')) {
        currentName = input.split(' ')[1];
    }
    var data = {
        sender: currentName,
        body: input
    };
    conn.send(JSON.stringify(data));
    document.getElementById('message').value = '';
});

getUsername();
waitForSocketConnection(conn, function() {
    waitForUsername()
});

function waitForUsername() {
    setTimeout(
        function() {
            if (document.getElementById('name').value != '') {
                console.log("Username is made")
                var data = {
                    sender: document.getElementById('name').value,
                    body: "/getUsername"
                };
                conn.send(JSON.stringify(data));
                console.log(data);
            } else {
                console.log("wait for username...")
                waitForUsername();
            }

        }, 5); // wait 5 milisecond for the connection...
}

function waitForSocketConnection(socket, callback) {
    setTimeout(
        function() {
            if (socket.readyState === 1) {
                console.log("Connection is made")
                if (callback != null) {
                    callback();
                }
            } else {
                console.log("wait for connection...")
                waitForSocketConnection(socket, callback);
            }

        }, 5); // wait 5 milisecond for the connection...
}