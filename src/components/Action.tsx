import { Link } from 'react-router-dom';
import { ActionMarkdown } from '~/components/ActionMarkdown';
import { RevealActionButton } from '~/components/RevealActionButton';
import { Button } from '~/components/ui/button';
import type { Action } from '~/types';

type Props = {
  action: Action;
  isAdmin?: boolean;
};

export function Action({ action, isAdmin }: Props) {
  return (
    <div className="space-y-8">
      {isAdmin ? (
        <div className="flex items-center justify-end space-x-4">
          <RevealActionButton
            actionId={action.id}
            characterId={action.characterId}
            revealed={action.revealed}
          />
          <Button size="icon" asChild>
            <Link to={`/admin/${action.characterId}/action/${action.id}`}>
              Edit
            </Link>
          </Button>
        </div>
      ) : null}
      <ActionMarkdown>{action.content}</ActionMarkdown>
    </div>
  );
}
