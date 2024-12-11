let socket;
let retryCount = 0;
const maxRetries = 3;
let reconnectTimeout;
let playersUpdateInterval;

function createWebSocketAdmin() {
    socket = new WebSocket(`ws://${window.location.host}/ws/admin`);

    socket.onopen = () => {
        console.log("WebSocket подключен для админа");
        retryCount = 0; // Сбросить счетчик попыток
        clearTimeout(reconnectTimeout); // Очистить таймер переподключения

        // Установить таймер для регулярного обновления списка игроков
        if (!playersUpdateInterval) {
            playersUpdateInterval = setInterval(() => {
                if (socket.readyState === WebSocket.OPEN) {
                    socket.send(JSON.stringify({ Action: "get_players" }));
                    console.log("Отправлен запрос на получение списка игроков");
                }
            }, 1000);
        }
    };

    socket.onmessage = (event) => {
        try {
            const data = JSON.parse(event.data);
            console.log("Сообщение получено:", data);

            if (data.action === "players_list") {
                displayPlayers(data.content); // Обновление списка игроков
            } else if (data.action === "player_accepted") {
                console.log("Игрок принят");
            } else if (data.action === "player_rejected") {
                console.log("Игрок отклонён");
            }
        } catch (error) {
            console.error("Ошибка обработки сообщения:", error);
        }
    };

    socket.onclose = () => {
        console.log("WebSocket закрыт");
        clearInterval(playersUpdateInterval);
        playersUpdateInterval = null;

        // Попробовать переподключиться
        if (retryCount < maxRetries) {
            retryCount++;
            console.log(`Попытка переподключения (${retryCount}/${maxRetries})...`);
            reconnectTimeout = setTimeout(createWebSocketAdmin, 2000 * retryCount); // Увеличиваем задержку
        } else {
            console.error("Превышено количество попыток переподключения.");
        }
    };

    socket.onerror = (error) => {
        console.error("Ошибка WebSocket:", error);
    };
}

document.addEventListener("DOMContentLoaded", () => {
    const playersList = document.getElementById("playersList");

    playersList.addEventListener("click", (event) => {
        const target = event.target;

        if (target.classList.contains("accept-btn")) {
            const userId = target.closest("li").dataset.userId;
            alert(`Игрок с ID ${userId} принят`);
            sendWebSocketMessage({ Action: "accept_player", Data: userId });
        }

        if (target.classList.contains("remove-btn")) {
            const userId = target.closest("li").dataset.userId;
            alert(`Игрок с ID ${userId} удалён`);
            sendWebSocketMessage({ Action: "delete_player", Data: userId });
            target.closest("li").remove(); // Удаляет элемент из списка
        }
    });
});

// Универсальная функция отправки сообщений
function sendWebSocketMessage(message) {
    try {
        if (socket && socket.readyState === WebSocket.OPEN) {
            socket.send(JSON.stringify(message));
        } else {
            console.error("WebSocket не подключен. Сообщение не отправлено:", message);
        }
    } catch (error) {
        console.error("Ошибка отправки сообщения через WebSocket:", error);
    }
}

// Функция для отображения списка игроков
function displayPlayers(players) {
    const playersList = document.getElementById("playersList");
    playersList.innerHTML = ""; // Очистка предыдущего содержимого

    if (Array.isArray(players) && players.length > 0) {
        players.forEach((player) => {
            const playerElement = document.createElement("li");
            playerElement.dataset.userId = player.id; // Уникальный идентификатор игрока
            playerElement.innerHTML = `
                <span class="player-name">${player.name}</span>
                <button class="accept-btn">Принять</button>
                <button class="remove-btn">Удалить</button>
            `;
            playersList.appendChild(playerElement);
        });
    } else {
        console.log("Нет игроков для отображения");
    }
}

// Запускаем WebSocket
createWebSocketAdmin();
