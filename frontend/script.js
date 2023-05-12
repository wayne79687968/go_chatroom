var conn = new WebSocket('localhost:8080');
var chatbox = document.getElementById('chatbox');
var messageForm = document.getElementById('messageForm');

conn.onmessage = function(e) {
    var data = JSON.parse(e.data);
    console.log(data);
    var msg = document.createElement('div');
    msg.innerText = data.sender + ': ' + data.body;
    chatbox.appendChild(msg);
    chatbox.scrollTop = chatbox.scrollHeight;
};

messageForm.addEventListener('submit', function(e) {
    e.preventDefault();
    var input = document.getElementById('message').value;
    var currentName = document.getElementById('name').value;
    var data = {
        name: currentName,
        body: input
    };
    conn.send(JSON.stringify(data));
    document.getElementById('message').value = '';
});