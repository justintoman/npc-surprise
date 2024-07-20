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
  race: string;
  gender: string;
  age: string;
  description: string;
  appearance: string;
  actions: Action[];
};

export type CharacterRevealedFields = {
  [Key in keyof Omit<Character, 'id' | 'playerId' | 'actions'>]: boolean;
};

export type Action = {
  id: number;
  revealed: boolean;
  characterId: number;
  type: string;
  content: string;
  direction: string;
};
