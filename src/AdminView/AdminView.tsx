import { useEffect } from 'react';
import { initStream } from '~/state';

export function AdminView() {
  useEffect(() => {
    return initStream();
  }, []);

  return (
    <div>
      <header>
        <h1>Admin View</h1>
      </header>
    </div>
  );
}
