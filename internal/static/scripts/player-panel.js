let socket;
const maxRetries = 1; 
let retryCount = 0;

function createWebSocket() {
    socket = new WebSocket(`ws://${window.location.host}/ws/player`);

    socket.onopen = () => {
        console.log("WebSocket подключен");
        retryCount = 0; 

    socket.onmessage = (event) => {
 
        const data = JSON.parse(event.data);

        if (data.Action === "error") {
            console.log("Ошибка:", data.Message);
            
            if (data.Message === "Admin is not logged in yet") {
                alert("Администратор еще не присоединился")
                window.location.href = "/role"
            } else if (data.Message === "Players Limit Exceeded") {
                alert("Избыток игроков")
                window.location.href = "/role"
            } else {
                alert("Ошибка")
                window.location.href = "/role"
            }
        } else if (data.Action === "message") {
            console.log("Сообщение: ", data.Message)
            alert("Админимстратор добавил вас")
        }
    }

    socket.onclose = () => {
        console.log("WebSocket закрыт.");
    };

    socket.onerror = (error) => {
        console.error("Ошибка WebSocket:", error);
    };
    }
}

function displayMessage(message) {
    const messagesOutput = document.getElementById("messagesOutput");
    const messageElement = document.createElement("div");
    messageElement.textContent = message;
    messagesOutput.appendChild(messageElement);
    messagesOutput.scrollTop = messagesOutput.scrollHeight;
}


document.addEventListener("DOMContentLoaded", () => {
    const sendButton = document.getElementById("sendButton");
    const messageInput = document.getElementById("messageInput");

    sendButton.addEventListener("click", () => {
        const message = messageInput.value.trim();
        if (message) {
            if (socket.readyState === WebSocket.OPEN) {
                socket.send(message); 
                displayMessage(`Вы: ${message}`); 
                messageInput.value = "";
            } else {
                console.error("WebSocket не подключен");
                displayMessage("Сообщение не отправлено: WebSocket не подключен.");
            }
        }
    });

});


createWebSocket();