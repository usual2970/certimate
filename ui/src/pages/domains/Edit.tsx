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
import { Plus, Trash2, Edit as EditIcon } from "lucide-react";
import { AccessEdit } from "@/components/certimate/AccessEdit";
import { accessTypeMap } from "@/domain/access";
import EmailsEdit from "@/components/certimate/EmailsEdit";
import { Textarea } from "@/components/ui/textarea";
import { cn } from "@/lib/utils";
import { EmailsSetting } from "@/domain/settings";
import { useTranslation } from "react-i18next";
import StringList from "@/components/certimate/StringList";
import { Input } from "@/components/ui/input";

const Edit = () => {
  const {
    config: { accesses, emails, accessGroups },
  } = useConfig();

  const [domain, setDomain] = useState<Domain>();

  const location = useLocation();
  const { t } = useTranslation();

  const [tab, setTab] = useState<"apply" | "deploy">("apply");

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
    nameservers: z.string().optional(),
    timeout: z.number().optional(),
  });

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      id: "",
      domain: "",
      email: "",
      access: "",
      nameservers: "",
      timeout: 60,
    },
  });

  useEffect(() => {
    if (domain) {
      form.reset({
        id: domain.id,
        domain: domain.domain,
        email: domain.applyConfig?.email,
        access: domain.applyConfig?.access,

        nameservers: domain.applyConfig?.nameservers,
        timeout: domain.applyConfig?.timeout,
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
    console.log(data);
    const req: Domain = {
      id: data.id as string,
      crontab: "0 0 * * *",
      domain: data.domain,
      email: data.email,
      access: data.access,
      applyConfig: {
        email: data.email ?? "",
        access: data.access,
        nameservers: data.nameservers,
        timeout: data.timeout,
      },
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
      if (!domain?.id) setTab("deploy");
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
          <div className="w-full md:w-[200px] text-muted-foreground space-x-3 md:space-y-3 flex-row md:flex-col flex md:mt-5">
            <div
              className={cn(
                "cursor-pointer text-right",
                tab === "apply" ? "text-primary" : ""
              )}
              onClick={() => {
                setTab("apply");
              }}
            >
              {t("apply.setting")}
            </div>
            <div
              className={cn(
                "cursor-pointer text-right",
                tab === "deploy" ? "text-primary" : ""
              )}
              onClick={() => {
                if (!domain?.id) {
                  toast({
                    title: t("domain.management.edit.deploy.error"),
                    description: t("domain.management.edit.deploy.error"),
                    variant: "destructive",
                  });
                  return;
                }
                setTab("deploy");
              }}
            >
              {t("deploy.setting")}
            </div>
          </div>
          <div className="flex flex-col">
            <div
              className={cn(
                "w-full md:w-[35em]  p-5 rounded mt-3 md:mt-0",
                tab == "deploy" && "hidden"
              )}
            >
              <Form {...form}>
                <form
                  onSubmit={form.handleSubmit(onSubmit)}
                  className="space-y-8 dark:text-stone-200"
                >
                  {/* 域名 */}
                  <FormField
                    control={form.control}
                    name="domain"
                    render={({ field }) => (
                      <FormItem>
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
                  {/* 邮箱 */}
                  <FormField
                    control={form.control}
                    name="email"
                    render={({ field }) => (
                      <FormItem>
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
                  {/* 授权 */}
                  <FormField
                    control={form.control}
                    name="access"
                    render={({ field }) => (
                      <FormItem>
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

                  {/* 超时时间 */}
                  <FormField
                    control={form.control}
                    name="timeout"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{t("timeout")}</FormLabel>
                        <FormControl>
                          <Input
                            type="number"
                            placeholder={t(
                              "domain.management.edit.timeout.placeholder"
                            )}
                            {...field}
                            value={field.value}
                            onChange={(e) => {
                              form.setValue(
                                "timeout",
                                parseInt(e.target.value)
                              );
                            }}
                          />
                        </FormControl>

                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  {/* nameservers */}
                  <FormField
                    control={form.control}
                    name="nameservers"
                    render={({ field }) => (
                      <FormItem>
                        <StringList
                          value={field.value ?? ""}
                          onValueChange={(val: string) => {
                            form.setValue("nameservers", val);
                          }}
                          valueType="dns"
                        ></StringList>

                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <div className="flex justify-end">
                    <Button type="submit">
                      {domain?.id ? t("save") : t("next")}
                    </Button>
                  </div>
                </form>
              </Form>
            </div>

            <div
              className={cn(
                "flex flex-col space-y-5",
                tab == "apply" && "hidden"
              )}
            >
              <div className="flex justify-end py-2 border-b">
                <Button size={"sm"} variant={"secondary"}>
                  添加部署
                </Button>
              </div>
              <div className="w-full md:w-[35em] rounded mt-5 border">
                <div className="">
                  <div className="flex justify-between text-sm p-3 items-center text-stone-700">
                    <div>ALIYUN-CDN</div>
                    <div className="flex space-x-2">
                      <EditIcon size={16} />
                      <Trash2 size={16} />
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  );
};

export default Edit;
