let socket = new WebSocket("ws://localhost:8080/ws");

export const connect = () => {
  if (socket && socket.readyState === WebSocket.OPEN) {
    console.log("WebSocket already connected");
    return;
  }

  socket = new WebSocket("ws://localhost:8080/ws");

  socket.onopen = () => {
    console.log("Successfully connected to WebSocket");
  };

  socket.onclose = (event) => {
    console.log("Socket closed connection:", event);
  };

  socket.onerror = (err) => {
    console.log("Socket error:", err);
  };
};

let sendMsg = (msg) => {
  console.log("Sending message: ", msg);
  socket.send(msg);
};

export const sendPlayerData = (playerData) => {
  if (socket.readyState === WebSocket.OPEN) {
    socket.send(JSON.stringify(playerData)); // Ensures JSON format
  } else {
    console.warn("WebSocket not open, cannot send data.");
  }
};

export const subscribeToUpdates = (callback) => {
  if (!socket) {
    console.error("WebSocket is not initialized.");
    return;
  }

  if (socket.onmessage) {
    console.warn(
      "WebSocket onmessage handler already set, avoiding duplicate listeners."
    );
    return; // Prevents overwriting the message handler
  }

  socket.onmessage = (event) => {
    try {
      if (typeof event.data === "string") {
        const updatedPlayers = JSON.parse(event.data);
        console.log("Updating players in subscribe function", updatedPlayers);
        callback((prevPlayers) => {
          // Remove players who are no longer in updatedPlayers
          const filteredPlayers = Object.fromEntries(
            Object.entries(prevPlayers).filter(([id]) => updatedPlayers[id])
          );

          return { ...filteredPlayers, ...updatedPlayers };
        });
      } else {
        console.error("Received non-string data from WebSocket:", event.data);
      }
    } catch (error) {
      console.error("Error parsing WebSocket data:", error, event.data);
    }
  };
};

export { sendMsg };
