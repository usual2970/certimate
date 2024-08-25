import { Access, accessFormType, SSHConfig } from "@/domain/access";
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
import { readFileContent } from "@/lib/file";

const AccessSSHForm = ({
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
    configType: accessFormType,
    host: z.string().ip({
      message: "请输入合法的IP地址",
    }),
    port: z.string().min(1).max(5),
    username: z.string().min(1).max(64),
    password: z.string().min(0).max(64),
    key: z.string().min(0).max(20480),
    keyFile: z.string().optional(),
    command: z.string().min(1).max(2048),
    certPath: z.string().min(0).max(2048),
    keyPath: z.string().min(0).max(2048),
  });

  let config: SSHConfig = {
    host: "127.0.0.1",
    port: "22",
    username: "root",
    password: "",
    key: "",
    keyFile: "",
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
      configType: "ssh",
      host: config.host,
      port: config.port,
      username: config.username,
      password: config.password,
      key: config.key,
      keyFile: config.keyFile,
      certPath: config.certPath,
      keyPath: config.keyPath,
      command: config.command,
    },
  });

  const onSubmit = async (data: z.infer<typeof formSchema>) => {
    console.log(data);
    const req: Access = {
      id: data.id as string,
      name: data.name,
      configType: data.configType,
      config: {
        host: data.host,
        port: data.port,
        username: data.username,
        password: data.password,
        key: data.key,
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

      return;
    }
  };

  const handleFileChange = async (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    const file = event.target.files?.[0];
    if (!file) return;
    const content = await readFileContent(file);
    form.setValue("key", content);
    form.setValue("keyFile", "");
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
            <div className="flex space-x-2">
              <FormField
                control={form.control}
                name="host"
                render={({ field }) => (
                  <FormItem className="grow">
                    <FormLabel>服务器IP</FormLabel>
                    <FormControl>
                      <Input placeholder="请输入Host" {...field} />
                    </FormControl>

                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="port"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>SSH端口</FormLabel>
                    <FormControl>
                      <Input
                        placeholder="请输入Port"
                        {...field}
                        type="number"
                      />
                    </FormControl>

                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <FormField
              control={form.control}
              name="username"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>用户名</FormLabel>
                  <FormControl>
                    <Input placeholder="请输入用户名" {...field} />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="password"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>密码</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="请输入密码"
                      {...field}
                      type="password"
                    />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="key"
              render={({ field }) => (
                <FormItem hidden>
                  <FormLabel>Key（使用证书登录）</FormLabel>
                  <FormControl>
                    <Input placeholder="请输入Key" {...field} />
                  </FormControl>

                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="keyFile"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Key（使用证书登录）</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="请输入Key"
                      {...field}
                      type="file"
                      onChange={handleFileChange}
                    />
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
                  <FormLabel>证书上传路径</FormLabel>
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
                  <FormLabel>私钥上传路径</FormLabel>
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

export default AccessSSHForm;
