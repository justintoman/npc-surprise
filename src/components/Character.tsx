import { useAtomValue } from 'jotai';
import { Link } from 'react-router-dom';
import { RevealFieldButton } from '~/Admin/RevealFieldButton';
import { Action } from '~/components/Action';
import { AssignCharacterButton } from '~/components/AssignCharacterButton';
import { Button } from '~/components/ui/button';
import { Label } from '~/components/ui/label';
import { playerAtomFamily } from '~/state';
import type { Character } from '~/types';

type Props = {
  character: Character;
  isAdmin?: boolean;
};

export function Character({ character, isAdmin }: Props) {
  const player = useAtomValue(playerAtomFamily(character.playerId));
  return (
    <div className="max-w-lg">
      <div className="space-y-6 p-4">
        <header className="flex justify-between">
          <div className="flex items-center space-x-2">
            <h1 className="text-bold text-2xl">
              {character.name || 'Name hidden'}
            </h1>

            <RevealFieldButton field="name" characterId={character.id} />
          </div>

          {isAdmin ? (
            <div className="flex items-center justify-start space-x-4">
              <AssignCharacterButton id={character.id} />
              <Button size="sm" asChild>
                <Link to={`/admin/character/${character.id}`}>Edit</Link>
              </Button>
            </div>
          ) : null}
        </header>

        <div className="flex justify-between">
          <div className="flex space-x-6">
            <div>
              <div className="flex items-center space-x-2">
                <Label className="text-sm font-bold">Age</Label>
                <RevealFieldButton field="age" characterId={character.id} />
              </div>
              <p>{character.age || 'hidden'}</p>
            </div>
            <div>
              <div className="flex items-center space-x-2">
                <Label className="text-sm font-bold">Race</Label>
                <RevealFieldButton field="race" characterId={character.id} />
              </div>
              <p>{character.race || 'hidden'}</p>
            </div>
            <div>
              <div className="flex items-center space-x-2">
                <Label className="text-sm font-bold">Gender</Label>
                <RevealFieldButton field="gender" characterId={character.id} />
              </div>
              <p>{character.gender || 'hidden'}</p>
            </div>
          </div>
          {player && isAdmin ? (
            <div className="rounded-sm bg-secondary p-2">
              <Label className="text-sm font-bold">Assigned to</Label>
              <p>{player.name}</p>
            </div>
          ) : null}
        </div>
        <div>
          <div className="flex items-center space-x-2">
            <Label className="text-sm font-bold">Appearance</Label>
            <RevealFieldButton field="appearance" characterId={character.id} />
          </div>
          <p>{character.appearance || 'hidden'}</p>
        </div>
        <div>
          <div className="flex items-center space-x-2">
            <Label className="text-sm font-bold">Description</Label>
            <RevealFieldButton field="description" characterId={character.id} />
          </div>
          <p>{character.description || 'hidden'}</p>
        </div>
      </div>

      <header className="mt-8 flex items-center space-x-4 px-4">
        <h3 className="text-lg font-bold">Actions</h3>
        {isAdmin ? (
          <Button size="icon" asChild>
            <Link to={`/admin/${character.id}/action/new`}>Add</Link>
          </Button>
        ) : null}
      </header>
      {character.actions.length ? (
        <ul className="divide-y">
          {character.actions.map((action) => (
            <li className="p-4" key={action.id}>
              <Action action={action} isAdmin={isAdmin} />
            </li>
          ))}
        </ul>
      ) : (
        <p className="p-4 text-sm italic">You don't have any actions yet</p>
      )}
    </div>
  );
}
