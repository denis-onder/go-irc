const chat = document.getElementById("chat");
const users = document.getElementById("users");
const input = document.getElementById("message_input");

function generateRandomColor() {
  return `${Math.floor(Math.random() * 256)}, ${Math.floor(
    Math.random() * 256
  )}, ${Math.floor(Math.random() * 256)}`;
}

function generateTimestamp() {
  const d = new Date();
  return [d.getHours(), d.getMinutes(), d.getSeconds()]
    .map((v) => `${v >= 10 ? v : `0${v}`}`)
    .join(":");
}

function generateMessage(time, author, body, color = false) {
  return `<div class="chat_msg ${
    author === "root" ? "chat_msg--root" : false
  }" ${color ? `style="color: rgb(${color});"` : false}>
    <span class="chat_msg_timestamp">${time} | </span>
    <span class="chat_msg_author">[${author}]: </span>
    <span class="chat_msg_body">${body}</span>
   </div>
  `;
}

function writeUser(name, color) {
  users.innerHTML += `
    <div class="users_user">
      <span class="users_user_name">${name}</span>
      <div class="users_user_color" style="color: rgb(${color}); background-color: rgb(${color});" />
    </div>
  `;
}

input.addEventListener("keydown", (e) => {
  if (e.keyCode !== 13) return;
  socket.emit("message_sent", e.target.value);
  e.target.value = "";
});

const socket = io(`http://${window.location.host}`);

const username = `user_${Math.floor(Math.random() * 10)}${Math.floor(
  Math.random() * 10
)}${Math.floor(Math.random() * 10)}`;

socket.on("messages", (json) => {
  const messages = JSON.parse(json);
  if (!messages) return;

  messages.forEach(({ User: user, Body: body }) => {
    chat.innerHTML += generateMessage(
      generateTimestamp(),
      user.Name,
      body,
      user.Name === username ? "255,255,255" : user.Color
    );
  });
});

socket.on("connect", () =>
  socket.emit(
    "new_user",
    JSON.stringify({ Name: username, Color: generateRandomColor() })
  )
);

socket.on(
  "admin",
  (msg) => (chat.innerHTML += generateMessage(generateTimestamp(), "root", msg))
);

socket.on("users", (data) => {
  users.innerHTML = "";
  const usersObj = JSON.parse(data);
  Object.keys(usersObj).forEach((k) => {
    const { Name, Color } = usersObj[k];
    writeUser(Name, Color);
  });
});

socket.on("new_message", (json) => {
  const { User: user, Body: body } = JSON.parse(json);
  chat.innerHTML += generateMessage(
    generateTimestamp(),
    user.Name,
    body,
    user.Name === username ? "255,255,255" : user.Color
  );
});
