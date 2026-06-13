import { render, fireEvent } from '@testing-library/react';
import { ZoneTable } from './ZoneTable';
import type { SnapshotData } from '../../types';

const mockData: SnapshotData = {
  '42': { TripCount: 15, AvgFare: 23.5, LastSeen: '2026-06-13T12:00:00Z' },
  '161': { TripCount: 8, AvgFare: 67.2, LastSeen: '2026-06-13T11:55:00Z' },
  '7': { TripCount: 30, AvgFare: 12.0, LastSeen: '2026-06-13T12:01:00Z' },
};

describe('ZoneTable', () => {
  it('should render zone data', () => {
    const { getByText } = render(<ZoneTable data={mockData} />);
    expect(getByText('42')).toBeInTheDocument();
    expect(getByText('161')).toBeInTheDocument();
    expect(getByText('$23.50')).toBeInTheDocument();
    expect(getByText('$67.20')).toBeInTheDocument();
  });

  it('should show empty state when no data', () => {
    const { getByText } = render(<ZoneTable data={{}} />);
    expect(getByText('Waiting for data...')).toBeInTheDocument();
  });

  it('should sort by trips descending by default', () => {
    const { container } = render(<ZoneTable data={mockData} />);
    const rows = container.querySelectorAll('tbody tr');
    expect(rows[0]).toHaveTextContent('30');
    expect(rows[1]).toHaveTextContent('15');
    expect(rows[2]).toHaveTextContent('8');
  });

  it('should toggle sort direction on column click', () => {
    const { container, getByText } = render(<ZoneTable data={mockData} />);
    const tripsHeader = getByText(/Trips/);

    fireEvent.click(tripsHeader); // already desc, clicking toggles to asc

    const rows = container.querySelectorAll('tbody tr');
    expect(rows[0]).toHaveTextContent('8');
    expect(rows[2]).toHaveTextContent('30');
  });

  it('should sort by avg fare when header clicked', () => {
    const { container, getByText } = render(<ZoneTable data={mockData} />);
    const fareHeader = getByText(/Avg Fare/);

    fireEvent.click(fareHeader);

    const rows = container.querySelectorAll('tbody tr');
    expect(rows[0]).toHaveTextContent('$67.20');
  });
});
