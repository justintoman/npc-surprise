import { atom, createStore } from 'jotai';
import { atomWithRefresh, loadable } from 'jotai/utils';
import { NpcSurpriseApi } from '~/api';
import { Action, Character, Player } from '~/types';

export const store = createStore();

export function initStream() {
  const eventSource = new EventSource('/api/stream');
  eventSource.onopen = () => {
    console.log('connected');
  };

  eventSource.onmessage = (event) => {
    const message = JSON.parse(event.data);
    console.log(message);
    handleEvents(message);
  };

  eventSource.onerror = (e) => {
    console.error(e);
    store.set(statusAtom);
  };

  return function close() {
    eventSource.close();
  };
}

type AssignActionMessage = {
  type: 'assign-action';
  data: Action;
};

type AssignCharacterMessage = {
  type: 'assign-character';
  data: Character;
};

type UnassignMessage = {
  type: 'unassign-action' | 'unassign-character';
  data: number; // id
};

type CharacterMessage = {
  type: 'character';
  data: Character;
};

type ActionMessage = {
  type: 'action';
  data: Action;
};

type ConnectedMessage = {
  type: 'connected';
  data: {
    players: Player[];
    characters: Character[];
  };
};

type PlayerConnectedMessage = {
  type: 'player-connected';
  data: Player;
};

type PlayerDisconnectedMessage = {
  type: 'player-disconnected';
  data: number; // player id
};

type DeleteMessage = {
  type: 'delete-action' | 'delete-character' | 'delete-player';
  data: number;
};

type Message =
  | AssignActionMessage
  | AssignCharacterMessage
  | UnassignMessage
  | CharacterMessage
  | ActionMessage
  | ConnectedMessage
  | PlayerConnectedMessage
  | PlayerDisconnectedMessage
  | DeleteMessage;

function handleEvents(message: Message) {
  switch (message.type) {
    case 'assign-action': {
      const characters = store.get(charactersAtomInternal);
      const character = characters.find(
        (char) => char.id === message.data.character_id,
      );
      if (!character) {
        return;
      }

      store.set(
        charactersAtomInternal,
        characters.map((char) =>
          char === character
            ? { ...character, actions: [...character.actions, message.data] }
            : char,
        ),
      );
      break;
    }
    case 'assign-character': {
      const characters = store.get(charactersAtomInternal);
      store.set(charactersAtomInternal, [...characters, message.data]);
      break;
    }
    case 'unassign-action': {
      const characters = store.get(charactersAtomInternal);
      const character = characters.find((char) =>
        char.actions.some((a) => a.id === message.data),
      );
      if (!character) {
        return;
      }

      store.set(
        charactersAtomInternal,
        characters.map((char) =>
          char === character
            ? {
                ...character,
                actions: character.actions.filter((a) => a.id !== message.data),
              }
            : char,
        ),
      );
      break;
    }
    case 'unassign-character': {
      const characters = store.get(charactersAtomInternal);
      store.set(
        charactersAtomInternal,
        characters.filter((char) => char.id !== message.data),
      );
      break;
    }
    case 'character': {
      const characters = store.get(charactersAtomInternal);
      const exists = characters.some((char) => char.id === message.data.id);
      store.set(
        charactersAtomInternal,
        exists
          ? characters.map((char) =>
              char.id === message.data.id ? message.data : char,
            )
          : [...characters, message.data],
      );
      break;
    }
    case 'action': {
      const characters = store.get(charactersAtomInternal);
      store.set(
        charactersAtomInternal,
        characters.map((character) => ({
          ...character,
          actions: character.actions.map((a) =>
            a.id === message.data.id ? message.data : a,
          ),
        })),
      );
      break;
    }
    case 'connected': {
      store.set(playersAtomInternal, message.data.players);
      store.set(charactersAtomInternal, message.data.characters);
      break;
    }
    case 'player-connected': {
      const players = store.get(playersAtomInternal);
      const exists = players.some((p) => p.id === message.data.id);
      store.set(
        playersAtomInternal,
        exists
          ? players.map((p) => (p.id === message.data.id ? message.data : p))
          : [...players, { ...message.data, is_online: true }],
      );
      break;
    }
    case 'player-disconnected': {
      const players = store.get(playersAtomInternal);
      store.set(
        playersAtomInternal,
        players.map((player) =>
          player.id === message.data ? { ...player, is_online: false } : player,
        ),
      );
      break;
    }
    case 'delete-action': {
      const characters = store.get(charactersAtomInternal);
      store.set(
        charactersAtomInternal,
        characters.map((char) => ({
          ...char,
          actions: char.actions.filter((a) => a.id !== message.data),
        })),
      );
      break;
    }
    case 'delete-character': {
      const characters = store.get(charactersAtomInternal);
      store.set(
        charactersAtomInternal,
        characters.filter((char) => char.id !== message.data),
      );
      break;
    }
    default: {
      console.error('Invalid message type', message);
    }
  }
}

const playersAtomInternal = atom<Player[]>([]);
export const playersAtom = atom<Player[]>((get) => get(playersAtomInternal));

const charactersAtomInternal = atom<Character[]>([]);
export const charactersAtom = atom<Character[]>((get) =>
  get(charactersAtomInternal),
);

export const statusAtom = atomWithRefresh(async () => {
  return NpcSurpriseApi.status();
});

export const isAdminAtom = atom((get) => {
  const loaded = get(loadable(statusAtom));
  return loaded.state === 'hasData' && loaded.data.is_admin;
});
