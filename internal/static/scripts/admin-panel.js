let socket;
const maxRetries = 3;
let retryCount = 0;

function createWebSocketAdmin() {
    socket = new WebSocket(`ws://${window.location.host}/ws-admin`);

    socket.onopen = () => {
        console.log('WebSocket подключен для админа');
        retryCount = 0; // Сбросить счетчик попыток при успешном подключении

        // Запрос списка игроков каждую секунду
        setInterval(() => {
            if (socket.readyState === WebSocket.OPEN) {
                socket.send(JSON.stringify({ Action: 'get_players' }));
                console.log('Отправлне запрос на получение списка игроков')
            }
        }, 1000);
    };

    socket.onmessage = event => {
        const data = JSON.parse(event.data);
        console.log('Сообщение получено:', data);

        if (data.type === 'players_list') {
            displayPlayers(data.content); // Обновление списка игроков
        } else {
            logSystemMessage(`Сервер: ${event.data}`); // Отобразить сообщение от сервера
        }
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

// Функция для отображения списка игроков
function displayPlayers(players) {
    const playersList = document.getElementById("playersList");
    playersList.innerHTML = ""; // Очистка предыдущего содержимого

    players.forEach(player => {
        const playerElement = document.createElement("div");
        playerElement.classList.add("player");
        playerElement.innerHTML = `
            <span class="player-name">${player.name}</span>
            <span class="player-status">${player.status}</span>
        `;
        playersList.appendChild(playerElement);
    });
}

createWebSocketAdmin();