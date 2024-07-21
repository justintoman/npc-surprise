import { createBrowserRouter, redirect } from 'react-router-dom';
import { AdminHome } from '~/Admin/AdminHome';
import { AdminPage } from '~/Admin/AdminPage';
import { EditAction } from '~/Admin/EditAction';
import { EditCharacter } from '~/Admin/EditCharacter';
import { NewAction } from '~/Admin/NewAction';
import { NewCharacter } from '~/Admin/NewCharacter';
import { NpcSurpriseApi } from '~/api';
import { App } from '~/App';
import { PlayerLogin } from '~/PlayerLogin/PlayerLogin';
import { PlayerView } from '~/PlayerView/PlayerView';

export const router = createBrowserRouter([
  {
    element: <App />,
    children: [
      {
        async loader() {
          const status = await NpcSurpriseApi.status();
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
        element: <AdminPage />,
        children: [
          {
            path: '',
            element: <AdminHome />,
          },
          {
            path: 'character/new',
            element: <NewCharacter />,
          },
          {
            path: 'character/:characterId',
            element: <EditCharacter />,
          },
          {
            path: ':characterId/action/new',
            element: <NewAction />,
          },
          {
            path: ':characterId/action/:actionId',
            element: <EditAction />,
          },
        ],
      },
    ],
  },
]);
