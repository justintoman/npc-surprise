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
import { Textarea } from '~/components/ui/textarea';

type Props = {
  id?: number;
  defaultValues?: CharacterFormInput;
  onClose(): void;
};

export function CharacterForm({ id, defaultValues, onClose }: Props) {
  const methods = useForm<CharacterFormInput>({
    resolver: zodResolver(schema),
    defaultValues,
  });

  async function submit(data: CharacterFormInput) {
    if (id) {
      await NpcSurpriseApi.updateCharacter({ id, ...data });
    } else {
      await NpcSurpriseApi.createCharacter(data);
    }

    onClose();
  }

  return (
    <Form {...methods}>
      <form onSubmit={methods.handleSubmit(submit)}>
        <FormField
          control={methods.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Name</FormLabel>
              <FormControl>
                <Input {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={methods.control}
          name="race"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Race</FormLabel>
              <FormControl>
                <Input {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={methods.control}
          name="gender"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Gender</FormLabel>
              <FormControl>
                <Input {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={methods.control}
          name="age"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Age</FormLabel>
              <FormControl>
                <Input {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={methods.control}
          name="appearance"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Appearance</FormLabel>
              <FormControl>
                <Textarea {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={methods.control}
          name="description"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Description</FormLabel>
              <FormControl>
                <Textarea {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button variant="secondary" onClick={onClose}>
          Cancel
        </Button>
        <Button type="submit">Submit</Button>
      </form>
    </Form>
  );
}

type CharacterFormInput = z.infer<typeof schema>;

const schema = z.object({
  name: z.string(),
  race: z.string(),
  gender: z.string(),
  age: z.string(),
  description: z.string(),
  appearance: z.string(),
});
