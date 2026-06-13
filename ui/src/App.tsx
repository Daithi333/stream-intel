import { useWebSocket } from './hooks/use-websocket';
import { StatusBar } from './components/StatusBar/StatusBar';
import { ZoneTable } from './components/ZoneTable/ZoneTable';

export function App() {
  const { data, status, sendReplay } = useWebSocket();
  const zoneCount = Object.keys(data).length;

  return (
    <div>
      <StatusBar status={status} zoneCount={zoneCount} onReplay={sendReplay} />
      <ZoneTable data={data} />
    </div>
  );
}
