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
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useConfig } from "@/providers/config";
import { useEffect, useState } from "react";
import { Domain, targetTypeKeys, targetTypeMap } from "@/domain/domain";
import { save, get } from "@/repository/domains";
import { ClientResponseError } from "pocketbase";
import { PbErrorData } from "@/domain/base";
import { useToast } from "@/components/ui/use-toast";
import { Toaster } from "@/components/ui/toaster";
import { useLocation, useNavigate } from "react-router-dom";
import { Plus } from "lucide-react";
import { AccessEdit } from "@/components/certimate/AccessEdit";
import { accessTypeMap } from "@/domain/access";
import EmailsEdit from "@/components/certimate/EmailsEdit";
import { Textarea } from "@/components/ui/textarea";
import { cn } from "@/lib/utils";
import { EmailsSetting } from "@/domain/settings";
import { useTranslation } from "react-i18next";
import StringList from "@/components/certimate/StringList";

const Edit = () => {
  const {
    config: { accesses, emails, accessGroups },
  } = useConfig();

  const [domain, setDomain] = useState<Domain>();

  const location = useLocation();
  const { t } = useTranslation();

  const [tab, setTab] = useState<"base" | "advance">("base");

  const [targetType, setTargetType] = useState(domain ? domain.targetType : "");

  useEffect(() => {
    // Parsing query parameters
    const queryParams = new URLSearchParams(location.search);
    const id = queryParams.get("id");
    if (id) {
      const fetchData = async () => {
        const data = await get(id);
        setDomain(data);
        setTargetType(data.targetType);
      };
      fetchData();
    }
  }, [location.search]);

  const formSchema = z.object({
    id: z.string().optional(),
    domain: z.string().min(1, {
      message: "domain.not.empty.verify.message",
    }),
    email: z.string().email("email.valid.message").optional(),
    access: z.string().regex(/^[a-zA-Z0-9]+$/, {
      message: "domain.management.edit.dns.access.not.empty.message",
    }),
    targetAccess: z.string().optional(),
    targetType: z.string().regex(/^[a-zA-Z0-9-]+$/, {
      message: "domain.management.edit.target.type.not.empty.message",
    }),
    variables: z.string().optional(),
    group: z.string().optional(),
    nameservers: z.string().optional(),
  });

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      id: "",
      domain: "",
      email: "",
      access: "",
      targetAccess: "",
      targetType: "",
      variables: "",
      group: "",
      nameservers: "",
    },
  });

  useEffect(() => {
    if (domain) {
      form.reset({
        id: domain.id,
        domain: domain.domain,
        email: domain.email,
        access: domain.access,
        targetAccess: domain.targetAccess,
        targetType: domain.targetType,
        variables: domain.variables,
        group: domain.group,
        nameservers: domain.nameservers,
      });
    }
  }, [domain, form]);

  const targetAccesses = accesses.filter((item) => {
    if (item.usage == "apply") {
      return false;
    }

    if (targetType == "") {
      return true;
    }
    const types = targetType.split("-");
    return item.configType === types[0];
  });

  const { toast } = useToast();

  const navigate = useNavigate();

  const onSubmit = async (data: z.infer<typeof formSchema>) => {
    const group = data.group == "emptyId" ? "" : data.group;
    const targetAccess =
      data.targetAccess === "emptyId" ? "" : data.targetAccess;
    if (group == "" && targetAccess == "") {
      form.setError("group", {
        type: "manual",
        message: "domain.management.edit.target.access.verify.msg",
      });
      form.setError("targetAccess", {
        type: "manual",
        message: "domain.management.edit.target.access.verify.msg",
      });
      return;
    }

    const req: Domain = {
      id: data.id as string,
      crontab: "0 0 * * *",
      domain: data.domain,
      email: data.email,
      access: data.access,
      group: group,
      targetAccess: targetAccess,
      targetType: data.targetType,
      variables: data.variables,
      nameservers: data.nameservers,
    };

    try {
      await save(req);
      let description = t("domain.management.edit.succeed.tips");
      if (req.id == "") {
        description = t("domain.management.add.succeed.tips");
      }

      toast({
        title: t("succeed"),
        description,
      });
      navigate("/domains");
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
      <div className="">
        <Toaster />
        <div className=" h-5 text-muted-foreground">
          {domain?.id ? t("domain.edit") : t("domain.add")}
        </div>
        <div className="mt-5 flex w-full justify-center md:space-x-10 flex-col md:flex-row">
          <div className="w-full md:w-[200px] text-muted-foreground space-x-3 md:space-y-3 flex-row md:flex-col flex">
            <div
              className={cn(
                "cursor-pointer text-right",
                tab === "base" ? "text-primary" : ""
              )}
              onClick={() => {
                setTab("base");
              }}
            >
              {t("basic.setting")}
            </div>
            <div
              className={cn(
                "cursor-pointer text-right",
                tab === "advance" ? "text-primary" : ""
              )}
              onClick={() => {
                setTab("advance");
              }}
            >
              {t("advanced.setting")}
            </div>
          </div>

          <div className="w-full md:w-[35em] bg-gray-100 dark:bg-gray-900 p-5 rounded mt-3 md:mt-0">
            <Form {...form}>
              <form
                onSubmit={form.handleSubmit(onSubmit)}
                className="space-y-8 dark:text-stone-200"
              >
                <FormField
                  control={form.control}
                  name="domain"
                  render={({ field }) => (
                    <FormItem hidden={tab != "base"}>
                      <>
                        <StringList
                          value={field.value}
                          valueType="domain"
                          onValueChange={(domain: string) => {
                            form.setValue("domain", domain);
                          }}
                        />
                      </>

                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="email"
                  render={({ field }) => (
                    <FormItem hidden={tab != "base"}>
                      <FormLabel className="flex w-full justify-between">
                        <div>
                          {t("email") +
                            t("domain.management.edit.email.description")}
                        </div>
                        <EmailsEdit
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
                          onValueChange={(value) => {
                            form.setValue("email", value);
                          }}
                        >
                          <SelectTrigger>
                            <SelectValue
                              placeholder={t(
                                "domain.management.edit.email.not.empty.message"
                              )}
                            />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectGroup>
                              <SelectLabel>{t("email.list")}</SelectLabel>
                              {(emails.content as EmailsSetting).emails.map(
                                (item) => (
                                  <SelectItem key={item} value={item}>
                                    <div>{item}</div>
                                  </SelectItem>
                                )
                              )}
                            </SelectGroup>
                          </SelectContent>
                        </Select>
                      </FormControl>

                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="access"
                  render={({ field }) => (
                    <FormItem hidden={tab != "base"}>
                      <FormLabel className="flex w-full justify-between">
                        <div>
                          {t("domain.management.edit.dns.access.label")}
                        </div>
                        <AccessEdit
                          trigger={
                            <div className="font-normal text-primary hover:underline cursor-pointer flex items-center">
                              <Plus size={14} />
                              {t("add")}
                            </div>
                          }
                          op="add"
                        />
                      </FormLabel>
                      <FormControl>
                        <Select
                          {...field}
                          value={field.value}
                          onValueChange={(value) => {
                            form.setValue("access", value);
                          }}
                        >
                          <SelectTrigger>
                            <SelectValue
                              placeholder={t(
                                "domain.management.edit.access.not.empty.message"
                              )}
                            />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectGroup>
                              <SelectLabel>
                                {t("domain.management.edit.access.label")}
                              </SelectLabel>
                              {accesses
                                .filter((item) => item.usage != "deploy")
                                .map((item) => (
                                  <SelectItem key={item.id} value={item.id}>
                                    <div className="flex items-center space-x-2">
                                      <img
                                        className="w-6"
                                        src={
                                          accessTypeMap.get(
                                            item.configType
                                          )?.[1]
                                        }
                                      />
                                      <div>{item.name}</div>
                                    </div>
                                  </SelectItem>
                                ))}
                            </SelectGroup>
                          </SelectContent>
                        </Select>
                      </FormControl>

                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="targetType"
                  render={({ field }) => (
                    <FormItem hidden={tab != "base"}>
                      <FormLabel>
                        {t("domain.management.edit.target.type")}
                      </FormLabel>
                      <FormControl>
                        <Select
                          {...field}
                          onValueChange={(value) => {
                            setTargetType(value);
                            form.setValue("targetType", value);
                          }}
                        >
                          <SelectTrigger>
                            <SelectValue
                              placeholder={t(
                                "domain.management.edit.target.type.not.empty.message"
                              )}
                            />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectGroup>
                              <SelectLabel>
                                {t("domain.management.edit.target.type")}
                              </SelectLabel>
                              {targetTypeKeys.map((key) => (
                                <SelectItem key={key} value={key}>
                                  <div className="flex items-center space-x-2">
                                    <img
                                      className="w-6"
                                      src={targetTypeMap.get(key)?.[1]}
                                    />
                                    <div>
                                      {t(targetTypeMap.get(key)?.[0] || "")}
                                    </div>
                                  </div>
                                </SelectItem>
                              ))}
                            </SelectGroup>
                          </SelectContent>
                        </Select>
                      </FormControl>

                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="targetAccess"
                  render={({ field }) => (
                    <FormItem hidden={tab != "base"}>
                      <FormLabel className="w-full flex justify-between">
                        <div>{t("domain.management.edit.target.access")}</div>
                        <AccessEdit
                          trigger={
                            <div className="font-normal text-primary hover:underline cursor-pointer flex items-center">
                              <Plus size={14} />
                              {t("add")}
                            </div>
                          }
                          op="add"
                        />
                      </FormLabel>
                      <FormControl>
                        <Select
                          {...field}
                          onValueChange={(value) => {
                            form.setValue("targetAccess", value);
                          }}
                        >
                          <SelectTrigger>
                            <SelectValue
                              placeholder={t(
                                "domain.management.edit.target.access.not.empty.message"
                              )}
                            />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectGroup>
                              <SelectLabel>
                                {t(
                                  "domain.management.edit.target.access.content.label"
                                )}{" "}
                                {form.getValues().targetAccess}
                              </SelectLabel>
                              <SelectItem value="emptyId">
                                <div className="flex items-center space-x-2">
                                  --
                                </div>
                              </SelectItem>
                              {targetAccesses.map((item) => (
                                <SelectItem key={item.id} value={item.id}>
                                  <div className="flex items-center space-x-2">
                                    <img
                                      className="w-6"
                                      src={
                                        accessTypeMap.get(item.configType)?.[1]
                                      }
                                    />
                                    <div>{item.name}</div>
                                  </div>
                                </SelectItem>
                              ))}
                            </SelectGroup>
                          </SelectContent>
                        </Select>
                      </FormControl>

                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="group"
                  render={({ field }) => (
                    <FormItem hidden={tab != "advance" || targetType != "ssh"}>
                      <FormLabel className="w-full flex justify-between">
                        <div>{t("domain.management.edit.group.label")}</div>
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
                              placeholder={t(
                                "domain.management.edit.group.not.empty.message"
                              )}
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
                            {accessGroups
                              .filter((item) => {
                                return (
                                  item.expand && item.expand?.access.length > 0
                                );
                              })
                              .map((item) => (
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
                  name="variables"
                  render={({ field }) => (
                    <FormItem hidden={tab != "advance"}>
                      <FormLabel>{t("variables")}</FormLabel>
                      <FormControl>
                        <Textarea
                          placeholder={t(
                            "domain.management.edit.variables.placeholder"
                          )}
                          {...field}
                          className="placeholder:whitespace-pre-wrap"
                        />
                      </FormControl>

                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="nameservers"
                  render={({ field }) => (
                    <FormItem hidden={tab != "advance"}>
                      <FormLabel>{t("dns")}</FormLabel>
                      <FormControl>
                        <Textarea
                          placeholder={t(
                            "domain.management.edit.dns.placeholder"
                          )}
                          {...field}
                          className="placeholder:whitespace-pre-wrap"
                        />
                      </FormControl>

                      <FormMessage />
                    </FormItem>
                  )}
                />

                <div className="flex justify-end">
                  <Button type="submit">{t("save")}</Button>
                </div>
              </form>
            </Form>
          </div>
        </div>
      </div>
    </>
  );
};

export default Edit;
