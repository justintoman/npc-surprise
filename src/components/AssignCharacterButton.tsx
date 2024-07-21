import { useAtomValue } from 'jotai';
import { NpcSurpriseApi } from '~/api';
import { Button } from '~/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '~/components/ui/dropdown-menu';
import { playersAtom } from '~/state';

type Props = {
  id: number;
};

export function AssignCharacterButton({ id }: Props) {
  const players = useAtomValue(playersAtom);
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button size="sm" variant="secondary">
          Assign
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent>
        <DropdownMenuItem
          className="cursor-pointer"
          onClick={() => NpcSurpriseApi.unassignCharacter(id)}
        >
          Unassign
        </DropdownMenuItem>
        <DropdownMenuSeparator />
        {players.map((player) => (
          <DropdownMenuItem
            key={player.id}
            onClick={() => NpcSurpriseApi.assignCharacter(id, player.id)}
            data-online={player.isOnline}
            className="group flex cursor-pointer items-center justify-between rounded-sm"
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
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
