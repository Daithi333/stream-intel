import styles from './StatusBar.module.css';

interface StatusBarProps {
  status: 'connecting' | 'connected' | 'disconnected';
  zoneCount: number;
  onReplay: () => void;
}

export function StatusBar({ status, zoneCount, onReplay }: StatusBarProps) {
  return (
    <header className={styles.bar}>
      <div className={styles.left}>
        <h1 className={styles.title}>stream-intel</h1>
        <span className={styles.zoneCount}>{zoneCount} zones</span>
      </div>
      <div className={styles.right}>
        <button className={styles.replayButton} onClick={onReplay} disabled={status !== 'connected'}>
          Replay
        </button>
        <span className={`${styles.status} ${styles[status]}`}>
          {status}
        </span>
      </div>
    </header>
  );
}
