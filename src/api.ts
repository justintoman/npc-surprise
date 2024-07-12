import ky from 'ky';
import { CurrentPlayer, Player } from '~/types';

const client = ky.create({
  prefixUrl: 'api',
  hooks: {
    beforeRequest: [
      (request) => {
        const adminKey = localStorage.getItem('adminKey');
        if (adminKey) {
          request.headers.set('X-Npc-Surprise-Admin-Key', adminKey);
        }
      },
    ],
  },
});

export const NpcSurpriseApi = {
  login(name: string): Promise<CurrentPlayer> {
    const response = client.post('login', { json: { name: name } });
    return response.json<Player>();
  },
  getPlayers(): Promise<Player[]> {
    return client.get('players').json<Player[]>();
  },
  status(): Promise<StatusResponse> {
    return client.get('status').json();
  },
};

type StatusResponse = {
  is_admin: boolean;
  player_id?: string;
  player_name?: string;
};
