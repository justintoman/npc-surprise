import { AssignButton } from '~/components/AssignButton';
import { Button } from '~/components/ui/button';
import type { Character } from '~/types';

type Props = {
  character: Character;
  onEdit?: (id: number) => void;
};

export function Character({ character, onEdit }: Props) {
  return (
    <div>
      {onEdit ? (
        <>
          <Button onClick={() => onEdit(character.id)}>Edit</Button>
          <AssignButton type="character" id={character.id} />
        </>
      ) : null}
      <h3>name: {character.name ?? 'hidden'}</h3>
      <p>race: {(character.race ?? 'hidden') || '-'}</p>
      <p>gender: {(character.gender ?? 'hidden') || '-'}</p>
      <p>age: {(character.age ?? 'hidden') || '-'}</p>
      <p>description: {(character.description ?? 'hidden') || '-'}</p>
      <p>appearance: {(character.appearance ?? 'hidden') || '-'}</p>

      {character.actions.length ? (
        <>
          <h3>Actions</h3>
          <ul>
            {character.actions.map((action) => (
              <li key={action.id}>
                <p>{action.type}</p>
                <p>{action.direction}</p>
                <p>{action.content}</p>
              </li>
            ))}
          </ul>
        </>
      ) : null}
    </div>
  );
}
