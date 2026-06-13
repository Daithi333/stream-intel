import { useEffect, useRef, useState, useCallback } from 'react';
import type { SnapshotData } from '../types';

const WS_URL = 'ws://localhost:8080/ws';
const RECONNECT_DELAY = 3000;

type ConnectionStatus = 'connecting' | 'connected' | 'disconnected';

export function useWebSocket() {
  const [data, setData] = useState<SnapshotData>({});
  const [status, setStatus] = useState<ConnectionStatus>('connecting');
  const wsRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    let reconnectTimer: ReturnType<typeof setTimeout>;
    let closed = false;

    function connect() {
      if (closed) return;
      setStatus('connecting');
      const ws = new WebSocket(WS_URL);

      ws.onopen = () => {
        setStatus('connected');
      };

      ws.onmessage = (event) => {
        const parsed: SnapshotData = JSON.parse(event.data);
        setData(parsed);
      };

      ws.onclose = () => {
        setStatus('disconnected');
        if (!closed) {
          reconnectTimer = setTimeout(connect, RECONNECT_DELAY);
        }
      };

      ws.onerror = () => {
        ws.close();
      };

      wsRef.current = ws;
    }

    connect();

    return () => {
      closed = true;
      clearTimeout(reconnectTimer);
      wsRef.current?.close();
    };
  }, []);

  const sendReplay = useCallback(() => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify({ action: 'replay' }));
    }
  }, []);

  return { data, status, sendReplay };
}
