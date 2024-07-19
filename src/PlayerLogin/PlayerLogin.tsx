import { zodResolver } from '@hookform/resolvers/zod';
import { useSetAtom } from 'jotai';
import { Check } from 'lucide-react';
import { useForm } from 'react-hook-form';
import { useNavigate } from 'react-router-dom';
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
import { statusAtom } from '~/state';

export function PlayerLogin({ name }: { name?: string }) {
  const methods = useForm<LoginForm>({
    resolver: zodResolver(schema),
    defaultValues: { name },
  });
  const setStatus = useSetAtom(statusAtom);
  const navigate = useNavigate();
  async function submit({ name }: LoginForm) {
    try {
      const status = await NpcSurpriseApi.login(name);
      setStatus(status);
      if (status.isAdmin) {
        navigate('/admin');
      } else {
        navigate('/player');
      }
    } catch (error) {
      console.error(error);
    }
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
                  <Input {...field} />
                </FormControl>
                <Button type="submit" size="icon" className="shrink-0">
                  <Check className="h-4 w-4" />
                </Button>
              </div>
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
