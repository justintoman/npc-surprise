import { zodResolver } from '@hookform/resolvers/zod';
import { useSetAtom } from 'jotai';
import { Check } from 'lucide-react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { NpcSurpriseApi } from '~/api';
import { Button } from '~/components/ui/button';
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '~/components/ui/form';
import { Input } from '~/components/ui/input';
import { statusAtom } from '~/state';

export function PlayerLogin({ name }: { name: string }) {
  const methods = useForm<LoginForm>({
    resolver: zodResolver(schema),
    defaultValues: { name },
  });
  const refreshStatus = useSetAtom(statusAtom);
  async function submit({ name }: LoginForm) {
    try {
      await NpcSurpriseApi.login(name);
    } catch (error) {
      console.error(error);
    }
    refreshStatus();
  }

  return (
    <Form {...methods}>
      <form onSubmit={methods.handleSubmit(submit)} className="p-4">
        <FormField
          control={methods.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Player Name</FormLabel>
              <div className="flex items-center space-x-2">
                <FormControl>
                  <Input placeholder="Put yer name here" {...field} />
                </FormControl>
                <Button type="submit" size="icon" className="shrink-0">
                  <Check className="h-4 w-4" />
                </Button>
              </div>
              <FormDescription>
                Only Justin will be able to see this for now.
              </FormDescription>
              <FormMessage />
            </FormItem>
          )}
        />
      </form>
    </Form>
  );
}

const schema = z.object({
  name: z.string().min(1, 'Name is required'),
});

type LoginForm = z.infer<typeof schema>;
