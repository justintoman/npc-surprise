import { useAtomValue } from 'jotai';
import { PlusCircle } from 'lucide-react';
import { useState } from 'react';
import { ActionForm } from '~/AdminView/ActionForm';
import { CharacterForm } from '~/AdminView/CharacterForm';
import { Action } from '~/components/Action';
import { AssignCharacterButton } from '~/components/AssignCharacterButton';
import { Button } from '~/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTrigger,
} from '~/components/ui/dialog';
import { getNewAction } from '~/lib/utils';
import { isAdminAtom } from '~/state';
import type { Character } from '~/types';

type Props = {
  character: Character;
};

export function Character({ character }: Props) {
  const isAdmin = useAtomValue(isAdminAtom);
  const [isEditing, setIsEditing] = useState(false);
  const [isAddingAction, setIsAddingAction] = useState(false);

  return (
    <div className="space-y-2">
      {isAdmin ? (
        <div className="flex items-center justify-start space-x-4">
          <AssignCharacterButton id={character.id} />
          <Dialog open={isEditing} onOpenChange={setIsEditing}>
            <DialogTrigger asChild>
              <Button size="sm">Edit</Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>Edit Character</DialogHeader>
              <CharacterForm
                id={character.id}
                defaultValues={character}
                onClose={() => setIsEditing(false)}
              />
            </DialogContent>
          </Dialog>
        </div>
      ) : null}
      <h3>name: {character.name ?? 'hidden'}</h3>
      <p>race: {(character.race ?? 'hidden') || '-'}</p>
      <p>gender: {(character.gender ?? 'hidden') || '-'}</p>
      <p>age: {(character.age ?? 'hidden') || '-'}</p>
      <p>description: {(character.description ?? 'hidden') || '-'}</p>
      <p>appearance: {(character.appearance ?? 'hidden') || '-'}</p>

      <header className="flex items-center space-x-4">
        <h3 className="text-lg font-bold">Actions</h3>
        {isAdmin ? (
          <Dialog open={isAddingAction} onOpenChange={setIsAddingAction}>
            <DialogTrigger asChild>
              <Button size="icon">
                <PlusCircle className="h-4 w-4" />
              </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>Add Action</DialogHeader>
              <ActionForm
                characterId={character.id}
                defaultValues={getNewAction()}
                onClose={() => setIsAddingAction(false)}
              />
            </DialogContent>
          </Dialog>
        ) : null}
      </header>
      {character.actions.length ? (
        <ul>
          {character.actions.map((action) => (
            <Action key={action.id} action={action} />
          ))}
        </ul>
      ) : (
        <p className="text-sm italic">You don't have any actions yet</p>
      )}
    </div>
  );
}
