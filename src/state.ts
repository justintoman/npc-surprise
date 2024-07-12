import { atom, createStore } from 'jotai';
import { atomFamily, atomWithRefresh } from 'jotai/utils';
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

type UnassignActionMessage = {
  type: 'unassign-action';
  data: number; // id
};

type UnassignCharacterMessage = {
  type: 'unassign-character';
  data: number; // id
};

type CharacterMessage = {
  type: 'character';
  data: Character;
};

type ConnectedMessage = {
  type: 'connected';
  data: null;
};

type Message =
  | AssignActionMessage
  | AssignCharacterMessage
  | UnassignActionMessage
  | UnassignCharacterMessage
  | CharacterMessage
  | ConnectedMessage;

function handleEvents(message: Message) {
  switch (message.type) {
    case 'assign-action': {
      const actions = store.get(actionsAtom);
      store.set(actionsAtom, [...actions, message.data]);
      break;
    }
    case 'assign-character': {
      const characters = store.get(assignedCharactersAtom);
      store.set(assignedCharactersAtom, [...characters, message.data]);
      break;
    }
    case 'unassign-action': {
      const actions = store.get(actionsAtom);
      store.set(
        actionsAtom,
        actions.filter((action) => action.id !== message.data),
      );
      break;
    }
    case 'unassign-character': {
      const characters = store.get(assignedCharactersAtom);
      store.set(
        assignedCharactersAtom,
        characters.filter((char) => char.id !== message.data),
      );
      break;
    }
    case 'character': {
      const characters = store.get(assignedCharactersAtom);
      store.set(
        assignedCharactersAtom,
        characters.map((char) =>
          char.id === message.data.id ? message.data : char,
        ),
      );
      break;
    }
    case 'connected': {
      console.log('connected');
      break;
    }
    default: {
      console.error('Invalid message type', message);
    }
  }
}

export const playersAtom = atom<Player[]>([]);

export const allCharactersAtom = atom<Array<Character & { actions: Action[] }>>(
  [],
);

export const assignedCharactersAtom = atom<Character[]>([]);

export const playerAtom = atom<Player | null>(null);

export const actionsAtom = atom<Action[]>([]);

export const characterActions = atomFamily((characterId: number) =>
  atom((get) => {
    return get(actionsAtom).filter(
      (action) => action.character_id === characterId,
    );
  }),
);

export const statusAtom = atomWithRefresh(async () => {
  return NpcSurpriseApi.status();
});

export const currentPlayerAtom = atom<Omit<Player, 'is_online'> | null>(null);
