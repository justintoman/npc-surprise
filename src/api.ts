import ky from 'ky';
import { Action, Character, CharacterRevealedFields } from '~/types';

const client = ky.create({
  prefixUrl: '/api',
});

export const NpcSurpriseApi = {
  // auth
  login(name: string): Promise<Required<StatusResponse>> {
    const response = client.post('login', { json: { name: name } });
    return response.json();
  },
  status(): Promise<StatusResponse> {
    return client.get('status').json();
  },

  /**
   * Characters
   */

  createCharacter(character: Omit<Character, 'id' | 'actions'>) {
    return client.post('characters', { json: character }).json<Character>();
  },

  updateCharacter(character: Omit<Character, 'actions'>) {
    return client
      .put(`characters/${character.id}`, { json: character })
      .json<Character>();
  },

  assignCharacter(id: number, playerId: number) {
    return client.put(`characters/${id}/assign/${playerId}`).json();
  },

  unassignCharacter(id: number) {
    return client.put(`characters/${id}/unassign`).json();
  },

  updateRevealedFields(charaterId: number, fields: CharacterRevealedFields) {
    return client
      .put(`characters/${charaterId}/reveal`, { json: fields })
      .json();
  },

  deleteCharacter(id: number) {
    return client.delete(`characters/${id}`).json();
  },

  /**
   * Actions
   */

  createAction(action: Omit<Action, 'id'>) {
    return client.post(`characters/${action.characterId}/actions`, {
      json: action,
    });
  },

  updateAction(action: Action) {
    return client.put(`characters/${action.characterId}/actions/${action.id}`, {
      json: action,
    });
  },

  revealAction(characterId: number, actionId: number) {
    return client.put(`characters/${characterId}/actions/${actionId}/reveal`);
  },

  hideAction(characterId: number, actionId: number) {
    return client.put(`characters/${characterId}/actions/${actionId}/hide`);
  },

  deleteAction(characterId: number, actionId: number) {
    return client
      .delete(`characters/${characterId}/actions/${actionId}`)
      .json();
  },

  deletePlayer(playerId: number): Promise<void> {
    return client.delete(`players/${playerId}`).json();
  },
};

export type StatusResponse = {
  isAdmin: boolean;
  playerId?: string;
  playerName?: string;
};
