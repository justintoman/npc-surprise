import { atom, createStore } from 'jotai';
import { atomFamily } from 'jotai/utils';
import { Action, Character, CharacterRevealedFields, Player } from '~/types';

export const store = createStore();

export function initStream() {
  const url = `${import.meta.env.VITE_API_PREFIX ?? ''}/stream`;
  const eventSource = new EventSource(url);
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

type CharacterMessage = {
  type: 'character';
  data: Character;
};

type CharacterWithFieldsMessage = {
  type: 'character-with-fields';
  data: {
    character: Character;
    fields: CharacterRevealedFields;
  };
};

type ActionMessage = {
  type: 'action';
  data: Action;
};

type InitPlayerMessage = {
  type: 'init-player';
  data: Character[];
};

type InitAdminMessage = {
  type: 'init-admin';
  data: {
    players: Player[];
    characters: Character[];
    fields: CharacterRevealedFields[];
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
  | InitPlayerMessage
  | InitAdminMessage
  | CharacterMessage
  | CharacterWithFieldsMessage
  | ActionMessage
  | PlayerConnectedMessage
  | PlayerDisconnectedMessage
  | DeleteMessage;

function handleEvents(message: Message) {
  switch (message.type) {
    case 'init-admin': {
      store.set(playersAtomInternal, message.data.players);
      store.set(charactersAtomInternal, message.data.characters);
      store.set(characterRevealedFieldsInternal, message.data.fields);
      break;
    }

    case 'init-player': {
      store.set(charactersAtomInternal, message.data);
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

    case 'character-with-fields': {
      const characters = store.get(charactersAtomInternal);
      const exists = characters.some(
        (char) => char.id === message.data.character.id,
      );
      store.set(
        charactersAtomInternal,
        exists
          ? characters.map((char) =>
              char.id === message.data.character.id
                ? message.data.character
                : char,
            )
          : [...characters, message.data.character],
      );
      const fields = store.get(characterRevealedFieldsInternal);
      const fieldExists = fields.some(
        (f) => f.characterId === message.data.character.id,
      );
      store.set(
        characterRevealedFieldsInternal,
        fieldExists
          ? fields.map((f) =>
              f.characterId === message.data.character.id
                ? message.data.fields
                : f,
            )
          : [...fields, message.data.fields],
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
export const playerAtomFamily = atomFamily((id: number | undefined) =>
  atom((get) => {
    if (id === undefined) {
      return undefined;
    }
    return get(playersAtomInternal).find((player) => player.id === id);
  }),
);

const charactersAtomInternal = atom<Character[]>([]);
export const charactersAtom = atom<Character[]>((get) =>
  get(charactersAtomInternal),
);

export const characterAtomFamily = atomFamily((id: number) =>
  atom((get) => get(charactersAtomInternal).find((char) => char.id === id)),
);

export const actionAtomFamily = atomFamily((id: number) =>
  atom((get) =>
    get(charactersAtom)
      .flatMap((char) => char.actions)
      .find((action) => action.id === id),
  ),
);

const characterRevealedFieldsInternal = atom<CharacterRevealedFields[]>([]);

export const characterRevealedFieldsAtomFamily = atomFamily((id: number) =>
  atom<CharacterRevealedFields | undefined>((get) => {
    const fields = get(characterRevealedFieldsInternal).find(
      (f) => f.characterId === id,
    );
    return fields;
  }),
);
