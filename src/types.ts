export type Player = {
  id: number;
  name: string;
  isOnline: boolean;
};

export type CurrentPlayer = Omit<Player, 'isOnline'>;

export type Character = {
  id: number;
  playerId?: number;
  name: string;
  age: string;
  race: string;
  gender: string;
  description: string;
  appearance: string;
  actions: Action[];
};

export type CharacterRevealedFields = {
  [Key in keyof Omit<Character, 'id' | 'playerId' | 'actions'>]: boolean;
} & { characterId: number };

export type Action = {
  id: number;
  revealed: boolean;
  characterId: number;
  content: string;
};
