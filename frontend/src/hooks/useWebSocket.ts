import { useEffect, useState } from "react";
import { getToken } from "../services/api";

export type SocketEvent = {
  type: string;
  payload: unknown;
  time: string;
};

export function useWebSocket() {
  const [connected, setConnected] = useState(false);
  const [lastEvent, setLastEvent] = useState<SocketEvent | null>(null);

  useEffect(() => {
    const token = getToken();
    if (!token) {
      return;
    }
    const protocol = window.location.protocol === "https:" ? "wss" : "ws";
    const socket = new WebSocket(`${protocol}://${window.location.host}/api/ws`);
    socket.onopen = () => setConnected(true);
    socket.onclose = () => setConnected(false);
    socket.onmessage = (event) => {
      try {
        setLastEvent(JSON.parse(event.data) as SocketEvent);
      } catch {
        setLastEvent(null);
      }
    };
    return () => socket.close();
  }, []);

  return { connected, lastEvent };
}
