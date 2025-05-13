import { onCleanup, onMount } from "solid-js";
import { getSessionID } from "../utils/session";

let ws: WebSocket;

const Spawn = () => {
    onMount(() => {
        const sessionID = getSessionID();
        ws = new WebSocket(`ws://localhost:8080/ws?session=${sessionID}`);

        ws.onopen = () => console.log("WebSocket connected");
        ws.onmessage = (e) => console.log("Message:", e.data);
        ws.onerror = (e) => console.error("WebSocket error", e);
        ws.onclose = () => console.log("WebSocket closed");
    });

    onCleanup(() => {
        ws?.close();
    });

    return (
        <div class="game-screen">
            <h2>Game Screen</h2>
            <p>WebSocket connected, waiting for game data...</p>
            {/* Game canvas or rectangle rendering logic goes here */}
        </div>
    );
};

export default Spawn;
