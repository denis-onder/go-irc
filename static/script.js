const output = document.getElementById("chat");

function generateRandomColor() {
  return `${Math.floor(Math.random() * 256)}, ${Math.floor(
    Math.random() * 256
  )}, ${Math.floor(Math.random() * 256)}`;
}

function generateTimestamp() {
  const d = new Date();
  return [d.getHours(), d.getMinutes(), d.getSeconds()]
    .map((v) => `${v > 10 ? v : `0${v}`}`)
    .join(":");
}

function generateMessage(time, author, body) {
  return `<div class="chat_msg ${
    author === "admin" ? "chat_msg--admin" : false
  }">
    <span class="chat_msg_timestamp">${time} | </span>
    <span class="chat_msg_author">[${author}]: </span>
    <span class="chat_msg_body">${body}</span>
   </div>
  `;
}

const socket = io(`http://${window.location.host}`);

const username = `test_user${Math.floor(Math.random() * 10)}${Math.floor(
  Math.random() * 10
)}${Math.floor(Math.random() * 10)}`;

socket.on("connect", function () {
  socket.emit(
    "new_user",
    JSON.stringify({ Name: username, Color: generateRandomColor() })
  );
});

socket.on("user_joined", (msg) => {
  output.innerHTML += generateMessage(generateTimestamp(), "admin", msg);
});

socket.on("user_left", (msg) => {
  output.innerHTML += generateMessage(generateTimestamp(), "admin", msg);
});

socket.on("new_user", (msg) => console.log(msg));
