import { useEffect } from 'react';
import { Outlet } from 'react-router-dom';
import { initStream } from '~/state';

export function AdminPage() {
  useEffect(() => {
    return initStream();
  }, []);
  return <Outlet />;
}
