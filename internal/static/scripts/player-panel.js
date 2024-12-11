let socket;
const maxRetries = 1; // Максимальное количество попыток переподключения
let retryCount = 0;

function createWebSocket() {
    socket = new WebSocket(`ws://${window.location.host}/ws/player`);

    socket.onopen = () => {
        console.log("WebSocket подключен");
        retryCount = 0; // Сбросить счетчик попыток при успешном подключении
    };

    socket.onmessage = (event) => {
 
        const data = JSON.parse(event.data);

        if (data.action === "players_exceeded") {
            console.log("Ошибка: ", data.message)
            socket.onerror(new Error(data.message))
        } else if (data.action === "internal_error") {
            console.log("Ошибка: ", data.message)
            socket.onerror(new Error(data.message))
        } else if (data.action === "admin_not_logged") {
            console.log("Ошибка: ", data.message)
            alert("Администратор еще не присоединился")
            window.location.href = "/home/role" // TODO: 
           // socket.onerror(new Error(data.message))
        }
    }

    socket.onclose = () => {
        console.log("WebSocket закрыт.");
    };

    socket.onerror = (error) => {
        console.error("Ошибка WebSocket:", error);
    };
    }


// Функция для отображения сообщений в UI
function displayMessage(message) {
    const messagesOutput = document.getElementById("messagesOutput");
    const messageElement = document.createElement("div");
    messageElement.textContent = message;
    messagesOutput.appendChild(messageElement);
    messagesOutput.scrollTop = messagesOutput.scrollHeight; // Автопрокрутка к последнему сообщению
}

// Добавляем обработчик на кнопку "Отправить"
document.addEventListener("DOMContentLoaded", () => {
    const sendButton = document.getElementById("sendButton");
    const messageInput = document.getElementById("messageInput");

    sendButton.addEventListener("click", () => {
        const message = messageInput.value.trim();
        if (message) {
            if (socket.readyState === WebSocket.OPEN) {
                socket.send(message); // Отправить сообщение через WebSocket
                displayMessage(`Вы: ${message}`); // Отобразить сообщение в UI
                messageInput.value = ""; // Очистить поле ввода
            } else {
                console.error("WebSocket не подключен");
                displayMessage("Сообщение не отправлено: WebSocket не подключен.");
            }
        }
    });

});


createWebSocket();