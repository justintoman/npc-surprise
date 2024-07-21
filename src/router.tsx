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
          console.log('login loader');
          const status = await NpcSurpriseApi.status();
          if (status.isAdmin) {
            console.log('redirect to admin');
            return redirect('/admin');
          }
          if (status.id) {
            console.log('redirect to player');
            return redirect('/');
          }
          return null;
        },
        path: '/login',
        element: <PlayerLogin />,
      },
      {
        async loader() {
          console.log('player loader');
          const status = await NpcSurpriseApi.status();
          if (status.isAdmin) {
            console.log('redirect to admin');
            return redirect('/admin');
          }
          if (status.id) {
            return null;
          }

          console.log('redirect to login');
          return redirect('/login');
        },
        path: '/',
        element: <PlayerView />,
      },
      {
        async loader() {
          console.log('admin loader');
          const status = await NpcSurpriseApi.status();
          if (status.isAdmin) {
            return null;
          }
          if (status.id) {
            console.log('redirect to player');
            return redirect('/');
          }

          console.log('redirect to login');
          return redirect('/login');
        },
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
