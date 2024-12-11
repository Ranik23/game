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
        console.log("Сообщение получено:", event.data);
        displayMessage(`Сервер: ${event.data}`); // Отобразить сообщение от сервера
    };

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