import { useState } from 'react';
import type { SnapshotData, ZoneSnapshot } from '../../types';
import styles from './ZoneTable.module.css';

interface ZoneTableProps {
  data: SnapshotData;
}

type SortKey = 'zone' | 'TripCount' | 'AvgFare' | 'LastSeen';
type SortDir = 'asc' | 'desc';

function formatTime(isoString: string): string {
  if (!isoString) return '-';
  const date = new Date(isoString);
  return date.toLocaleTimeString();
}

function formatFare(fare: number): string {
  return `$${fare.toFixed(2)}`;
}

function compare(a: [string, ZoneSnapshot], b: [string, ZoneSnapshot], key: SortKey, dir: SortDir): number {
  let result: number;
  switch (key) {
    case 'zone':
      result = Number(a[0]) - Number(b[0]);
      break;
    case 'TripCount':
      result = a[1].TripCount - b[1].TripCount;
      break;
    case 'AvgFare':
      result = a[1].AvgFare - b[1].AvgFare;
      break;
    case 'LastSeen':
      result = new Date(a[1].LastSeen).getTime() - new Date(b[1].LastSeen).getTime();
      break;
  }
  return dir === 'asc' ? result : -result;
}

export function ZoneTable({ data }: ZoneTableProps) {
  const [sortKey, setSortKey] = useState<SortKey>('TripCount');
  const [sortDir, setSortDir] = useState<SortDir>('desc');

  const handleSort = (key: SortKey) => {
    if (key === sortKey) {
      setSortDir(sortDir === 'asc' ? 'desc' : 'asc');
    } else {
      setSortKey(key);
      setSortDir('desc');
    }
  };

  const zones = Object.entries(data).sort((a, b) => compare(a, b, sortKey, sortDir));

  if (zones.length === 0) {
    return (
      <div className={styles.empty}>
        Waiting for data...
      </div>
    );
  }

  const indicator = (key: SortKey) => {
    if (key !== sortKey) return '';
    return sortDir === 'asc' ? ' \u25B2' : ' \u25BC';
  };

  return (
    <div className={styles.wrapper}>
      <table className={styles.table}>
        <thead>
          <tr>
            <th className={styles.sortable} onClick={() => handleSort('zone')}>
              Zone{indicator('zone')}
            </th>
            <th className={styles.sortable} onClick={() => handleSort('TripCount')}>
              Trips{indicator('TripCount')}
            </th>
            <th className={styles.sortable} onClick={() => handleSort('AvgFare')}>
              Avg Fare{indicator('AvgFare')}
            </th>
            <th className={styles.sortable} onClick={() => handleSort('LastSeen')}>
              Last Seen{indicator('LastSeen')}
            </th>
          </tr>
        </thead>
        <tbody>
          {zones.map(([zone, stats]) => (
            <tr key={zone}>
              <td className={styles.zone}>{zone}</td>
              <td className={styles.mono}>{stats.TripCount}</td>
              <td className={styles.mono}>{formatFare(stats.AvgFare)}</td>
              <td className={styles.muted}>{formatTime(stats.LastSeen)}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
