import { useAtomValue } from 'jotai';
import { NpcSurpriseApi } from '~/api';
import { Button } from '~/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
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
        {players.map((player) => (
          <DropdownMenuItem
            key={player.id}
            onClick={() => NpcSurpriseApi.assign('character', id, player.id)}
          >
            {player.name}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
