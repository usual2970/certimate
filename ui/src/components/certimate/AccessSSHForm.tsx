import {
  Access,
  accessFormType,
  getUsageByConfigType,
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
import { readFileContent } from "@/lib/file";
import { useRef, useState } from "react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../ui/select";
import { cn } from "@/lib/utils";
import AccessGroupEdit from "./AccessGroupEdit";
import { Plus } from "lucide-react";
import { updateById } from "@/repository/access_group";

const AccessSSHForm = ({
  data,
  onAfterReq,
}: {
  data?: Access;
  onAfterReq: () => void;
}) => {
  const {
    addAccess,
    updateAccess,
    reloadAccessGroups,
    config: { accessGroups },
  } = useConfig();

  const fileInputRef = useRef<HTMLInputElement | null>(null);

  const [fileName, setFileName] = useState("");

  const originGroup = data ? (data.group ? data.group : "") : "";

  const domainReg = /^(?:\*\.)?([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$/;
  const ipReg =
    /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;

  const formSchema = z.object({
    id: z.string().optional(),
    name: z.string().min(1).max(64),
    configType: accessFormType,
    host: z.string().refine(
      (str) => {
        return ipReg.test(str) || domainReg.test(str);
      },
      {
        message: "请输入正确的域名或IP",
      }
    ),
    group: z.string().optional(),
    port: z.string().min(1).max(5),
    username: z.string().min(1).max(64),
    password: z.string().min(0).max(64),
    key: z.string().min(0).max(20480),
    keyFile: z.any().optional(),
    command: z.string().min(1).max(2048),
    preCommand: z.string().min(0).max(2048).optional(),
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
    preCommand: "",
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
      group: data?.group,
      host: config.host,
      port: config.port,
      username: config.username,
      password: config.password,
      key: config.key,
      keyFile: config.keyFile,
      certPath: config.certPath,
      keyPath: config.keyPath,
      command: config.command,
      preCommand: config.preCommand,
    },
  });

  const onSubmit = async (data: z.infer<typeof formSchema>) => {
    console.log(data);
    let group = data.group;
    if (group == "emptyId") group = "";

    const req: Access = {
      id: data.id as string,
      name: data.name,
      configType: data.configType,
      usage: getUsageByConfigType(data.configType),
      group: group,
      config: {
        host: data.host,
        port: data.port,
        username: data.username,
        password: data.password,
        key: data.key,
        command: data.command,
        preCommand: data.preCommand,
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

      // 同步更新授权组
      if (group != originGroup) {
        if (originGroup) {
          await updateById({
            id: originGroup,
            "access-": req.id,
          });
        }

        if (group) {
          await updateById({
            id: group,
            "access+": req.id,
          });
        }
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

  const handleFileChange = async (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    const file = event.target.files?.[0];
    if (!file) return;
    const savedFile = file;
    setFileName(savedFile.name);
    const content = await readFileContent(savedFile);
    form.setValue("key", content);
  };

  const handleSelectFileClick = () => {
    console.log(fileInputRef.current);

    fileInputRef.current?.click();
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
              name="group"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="w-full flex justify-between">
                    <div>授权配置组(用于将一个域名证书部署到多个 ssh 主机)</div>
                    <AccessGroupEdit
                      trigger={
                        <div className="font-normal text-primary hover:underline cursor-pointer flex items-center">
                          <Plus size={14} />
                          新增
                        </div>
                      }
                    />
                  </FormLabel>
                  <FormControl>
                    <Select
                      {...field}
                      value={field.value}
                      defaultValue="emptyId"
                      onValueChange={(value) => {
                        form.setValue("group", value);
                      }}
                    >
                      <SelectTrigger>
                        <SelectValue placeholder="请选择分组" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="emptyId">
                          <div
                            className={cn(
                              "flex items-center space-x-2 rounded cursor-pointer"
                            )}
                          >
                            --
                          </div>
                        </SelectItem>
                        {accessGroups.map((item) => (
                          <SelectItem
                            value={item.id ? item.id : ""}
                            key={item.id}
                          >
                            <div
                              className={cn(
                                "flex items-center space-x-2 rounded cursor-pointer"
                              )}
                            >
                              {item.name}
                            </div>
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
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
                    <FormLabel>服务器HOST</FormLabel>
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
                    <div>
                      <Button
                        type={"button"}
                        variant={"secondary"}
                        size={"sm"}
                        className="w-48"
                        onClick={handleSelectFileClick}
                      >
                        {fileName ? fileName : "请选择文件"}
                      </Button>
                      <Input
                        placeholder="请输入Key"
                        {...field}
                        ref={fileInputRef}
                        className="hidden"
                        hidden
                        type="file"
                        onChange={handleFileChange}
                      />
                    </div>
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
              name="preCommand"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>前置 Command</FormLabel>
                  <FormControl>
                    <Textarea placeholder="请输入要在部署证书前执行的前置命令" {...field} />
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
