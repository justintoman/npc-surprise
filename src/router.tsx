import { createBrowserRouter, redirect } from 'react-router-dom';
import { AdminView } from '~/AdminView/AdminView';
import { NpcSurpriseApi } from '~/api';
import { App } from '~/App';
import { PlayerLogin } from '~/PlayerLogin/PlayerLogin';
import { PlayerView } from '~/PlayerView/PlayerView';
import { statusAtom, store } from '~/state';

export const router = createBrowserRouter([
  {
    async loader() {
      const status = await NpcSurpriseApi.status();
      store.set(statusAtom, status);
      return null;
    },
    element: <App />,
    children: [
      {
        async loader() {
          const status = await NpcSurpriseApi.status();
          store.set(statusAtom, status);
          if (status.isAdmin) {
            return redirect('/admin');
          }
          if (!status.playerId) {
            return redirect('/login');
          }

          return redirect('/player');
        },
        path: '/',
      },
      {
        path: '/login',
        element: <PlayerLogin />,
      },
      {
        path: '/player',
        element: <PlayerView />,
      },
      {
        path: '/admin',
        element: <AdminView />,
      },
    ],
  },
]);
