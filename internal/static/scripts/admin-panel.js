let socket;
const maxRetries = 1;
let retryCount = 0;


function createWebSocketAdmin() {
    socket = new WebSocket(`ws://${window.location.host}/wsmain`)
    socket.onopen = () => {
        console.log('WebSocket подключен для админа');
        retryCount = 0; // Сбросить счетчик попыток при успешном подключении
      };
    
      socket.onmessage = event => {
        console.log('Сообщение получено:', event.data);
        logSystemMessage(`Сервер: ${event.data}`); // Отобразить сообщение от сервера
      };
    
      socket.onclose = () => {
        console.log('WebSocket закрыт. Попытка переподключения...');
        if (retryCount < maxRetries) {
          retryCount++;
          setTimeout(createWebSocketAdmin, 1000 * retryCount); // Увеличение времени ожидания перед повторной попыткой
        } else {
          console.error('Достигнуто максимальное количество попыток подключения.');
        }
      };
    
      socket.onerror = error => {
        console.error('Ошибка WebSocket:', error);
      };
}

function logSystemMessage(message) {
  const systemLog = document.getElementById('system-log');
  const logElement = document.createElement('div');
  logElement.textContent = message;
  systemLog.appendChild(logElement);
}

document.addEventListener('DOMContentLoaded', () => {
  const kickButton = document.getElementById('kickUserButton');
  kickButton.addEventListener('click', () => {
      if (socket.readyState === WebSocket.OPEN) {
          const userId = prompt('Введите ID пользователя для удаления:');
          if (userId) {
              socket.send(JSON.stringify({ type: 'kick', userId }));
              logSystemMessage(`Пользователь с ID ${userId} удален.`);
          }
      }
  });
});


createWebSocketAdmin()