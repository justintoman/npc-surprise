import { type ClassValue, clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';
import { Action, Character } from '~/types';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function getNewAction(): Omit<Action, 'id' | 'characterId'> {
  return {
    type: '',
    direction: '',
    content: '',
  };
}

export function getNewCharacter(): Omit<Character, 'id' | 'actions'> {
  return {
    name: '',
    race: '',
    gender: '',
    age: '',
    description: '',
    appearance: '',
  };
}
