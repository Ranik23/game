let socket;


function createWebSocketAdmin() {
    
    socket = new WebSocket(`ws://${window.location.host}/ws/leader`);

    socket.onopen = () => {
        console.log("WebSocket подключен");
    };

    socket.onmessage = (event) => {
        try {
            
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