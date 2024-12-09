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
                console.log('Отправление запрос на получение списка игроков')
            }
        }, 1000);
    };

    socket.onmessage = event => {
        const data = JSON.parse(event.data);
        console.log('Сообщение получено:', data);

        if (data.action === 'players_list') {
            displayPlayers(data.content); // Обновление списка игроков
        }
        if (data.action === 'player_accepted') {
            console.log('Игрок принят')
        }
        if (data.action === 'player_rejected') {
            console.log('Игрок отклонен')
        }
    };

    socket.onclose = () => {
        console.log('WebSocket закрыт')
    };

    socket.onerror = error => {
        console.error('Ошибка WebSocket:', error);
    };
}

document.addEventListener("DOMContentLoaded", () => {
    const playersList = document.getElementById("playersList");
  
    playersList.addEventListener("click", (event) => {
      const target = event.target;
  
      if (target.classList.contains("accept-btn")) {
        const userId = target.closest("li").dataset.userId;
        alert(`Игрок с ID ${userId} принят`);
        socket.send(JSON.stringify({ Action: 'accept_player', Data: userId }));
        // Здесь вы можете отправить запрос на сервер для обработки принятия
      }
  
      // Проверяем, нажата ли кнопка "Удалить"
      if (target.classList.contains("remove-btn")) {
        const userId = target.closest("li").dataset.userId;
        alert(`Игрок с ID ${userId} удалён`);
        socket.send(JSON.stringify({ Action: 'delete_player', Data: userId }));
        target.closest("li").remove(); // Удаляет элемент из списка
      }
    });
  });
  

// Функция для отображения списка игроков
function displayPlayers(players) {
    const playersList = document.getElementById("playersList");
    playersList.innerHTML = ""; // Очистка предыдущего содержимого

    // Проверка, что players существует и является массивом
    if (Array.isArray(players) && players.length > 0) {
        players.forEach(player => {
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
        // Если игроков нет, можно не выводить ничего или показать сообщение
        console.log('Нет игроков для отображения');
    }
}



createWebSocketAdmin();