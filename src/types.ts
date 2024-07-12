export type Player = {
  id: number;
  name: string;
  is_online: boolean;
};

export type CurrentPlayer = Omit<Player, 'is_online'>;

export type Character = {
  id: number;
  name: string;
  race: string;
  gender: string;
  age: string;
  description: string;
  appearance: string;
};

export type Action = {
  id: number;
  character_id: number;
  type: string;
  content: string;
  direction: string;
};
