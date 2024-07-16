import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { NpcSurpriseApi } from '~/api';
import { Button } from '~/components/ui/button';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '~/components/ui/form';
import { Input } from '~/components/ui/input';

type Props = {
  character_id: number;
  id?: number;
  defaultValues?: ActionFormInput;
  onClose(): void;
};

export function ActionForm({ id, character_id, defaultValues, onClose }: Props) {
  const methods = useForm<ActionFormInput>({
    resolver: zodResolver(schema),
    defaultValues,
  });

  async function submit(data: ActionFormInput) {
    if (id) {
      await NpcSurpriseApi.updateAction({ id, character_id, ...data });
    } else {
      await NpcSurpriseApi.createAction({ character_id, ...data });
    }

    onClose();
  }

  return (
    <Form {...methods}>
      <form onSubmit={methods.handleSubmit(submit)} className="p-4">
        <FormField
          control={methods.control}
          name="type"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Type</FormLabel>
              <FormControl>
                <Input {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={methods.control}
          name="content"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Content</FormLabel>
              <FormControl>
                <Input {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={methods.control}
          name="direction"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Direction</FormLabel>
              <FormControl>
                <Input {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button variant="secondary" onClick={onClose}>Cancel</Button>
        <Button type="submit">Submit</Button>
      </form>
    </Form>
  );
}

type ActionFormInput = z.infer<typeof schema>;

const schema = z.object({
  type: z.string().min(1, 'Type is required'),
  content: z.string().min(1, 'Content is required'),
  direction: z.string().min(1, 'Direction is required'),
});
