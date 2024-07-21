import { useNavigate, useParams } from 'react-router-dom';
import { NpcSurpriseApi } from '~/api';
import { getNewAction } from '~/lib/utils';
import { ActionForm } from './ActionForm';

export function NewAction() {
  const { characterId } = useParams();
  if (!characterId) {
    throw new Error('characterId is required');
  }
  const navigate = useNavigate();

  function close() {
    navigate(`/admin`);
  }

  return (
    <ActionForm
      defaultValues={getNewAction()}
      onClose={() => close()}
      onSubmit={async (data) => {
        await NpcSurpriseApi.createAction({
          characterId: Number(characterId),
          revealed: false,
          ...data,
        });
        close();
      }}
    />
  );
}
