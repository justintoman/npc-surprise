import { useAtomValue } from 'jotai';
import { Plus, X } from 'lucide-react';
import { useEffect, useState } from 'react';
import { ActionForm } from '~/AdminView/ActionForm';
import { CharacterForm } from '~/AdminView/CharacterForm';
import { NpcSurpriseApi } from '~/api';
import { Character } from '~/components/Character';
import { Button } from '~/components/ui/button';
import { charactersAtom, initStream, playersAtom } from '~/state';
import { Action, Character as CharacterType } from '~/types';

export function AdminView() {
  const characters = useAtomValue(charactersAtom);

  const [characterId, setCharacterId] = useState<number | 'new' | null>(null);
  const [actionId, setActionId] = useState<number | 'new' | null>(null);

  useEffect(() => {
    return initStream();
  }, []);

  if (characterId !== null) {
    const defaultValues =
      characters.find((c) => c.id === characterId) ?? getNewCharacter();
    return (
      <CharacterForm
        id={typeof characterId === 'number' ? characterId : undefined}
        defaultValues={defaultValues}
        onClose={() => setCharacterId(null)}
      />
    );
  }
  if (typeof characterId === 'number' && actionId !== null) {
    const defaultValues =
      characters
        .find((c) => c.id === characterId)
        ?.actions.find((a) => a.id === actionId) ?? getNewAction();
    return (
      <ActionForm
        id={typeof actionId === 'number' ? actionId : undefined}
        defaultValues={defaultValues}
        character_id={characterId}
        onClose={() => setActionId(null)}
      />
    );
  }

  return (
    <div className="flex">
      <div className="grow space-y-4">
        <ul>
          <h2>Characters</h2>
          <Button onClick={() => setCharacterId('new')}>
            Add Character <Plus />
          </Button>
          <ul>
            {characters.map((char) => (
              <li key={char.id}>
                <Character
                  character={char}
                  onEdit={(id: number) => setCharacterId(id)}
                />
              </li>
            ))}
          </ul>
        </ul>
      </div>
      <PlayersList />
    </div>
  );
}

function getNewCharacter(): Omit<CharacterType, 'id' | 'actions'> {
  return {
    name: '',
    race: '',
    gender: '',
    age: '',
    description: '',
    appearance: '',
  };
}

function getNewAction(): Omit<Action, 'id' | 'character_id'> {
  return {
    type: '',
    direction: '',
    content: '',
  };
}

function PlayersList() {
  const players = useAtomValue(playersAtom);

  return (
    <div className="space-y-4 rounded-sm">
      <h2 className="text-lg font-bold">Players</h2>
      <ul className="text-md w-64 space-y-2">
        {players.map((player) => (
          <li
            key={player.id}
            data-online={player.is_online}
            className="group flex items-center justify-between space-x-2 rounded-sm px-4 py-2"
          >
            <div className="leading-0 flex h-full items-center space-x-2">
              <div
                title={player.is_online ? 'Online' : 'Offline'}
                className="h-3 w-3 rounded-full bg-gray-500 group-data-[online=true]:bg-teal-500"
              />
              <span className="text-muted-foreground group-data-[online=true]:font-bold group-data-[online=true]:text-secondary-foreground">
                {player.name}
              </span>
            </div>
            <Button
              size="icon"
              variant="outline"
              onClick={() => NpcSurpriseApi.deletePlayer(player.id)}
            >
              <X className="h-4 w-4" />
            </Button>
          </li>
        ))}
      </ul>
    </div>
  );
}
