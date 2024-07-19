import ky from 'ky';
import { Action, Character, Player } from '~/types';

const client = ky.create({
  prefixUrl: '/api',
});

export const NpcSurpriseApi = {
  login(name: string): Promise<Required<StatusResponse>> {
    const response = client.post('login', { json: { name: name } });
    return response.json();
  },
  getPlayers(): Promise<Player[]> {
    return client.get('players').json<Player[]>();
  },
  status(): Promise<StatusResponse> {
    return client.get('status').json();
  },

  getCharacters(): Promise<Array<Character>> {
    return client.get('characters').json<Array<Character>>();
  },

  createCharacter(character: Omit<Character, 'id' | 'actions'>) {
    return client.post('characters', { json: character }).json<Character>();
  },

  updateCharacter(character: Omit<Character, 'actions'>) {
    return client
      .put(`characters/${character.id}`, { json: character })
      .json<Character>();
  },

  deleteCharacter(id: number) {
    return client.delete(`characters/${id}`).json();
  },

  createAction(action: Omit<Action, 'id'>) {
    return client.post('actions', { json: action }).json<Action>();
  },

  updateAction(action: Action) {
    return client.put(`actions/${action.id}`, { json: action }).json<Action>();
  },

  deleteAction(id: number) {
    return client.delete(`actions/${id}`).json();
  },

  assign(type: 'action' | 'character', id: number, playerId: number) {
    return client.post('assign', { json: { id, playerId, type } }).json();
  },

  deletePlayer(id: number): Promise<void> {
    return client.delete(`players/${id}`).json();
  },
};

export type StatusResponse = {
  isAdmin: boolean;
  playerId?: string;
  playerName?: string;
};
