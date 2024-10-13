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
import { save } from "@/repository/access";
import { ClientResponseError } from "pocketbase";
import { PbErrorData } from "@/domain/base";
import { readFileContent } from "@/lib/file";
import { useRef, useState } from "react";
import { useTranslation } from "react-i18next";
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
  op,
  onAfterReq,
}: {
  data?: Access;
  op: "add" | "edit" | "copy";
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
  const { t } = useTranslation();

  const originGroup = data ? (data.group ? data.group : "") : "";

  const domainReg = /^(?:\*\.)?([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$/;
  const ipReg =
    /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;

  const formSchema = z.object({
    id: z.string().optional(),
    name: z
      .string()
      .min(1, "access.form.name.not.empty")
      .max(64, t("zod.rule.string.max", { max: 64 })),
    configType: accessFormType,
    host: z.string().refine(
      (str) => {
        return ipReg.test(str) || domainReg.test(str);
      },
      {
        message: "zod.rule.ssh.host",
      }
    ),
    group: z.string().optional(),
    port: z
      .string()
      .min(1, "access.form.ssh.port.not.empty")
      .max(5, t("zod.rule.string.max", { max: 5 })),
    username: z
      .string()
      .min(1, "username.not.empty")
      .max(64, t("zod.rule.string.max", { max: 64 })),
    password: z
      .string()
      .min(0, "password.not.empty")
      .max(64, t("zod.rule.string.max", { max: 64 })),
    key: z
      .string()
      .min(0, "access.form.ssh.key.not.empty")
      .max(20480, t("zod.rule.string.max", { max: 20480 })),
    keyFile: z.any().optional(),
  });

  let config: SSHConfig = {
    host: "127.0.0.1",
    port: "22",
    username: "root",
    password: "",
    key: "",
    keyFile: "",
  };
  if (data) config = data.config as SSHConfig;

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      id: data?.id,
      name: data?.name || "",
      configType: "ssh",
      group: data?.group,
      host: config.host,
      port: config.port,
      username: config.username,
      password: config.password,
      key: config.key,
      keyFile: config.keyFile,
    },
  });

  const onSubmit = async (data: z.infer<typeof formSchema>) => {
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
      },
    };

    try {
      req.id = op == "copy" ? "" : req.id;
      const rs = await save(req);

      onAfterReq();

      req.id = rs.id;
      req.created = rs.created;
      req.updated = rs.updated;
      if (data.id && op == "edit") {
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
                  <FormLabel>{t("name")}</FormLabel>
                  <FormControl>
                    <Input
                      placeholder={t("access.form.name.not.empty")}
                      {...field}
                    />
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
                    <div>{t("access.form.ssh.group.label")}</div>
                    <AccessGroupEdit
                      trigger={
                        <div className="font-normal text-primary hover:underline cursor-pointer flex items-center">
                          <Plus size={14} />
                          {t("add")}
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
                        <SelectValue
                          placeholder={t("access.group.not.empty")}
                        />
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
                  <FormLabel>{t("access.form.config.field")}</FormLabel>
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
                  <FormLabel>{t("access.form.config.field")}</FormLabel>
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
                    <FormLabel>{t("access.form.ssh.host")}</FormLabel>
                    <FormControl>
                      <Input
                        placeholder={t("access.form.ssh.host.not.empty")}
                        {...field}
                      />
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
                    <FormLabel>{t("access.form.ssh.port")}</FormLabel>
                    <FormControl>
                      <Input
                        placeholder={t("access.form.ssh.port.not.empty")}
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
                  <FormLabel>{t("username")}</FormLabel>
                  <FormControl>
                    <Input placeholder={t("username.not.empty")} {...field} />
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
                  <FormLabel>{t("password")}</FormLabel>
                  <FormControl>
                    <Input
                      placeholder={t("password.not.empty")}
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
                  <FormLabel>{t("access.form.ssh.key")}</FormLabel>
                  <FormControl>
                    <Input
                      placeholder={t("access.form.ssh.key.not.empty")}
                      {...field}
                    />
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
                  <FormLabel>{t("access.form.ssh.key")}</FormLabel>
                  <FormControl>
                    <div>
                      <Button
                        type={"button"}
                        variant={"secondary"}
                        size={"sm"}
                        className="w-48"
                        onClick={handleSelectFileClick}
                      >
                        {fileName
                          ? fileName
                          : t("access.form.ssh.key.file.not.empty")}
                      </Button>
                      <Input
                        placeholder={t("access.form.ssh.key.not.empty")}
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

            <FormMessage />

            <div className="flex justify-end">
              <Button type="submit">{t("save")}</Button>
            </div>
          </form>
        </Form>
      </div>
    </>
  );
};

export default AccessSSHForm;
