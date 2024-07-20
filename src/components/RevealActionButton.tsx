import { NpcSurpriseApi } from '~/api';
import { Button } from '~/components/ui/button';

type Props = {
  characterId: number;
  actionId: number;
  revealed: boolean;
};

export function RevealActionButton({ characterId, revealed, actionId }: Props) {
  return (
    <Button
      size="sm"
      variant="secondary"
      onClick={() => {
        if (revealed) {
          NpcSurpriseApi.hideAction(characterId, actionId);
        } else {
          NpcSurpriseApi.revealAction(characterId, actionId);
        }
      }}
    >
      {revealed ? 'Hide' : 'Reveal'}
    </Button>
  );
}
