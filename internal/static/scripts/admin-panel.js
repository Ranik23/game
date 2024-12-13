let socket;
let playersUpdateInterval;

// Функция для создания WebSocket соединения для админа
function createWebSocketAdmin() {
    
    socket = new WebSocket(`ws://${window.location.host}/ws/admin`);

    socket.onopen = () => {
        console.log("WebSocket подключен для админа");

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
                displayPlayers(data.content);
                console.log("Список игроков принят")
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
    };

    socket.onerror = (error) => {
        console.error("Ошибка WebSocket:", error);
    };
}

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

function displayPlayers(players) {
    const playersList = document.getElementById("playersList");
    playersList.innerHTML = ""; 

    if (Array.isArray(players) && players.length > 0) {
        players.forEach((player) => {
            const playerElement = document.createElement("li");
            playerElement.dataset.userId = player.id;
            playerElement.innerHTML = `
                <span class="player-name">${player.userName}</span>
                <button class="accept-btn">Принять</button>
                <button class="remove-btn">Удалить</button>
            `;
            playersList.appendChild(playerElement);
        });
    } else {
        console.log("Нет игроков для отображения");
    }
}

document.addEventListener("DOMContentLoaded", () => {
    const playersList = document.getElementById("playersList");

    playersList.addEventListener("click", (event) => {
        const target = event.target;

        if (target.classList.contains("accept-btn")) {
            const userId = target.closest("li").dataset.userId;
            alert(`Игрок с ID ${userId} принят`);
            sendWebSocketMessage({ Action: "accept_player", Data: userId });
            target.closest("li").remove();
        }

        if (target.classList.contains("remove-btn")) {
            const userId = target.closest("li").dataset.userId;
            alert(`Игрок с ID ${userId} удалён`);
            sendWebSocketMessage({ Action: "delete_player", Data: userId });
            target.closest("li").remove();
        }
    });
});

// Запускаем WebSocket
createWebSocketAdmin();
