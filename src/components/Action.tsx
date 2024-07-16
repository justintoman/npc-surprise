import type { Action } from '~/types';

export function Action({ action }: { action: Action }) {
  return (
    <div>
      <p>{action.type}</p>
      <p>{action.direction}</p>
      <p>{action.content}</p>
    </div>
  );
}
