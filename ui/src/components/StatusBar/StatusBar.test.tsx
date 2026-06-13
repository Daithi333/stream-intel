import { render, fireEvent } from '@testing-library/react';
import { StatusBar } from './StatusBar';

describe('StatusBar', () => {
  it('should display zone count', () => {
    const { getByText } = render(
      <StatusBar status="connected" zoneCount={42} onReplay={() => {}} />
    );
    expect(getByText('42 zones')).toBeInTheDocument();
  });

  it('should show connected status', () => {
    const { getByText } = render(
      <StatusBar status="connected" zoneCount={0} onReplay={() => {}} />
    );
    expect(getByText('connected')).toBeInTheDocument();
  });

  it('should show disconnected status', () => {
    const { getByText } = render(
      <StatusBar status="disconnected" zoneCount={0} onReplay={() => {}} />
    );
    expect(getByText('disconnected')).toBeInTheDocument();
  });

  it('should call onReplay when button clicked', () => {
    const onReplay = vi.fn();
    const { getByText } = render(
      <StatusBar status="connected" zoneCount={0} onReplay={onReplay} />
    );
    fireEvent.click(getByText('Replay'));
    expect(onReplay).toHaveBeenCalledOnce();
  });

  it('should disable replay button when disconnected', () => {
    const { getByText } = render(
      <StatusBar status="disconnected" zoneCount={0} onReplay={() => {}} />
    );
    expect(getByText('Replay')).toBeDisabled();
  });
});
