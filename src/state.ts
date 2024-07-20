import { atom, createStore } from 'jotai';
import { StatusResponse } from '~/api';
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
    store.set(statusAtom, null);
  };

  return function close() {
    eventSource.close();
  };
}

type CharacterMessage = {
  type: 'character';
  data: Character;
};

type ActionMessage = {
  type: 'action';
  data: Action;
};

type InitMessage = {
  type: 'init';
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
  | InitMessage
  | CharacterMessage
  | ActionMessage
  | PlayerConnectedMessage
  | PlayerDisconnectedMessage
  | DeleteMessage;

function handleEvents(message: Message) {
  switch (message.type) {
    case 'init': {
      store.set(playersAtomInternal, message.data.players);
      store.set(charactersAtomInternal, message.data.characters);
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
      const character = characters.find(
        (c) => c.id === message.data.characterId,
      );
      if (!character) {
        console.error(
          "tried to reveal action for a character that doesn't exist",
        );
        return;
      }
      const exists = character.actions.some((a) => a.id === message.data.id);
      store.set(
        charactersAtomInternal,
        characters.map((c) =>
          c === character
            ? {
                ...character,
                actions: exists
                  ? character.actions.map((action) =>
                      action.id === message.data.id ? message.data : action,
                    )
                  : [...character.actions, message.data],
              }
            : c,
        ),
      );
      break;
    }

    case 'player-connected': {
      const players = store.get(playersAtomInternal);
      const player = players.find((p) => p.id === message.data.id);
      const updated = { ...(player ?? message.data), isOnline: true };
      store.set(
        playersAtomInternal,
        player
          ? players.map((p) => (p.id === message.data.id ? updated : p))
          : [...players, updated],
      );
      break;
    }

    case 'player-disconnected': {
      const players = store.get(playersAtomInternal);
      store.set(
        playersAtomInternal,
        players.map((player) =>
          player.id === message.data ? { ...player, isOnline: false } : player,
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

    case 'delete-player': {
      const players = store.get(playersAtomInternal);
      store.set(
        playersAtomInternal,
        players.filter((player) => player.id !== message.data),
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

export const statusAtom = atom<StatusResponse | null>(null);

export const isAdminAtom = atom((get) => {
  const status = get(statusAtom);
  return Boolean(status?.isAdmin);
});
