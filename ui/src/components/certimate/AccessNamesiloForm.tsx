import { Input } from "@/components/ui/input";
import { useForm } from "react-hook-form";

import { zodResolver } from "@hookform/resolvers/zod";

import z from "zod";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Button } from "@/components/ui/button";

import { Access, accessFormType, getUsageByConfigType, NamesiloConfig } from "@/domain/access";
import { save } from "@/repository/access";
import { useConfig } from "@/providers/config";
import { ClientResponseError } from "pocketbase";
import { PbErrorData } from "@/domain/base";

const AccessNamesiloForm = ({
  data,
  onAfterReq,
}: {
  data?: Access;
  onAfterReq: () => void;
}) => {
  const { addAccess, updateAccess } = useConfig();
  const formSchema = z.object({
    id: z.string().optional(),
    name: z.string().min(1).max(64),
    sslprovider: z.string().min(1).max(64),
    configType: accessFormType,
    apiKey: z.string().min(1).max(64),
  });

  let config: NamesiloConfig = {
    apiKey: "",
  };
  if (data) config = data.config as NamesiloConfig;

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      id: data?.id,
      name: data?.name,
      configType: "namesilo",
      apiKey: config.apiKey,
    },
  });

  const onSubmit = async (data: z.infer<typeof formSchema>) => {
    console.log(data);
    const req: Access = {
      id: data.id as string,
      name: data.name,
      configType: data.configType,
      usage: getUsageByConfigType(data.configType),
      config: {
        apiKey: data.apiKey,
      },
    };

    try {
      const rs = await save(req);

      onAfterReq();

      req.id = rs.id;
      req.created = rs.created;
      req.updated = rs.updated;
      if (data.id) {
        updateAccess(req);
        return;
      }
      addAccess(req);
    } catch (e) {
      const err = e as ClientResponseError;

      Object.entries(err.response.data as PbErrorData).forEach(
        ([key, value]) => {
          form.setError(key as keyof z.infer<typeof formSchema>, {
            type: "manual",
            message: value.message,
          });
        }
      );
    }
  };

  return (
    <>
      <div className="max-w-[35em] mx-auto mt-10">
        <Form {...form}>
          <form
            onSubmit={(e) => {
              console.log(e);
              e.stopPropagation();
              form.handleSubmit(onSubmit)(e);
            }}
            className="space-y-8"
          >
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>名称</FormLabel>
                  <FormControl>
                    <Input placeholder="请输入授权名称" {...field} />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="id"
              render={({ field }) => (
                <FormItem className="hidden">
                  <FormLabel>配置类型</FormLabel>
                  <FormControl>
                    <Input {...field} />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="configType"
              render={({ field }) => (
                <FormItem className="hidden">
                  <FormLabel>配置类型</FormLabel>
                  <FormControl>
                    <Input {...field} />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="apiKey"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>NAMESILO_API_KEY</FormLabel>
                  <FormControl>
                    <Input placeholder="请输入NAMESILO_API_KEY" {...field} />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />

            <div className="flex justify-end">
              <Button type="submit">保存</Button>
            </div>
          </form>
        </Form>
      </div>
    </>
  );
};

export default AccessNamesiloForm;