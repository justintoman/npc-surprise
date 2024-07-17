import { useAtomValue } from 'jotai';
import { PlusCircle, X } from 'lucide-react';
import { useEffect, useState } from 'react';
import { CharacterForm } from '~/AdminView/CharacterForm';
import { NpcSurpriseApi } from '~/api';
import { Character } from '~/components/Character';
import { Button } from '~/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTrigger,
} from '~/components/ui/dialog';
import { getNewCharacter } from '~/lib/utils';
import { charactersAtom, initStream, playersAtom } from '~/state';

export function AdminView() {
  const characters = useAtomValue(charactersAtom);
  const [isAddingCharacter, setIsAddingCharacter] = useState(false);

  useEffect(() => {
    return initStream();
  }, []);

  return (
    <div className="flex">
      <div className="grow space-y-4">
        <header className="flex items-center space-x-4">
          <h2 className="text-lg font-bold">Characters</h2>
          <Dialog open={isAddingCharacter} onOpenChange={setIsAddingCharacter}>
            <DialogTrigger asChild>
              <Button size="icon">
                <PlusCircle className="h-4 w-4" />
              </Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>Add Character</DialogHeader>
              <CharacterForm
                defaultValues={getNewCharacter()}
                onClose={() => setIsAddingCharacter(false)}
              />
            </DialogContent>
          </Dialog>
        </header>
        <ul>
          {characters.map((char) => (
            <li key={char.id}>
              <Character character={char} />
            </li>
          ))}
        </ul>
      </div>
      <PlayersList />
    </div>
  );
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
