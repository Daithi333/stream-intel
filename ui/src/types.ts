export interface ZoneSnapshot {
  TripCount: number;
  AvgFare: number;
  LastSeen: string;
}

export type SnapshotData = Record<string, ZoneSnapshot>;
