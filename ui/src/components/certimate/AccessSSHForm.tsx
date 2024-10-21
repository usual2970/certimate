import { useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Plus } from "lucide-react";
import { ClientResponseError } from "pocketbase";

import { Access, accessFormType, getUsageByConfigType, SSHConfig } from "@/domain/access";
import { Button } from "@/components/ui/button";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import AccessGroupEdit from "./AccessGroupEdit";
import { readFileContent } from "@/lib/file";
import { cn } from "@/lib/utils";
import { PbErrorData } from "@/domain/base";
import { save } from "@/repository/access";
import { updateById } from "@/repository/access_group";
import { useConfig } from "@/providers/config";

type AccessSSHFormProps = {
  op: "add" | "edit" | "copy";
  data?: Access;
  onAfterReq: () => void;
};

const AccessSSHForm = ({ data, op, onAfterReq }: AccessSSHFormProps) => {
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
      .min(1, "access.authorization.form.name.placeholder")
      .max(64, t("common.errmsg.string_max", { max: 64 })),
    configType: accessFormType,
    host: z.string().refine(
      (str) => {
        return ipReg.test(str) || domainReg.test(str);
      },
      {
        message: "common.errmsg.host_invalid",
      }
    ),
    group: z.string().optional(),
    port: z
      .string()
      .min(1, "access.authorization.form.ssh_port.placeholder")
      .max(5, t("common.errmsg.string_max", { max: 5 })),
    username: z
      .string()
      .min(1, "access.authorization.form.ssh_username.placeholder")
      .max(64, t("common.errmsg.string_max", { max: 64 })),
    password: z
      .string()
      .min(0, "access.authorization.form.ssh_password.placeholder")
      .max(64, t("common.errmsg.string_max", { max: 64 })),
    key: z
      .string()
      .min(0, "access.authorization.form.ssh_key.placeholder")
      .max(20480, t("common.errmsg.string_max", { max: 20480 })),
    keyFile: z.any().optional(),
    keyPassphrase: z
      .string()
      .min(0, "access.authorization.form.ssh_key_passphrase.placeholder")
      .max(2048, t("common.errmsg.string_max", { max: 2048 })),
  });

  let config: SSHConfig = {
    host: "127.0.0.1",
    port: "22",
    username: "root",
    password: "",
    key: "",
    keyFile: "",
    keyPassphrase: "",
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
      keyPassphrase: config.keyPassphrase,
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
        keyPassphrase: data.keyPassphrase,
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

      Object.entries(err.response.data as PbErrorData).forEach(([key, value]) => {
        form.setError(key as keyof z.infer<typeof formSchema>, {
          type: "manual",
          message: value.message,
        });
      });

      return;
    }
  };

  const handleFileChange = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;
    const savedFile = file;
    setFileName(savedFile.name);
    const content = await readFileContent(savedFile);
    form.setValue("key", content);
  };

  const handleSelectFileClick = () => {
    fileInputRef.current?.click();
  };

  return (
    <>
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
                <FormLabel>{t("access.authorization.form.name.label")}</FormLabel>
                <FormControl>
                  <Input placeholder={t("access.authorization.form.name.placeholder")} {...field} />
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
                  <div>{t("access.authorization.form.ssh_group.label")}</div>
                  <AccessGroupEdit
                    trigger={
                      <div className="font-normal text-primary hover:underline cursor-pointer flex items-center">
                        <Plus size={14} />
                        {t("common.add")}
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
                      <SelectValue placeholder={t("access.authorization.form.access_group.placeholder")} />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="emptyId">
                        <div className={cn("flex items-center space-x-2 rounded cursor-pointer")}>--</div>
                      </SelectItem>
                      {accessGroups.map((item) => (
                        <SelectItem value={item.id ? item.id : ""} key={item.id}>
                          <div className={cn("flex items-center space-x-2 rounded cursor-pointer")}>{item.name}</div>
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
                <FormLabel>{t("access.authorization.form.config.label")}</FormLabel>
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
                <FormLabel>{t("access.authorization.form.config.label")}</FormLabel>
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
                  <FormLabel>{t("access.authorization.form.ssh_host.label")}</FormLabel>
                  <FormControl>
                    <Input placeholder={t("access.authorization.form.ssh_host.placeholder")} {...field} />
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
                  <FormLabel>{t("access.authorization.form.ssh_port.label")}</FormLabel>
                  <FormControl>
                    <Input placeholder={t("access.authorization.form.ssh_port.placeholder")} {...field} type="number" />
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
                <FormLabel>{t("access.authorization.form.ssh_username.label")}</FormLabel>
                <FormControl>
                  <Input placeholder={t("access.authorization.form.ssh_username.placeholder")} {...field} />
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
                <FormLabel>{t("access.authorization.form.ssh_password.label")}</FormLabel>
                <FormControl>
                  <Input placeholder={t("access.authorization.form.ssh_password.placeholder")} {...field} type="password" />
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
                <FormLabel>{t("access.authorization.form.ssh_key.label")}</FormLabel>
                <FormControl>
                  <Input placeholder={t("access.authorization.form.ssh_key.placeholder")} {...field} />
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
                <FormLabel>{t("access.authorization.form.ssh_key.label")}</FormLabel>
                <FormControl>
                  <div>
                    <Button type={"button"} variant={"secondary"} size={"sm"} className="w-48" onClick={handleSelectFileClick}>
                      {fileName ? fileName : t("access.authorization.form.ssh_key_file.placeholder")}
                    </Button>
                    <Input
                      placeholder={t("access.authorization.form.ssh_key.placeholder")}
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
            name="keyPassphrase"
            render={({ field }) => (
              <FormItem>
                <FormLabel>{t("access.authorization.form.ssh_key_passphrase.label")}</FormLabel>
                <FormControl>
                  <Input placeholder={t("access.authorization.form.ssh_key_passphrase.placeholder")} {...field} type="password" />
                </FormControl>

                <FormMessage />
              </FormItem>
            )}
          />

          <FormMessage />

          <div className="flex justify-end">
            <Button type="submit">{t("common.save")}</Button>
          </div>
        </form>
      </Form>
    </>
  );
};

export default AccessSSHForm;
