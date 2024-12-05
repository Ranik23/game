let socket;
const maxRetries = 1;
let retryCount = 0;

function createWebSocket() {
  socket = new WebSocket(`ws://${window.location.host}/ws`);

  socket.onopen = () => {
    console.log('WebSocket подключен');
    retryCount = 0; // Сбросить счетчик попыток при успешном подключении
  };

  socket.onmessage = event => {
    console.log('Сообщение получено:', event.data);
    displayMessage(`Сервер: ${event.data}`); // Отобразить сообщение от сервера
  };

  socket.onclose = () => {
    console.log('WebSocket закрыт. Попытка переподключения...');
    if (retryCount < maxRetries) {
      retryCount++;
      setTimeout(createWebSocket, 1000 * retryCount); // Увеличение времени ожидания перед повторной попыткой
    } else {
      console.error('Достигнуто максимальное количество попыток подключения.');
    }
  };

  socket.onerror = error => {
    console.error('Ошибка WebSocket:', error);
  };
}

// Функция для отображения сообщений в UI
function displayMessage(message) {
  const messagesOutput = document.getElementById('messagesOutput');
  const messageElement = document.createElement('div');
  messageElement.textContent = message;
  messagesOutput.appendChild(messageElement);
}

// Добавляем обработчик на кнопку "Отправить"
document.addEventListener('DOMContentLoaded', () => {
  const sendButton = document.getElementById('sendButton');
  const messageInput = document.getElementById('messageInput');

  sendButton.addEventListener('click', () => {
    const message = messageInput.value.trim();
    if (message && socket.readyState === WebSocket.OPEN) {
      socket.send(message); // Отправить сообщение через WebSocket
      displayMessage(`Вы: ${message}`); // Отобразить сообщение в UI
      messageInput.value = ''; // Очистить поле ввода
    } else if (socket.readyState !== WebSocket.OPEN) {
      console.error('WebSocket не подключен');
    }
  });
});

// Вызов функции для создания WebSocket
createWebSocket();
