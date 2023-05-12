var conn = new WebSocket('localhost:8080');
var chatbox = document.getElementById('chatbox');
var messageForm = document.getElementById('messageForm');

getUsername();

function getUsername() {
    fetch('/api/getUsername')
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
    if (data.recipient) {
        msg.innerText = " ( To" + data.recipient + " ) " + data.sender + ': ' + data.body;
        msg.className = 'message to-command';
    } else if (data.sender) {
        msg.innerText = data.sender + ': ' + data.body;
        msg.className = 'message';
    } else {
        msg.innerText = data.body;
        msg.className = 'message chname-command';
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
        document.getElementById('name').value = currentName;
        var data = {
            name: currentName,
            body: document.getElementById('message').value
        }
        conn.send(JSON.stringify(data));
        document.getElementById('message').value = '';
    } else if (input.startsWith('/to ')) {
        var parts = input.split(' ');
        var recipient = parts[1];
        var message = parts.slice(2).join(' ');

        var data = {
            name: currentName,
            body: message,
            recipient: recipient
        };
        conn.send(JSON.stringify(data));
        document.getElementById('message').value = '';
    } else {
        var data = {
            name: currentName,
            body: input
        };
        conn.send(JSON.stringify(data));
        document.getElementById('message').value = '';
    }
});