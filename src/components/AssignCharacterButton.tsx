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
        <DropdownMenuItem onClick={() => NpcSurpriseApi.unassignCharacter(id)}>
          Unassign
        </DropdownMenuItem>
        <DropdownMenuSeparator />
        {players.map((player) => (
          <DropdownMenuItem
            key={player.id}
            onClick={() => NpcSurpriseApi.assignCharacter(id, player.id)}
          >
            {player.name}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
