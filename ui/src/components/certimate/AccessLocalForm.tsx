import {
  Access,
  accessFormType,
  getUsageByConfigType,
  LocalConfig,
  SSHConfig,
} from "@/domain/access";
import { useConfig } from "@/providers/config";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "../ui/form";
import { Input } from "../ui/input";
import { Button } from "../ui/button";
import { Textarea } from "../ui/textarea";
import { save } from "@/repository/access";
import { ClientResponseError } from "pocketbase";
import { PbErrorData } from "@/domain/base";

const AccessLocalForm = ({
  data,
  onAfterReq,
}: {
  data?: Access;
  onAfterReq: () => void;
}) => {
  const { addAccess, updateAccess, reloadAccessGroups } = useConfig();

  const formSchema = z.object({
    id: z.string().optional(),
    name: z.string().min(1).max(64),
    configType: accessFormType,

    command: z.string().min(1).max(2048),
    certPath: z.string().min(0).max(2048),
    keyPath: z.string().min(0).max(2048),
  });

  let config: LocalConfig = {
    command: "sudo service nginx restart",
    certPath: "/etc/nginx/ssl/certificate.crt",
    keyPath: "/etc/nginx/ssl/private.key",
  };
  if (data) config = data.config as SSHConfig;

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      id: data?.id,
      name: data?.name,
      configType: "local",
      certPath: config.certPath,
      keyPath: config.keyPath,
      command: config.command,
    },
  });

  const onSubmit = async (data: z.infer<typeof formSchema>) => {
    const req: Access = {
      id: data.id as string,
      name: data.name,
      configType: data.configType,
      usage: getUsageByConfigType(data.configType),

      config: {
        command: data.command,
        certPath: data.certPath,
        keyPath: data.keyPath,
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
      } else {
        addAccess(req);
      }

      reloadAccessGroups();
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

      return;
    }
  };

  return (
    <>
      <div className="max-w-[35em] mx-auto mt-10">
        <Form {...form}>
          <form
            onSubmit={(e) => {
              e.stopPropagation();
              form.handleSubmit(onSubmit)(e);
            }}
            className="space-y-3"
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
              name="certPath"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>证书保存路径</FormLabel>
                  <FormControl>
                    <Input placeholder="请输入证书上传路径" {...field} />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="keyPath"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>私钥保存路径</FormLabel>
                  <FormControl>
                    <Input placeholder="请输入私钥上传路径" {...field} />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="command"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Command</FormLabel>
                  <FormControl>
                    <Textarea placeholder="请输入要执行的命令" {...field} />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />

            <FormMessage />

            <div className="flex justify-end">
              <Button type="submit">保存</Button>
            </div>
          </form>
        </Form>
      </div>
    </>
  );
};

export default AccessLocalForm;
