import { useNavigate } from 'react-router-dom';
import { NpcSurpriseApi } from '~/api';
import { getNewCharacter } from '~/lib/utils';
import { CharacterForm } from './CharacterForm';

export function NewCharacter() {
  const navigate = useNavigate();

  function close() {
    navigate(`/admin`);
  }

  return (
    <CharacterForm
      defaultValues={getNewCharacter()}
      onClose={() => close()}
      onSubmit={async (data) => {
        await NpcSurpriseApi.createCharacter(data);
        close();
      }}
    />
  );
}
