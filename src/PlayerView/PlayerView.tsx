import { useEffect } from 'react';
import { initStream } from '~/state';

export function PlayerView() {
  useEffect(() => {
    return initStream();
  });

  return (
    <div>
      <header>
        <h1>Player View</h1>
      </header>
    </div>
  );
}
