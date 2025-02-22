import React, { useEffect, useState, useRef } from "react";
import "./App.css";
import { connect, sendMsg, sendPlayerData, subscribeToUpdates } from "./api";

import { Stage, Graphics } from "@pixi/react";

const WIDTH = 800;
const HEIGHT = 600;
const SPEED = 5;

const App = () => {
  const [players, setPlayers] = useState({});
  //const playerID = useRef(`player-${Math.random().toString(36).substring(7)}`);
  const storedID = localStorage.getItem("playerID");
  const playerID = useRef(
    storedID || `player-${Math.random().toString(36).substring(7)}`
  );

  useEffect(() => {
    localStorage.setItem("playerID", playerID.current);
  }, []);

  const position = useRef({ x: 100, y: 100 });

  const subscribed = useRef(false); // Prevents duplicate subscriptions

  useEffect(() => {
    connect();
    sendPlayerData({
      id: playerID.current,
      x: position.current.x,
      y: position.current.y,
    });

    if (!subscribed.current) {
      subscribeToUpdates((updatedPlayers) => {
        console.log("Updating players in subscribe function", updatedPlayers);
        setPlayers(updatedPlayers);
      });
      subscribed.current = true; // Mark as subscribed
    }
  }, []);

  const send = () => {
    console.log("hello");
    sendMsg("hello");
  };

  console.log(players);
  const movePlayer = (dx, dy) => {
    position.current.x += dx * SPEED;
    position.current.y += dy * SPEED;

    const playerData = {
      id: playerID.current,
      x: position.current.x,
      y: position.current.y,
    };
    console.log("Player data sent when moving: ", playerData);
    sendPlayerData(playerData);
  };

  useEffect(() => {
    const handleKeyDown = (e) => {
      if (e.key === "ArrowUp") movePlayer(0, -1);
      if (e.key === "ArrowDown") movePlayer(0, 1);
      if (e.key === "ArrowLeft") movePlayer(-1, 0);
      if (e.key === "ArrowRight") movePlayer(1, 0);
    };

    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, []);

  return (
    <>
      <div className="App">
        <button onClick={send}>Hit</button>
      </div>
      <Stage
        width={WIDTH}
        height={HEIGHT}
        options={{ backgroundColor: 0x333333 }}
      >
        {Object.values(players).map((p) => (
          <Graphics
            key={p.id} // Ensures React tracks each player uniquely
            draw={(g) => {
              g.clear();
              g.beginFill(p.id === playerID.current ? 0x00ff00 : 0xff0000);
              g.drawCircle(p.x, p.y, 10);
              g.endFill();
            }}
          />
        ))}
      </Stage>
    </>
  );
};

export default App;
