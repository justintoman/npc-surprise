import { useAtomValue } from 'jotai';
import { PlusCircle, X } from 'lucide-react';
import { Link } from 'react-router-dom';
import { NpcSurpriseApi } from '~/api';
import { Character } from '~/components/Character';
import { Button } from '~/components/ui/button';
import { charactersAtom, playersAtom } from '~/state';

export function AdminHome() {
  const characters = useAtomValue(charactersAtom);

  return (
    <div className="flex">
      <div className="grow space-y-4">
        <header className="flex items-center space-x-4">
          <Button asChild>
            <Link to="/admin/character/new" className="flex items-center gap-2">
              Add Character
              <PlusCircle className="h-4 w-4" />
            </Link>
          </Button>
        </header>
        <ul className="grid grid-cols-2 gap-4">
          {characters.map((char) => (
            <li key={char.id}>
              <div className="fit-content rounded-lg border border-secondary bg-secondary/30">
                <Character character={char} isAdmin />
              </div>
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
            data-online={player.isOnline}
            className="group flex items-center justify-between space-x-2 rounded-sm px-4 py-1"
          >
            <div className="leading-0 flex h-full items-center space-x-2">
              <div
                title={player.isOnline ? 'Online' : 'Offline'}
                className="h-3 w-3 rounded-full bg-muted group-data-[online=true]:bg-teal-500"
              />
              <span className="text-muted-foreground group-data-[online=true]:font-bold group-data-[online=true]:text-secondary-foreground">
                {player.name}
              </span>
            </div>
            <Button
              size="tiny"
              variant="outline"
              onClick={() => NpcSurpriseApi.deletePlayer(player.id)}
            >
              <X className="h-3 w-3" />
            </Button>
          </li>
        ))}
      </ul>
    </div>
  );
}
