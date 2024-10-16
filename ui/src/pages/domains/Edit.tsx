import { useEffect, useState } from "react";
import { useLocation } from "react-router-dom";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import z from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { Plus } from "lucide-react";
import { ClientResponseError } from "pocketbase";

import { Button } from "@/components/ui/button";
import { Breadcrumb, BreadcrumbItem, BreadcrumbLink, BreadcrumbList, BreadcrumbPage, BreadcrumbSeparator } from "@/components/ui/breadcrumb";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectGroup, SelectItem, SelectLabel, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Toaster } from "@/components/ui/toaster";
import { useToast } from "@/components/ui/use-toast";
import AccessEdit from "@/components/certimate/AccessEdit";
import DeployList from "@/components/certimate/DeployList";
import EmailsEdit from "@/components/certimate/EmailsEdit";
import StringList from "@/components/certimate/StringList";
import { cn } from "@/lib/utils";
import { PbErrorData } from "@/domain/base";
import { accessTypeMap } from "@/domain/access";
import { EmailsSetting } from "@/domain/settings";
import { DeployConfig, Domain } from "@/domain/domain";
import { save, get } from "@/repository/domains";
import { useConfig } from "@/providers/config";

const Edit = () => {
  const {
    config: { accesses, emails },
  } = useConfig();

  const [domain, setDomain] = useState<Domain>({} as Domain);

  const location = useLocation();
  const { t } = useTranslation();

  const [tab, setTab] = useState<"apply" | "deploy">("apply");

  useEffect(() => {
    // Parsing query parameters
    const queryParams = new URLSearchParams(location.search);
    const id = queryParams.get("id");
    if (id) {
      const fetchData = async () => {
        const data = await get(id);
        setDomain(data);
      };
      fetchData();
    }
  }, [location.search]);

  const formSchema = z.object({
    id: z.string().optional(),
    domain: z.string().min(1, {
      message: "common.errmsg.domain_invalid",
    }),
    email: z.string().email("common.errmsg.email_invalid").optional(),
    access: z.string().regex(/^[a-zA-Z0-9]+$/, {
      message: "domain.application.form.access.errmsg.empty",
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

  const { toast } = useToast();

  const onSubmit = async (data: z.infer<typeof formSchema>) => {
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
      const resp = await save(req);
      let description = t("domain.application.form.domain.changed.message");
      if (req.id == "") {
        description = t("domain.application.form.domain.added.message");
      }

      toast({
        title: t("common.save.succeeded.message"),
        description,
      });

      if (!domain?.id) setTab("deploy");
      setDomain({ ...resp });
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

  const handelOnDeployListChange = async (list: DeployConfig[]) => {
    const req = {
      ...domain,
      deployConfig: list,
    };
    try {
      const resp = await save(req);
      let description = t("domain.application.form.domain.changed.message");
      if (req.id == "") {
        description = t("domain.application.form.domain.added.message");
      }

      toast({
        title: t("common.save.succeeded.message"),
        description,
      });

      if (!domain?.id) setTab("deploy");
      setDomain({ ...resp });
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

  return (
    <>
      <div className="">
        <Toaster />
        <div className=" h-5 text-muted-foreground">
          <Breadcrumb>
            <BreadcrumbList>
              <BreadcrumbItem>
                <BreadcrumbLink href="#/domains">{t("domain.page.title")}</BreadcrumbLink>
              </BreadcrumbItem>
              <BreadcrumbSeparator />

              <BreadcrumbItem>
                <BreadcrumbPage>{domain?.id ? t("domain.edit") : t("domain.add")}</BreadcrumbPage>
              </BreadcrumbItem>
            </BreadcrumbList>
          </Breadcrumb>
        </div>
        <div className="mt-5 flex w-full justify-center md:space-x-10 flex-col md:flex-row">
          <div className="w-full md:w-[200px] text-muted-foreground space-x-3 md:space-y-3 flex-row md:flex-col flex md:mt-5">
            <div
              className={cn("cursor-pointer text-right", tab === "apply" ? "text-primary" : "")}
              onClick={() => {
                setTab("apply");
              }}
            >
              {t("domain.application.tab")}
            </div>
            <div
              className={cn("cursor-pointer text-right", tab === "deploy" ? "text-primary" : "")}
              onClick={() => {
                if (!domain?.id) {
                  toast({
                    title: t("domain.application.unsaved.message"),
                    description: t("domain.application.unsaved.message"),
                    variant: "destructive",
                  });
                  return;
                }
                setTab("deploy");
              }}
            >
              {t("domain.deployment.tab")}
            </div>
          </div>
          <div className="flex flex-col">
            <div className={cn("w-full md:w-[35em]  p-5 rounded mt-3 md:mt-0", tab == "deploy" && "hidden")}>
              <Form {...form}>
                <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8 dark:text-stone-200">
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
                          <div>{t("domain.application.form.email.label") + " " + t("domain.application.form.email.tips")}</div>
                          <EmailsEdit
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
                            onValueChange={(value) => {
                              form.setValue("email", value);
                            }}
                          >
                            <SelectTrigger>
                              <SelectValue placeholder={t("domain.application.form.email.errmsg.empty")} />
                            </SelectTrigger>
                            <SelectContent>
                              <SelectGroup>
                                <SelectLabel>{t("domain.application.form.email.list")}</SelectLabel>
                                {(emails.content as EmailsSetting).emails.map((item) => (
                                  <SelectItem key={item} value={item}>
                                    <div>{item}</div>
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
                  {/* 授权 */}
                  <FormField
                    control={form.control}
                    name="access"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel className="flex w-full justify-between">
                          <div>{t("domain.application.form.access.label")}</div>
                          <AccessEdit
                            trigger={
                              <div className="font-normal text-primary hover:underline cursor-pointer flex items-center">
                                <Plus size={14} />
                                {t("common.add")}
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
                              <SelectValue placeholder={t("domain.application.form.access.placeholder")} />
                            </SelectTrigger>
                            <SelectContent>
                              <SelectGroup>
                                <SelectLabel>{t("domain.application.form.access.list")}</SelectLabel>
                                {accesses
                                  .filter((item) => item.usage != "deploy")
                                  .map((item) => (
                                    <SelectItem key={item.id} value={item.id}>
                                      <div className="flex items-center space-x-2">
                                        <img className="w-6" src={accessTypeMap.get(item.configType)?.[1]} />
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
                        <FormLabel>{t("domain.application.form.timeout.label")}</FormLabel>
                        <FormControl>
                          <Input
                            type="number"
                            placeholder={t("ddomain.application.form.timeout.placeholder")}
                            {...field}
                            value={field.value}
                            onChange={(e) => {
                              form.setValue("timeout", parseInt(e.target.value));
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
                    <Button type="submit">{domain?.id ? t("common.save") : t("common.next")}</Button>
                  </div>
                </form>
              </Form>
            </div>

            <div className={cn("flex flex-col space-y-5 w-full md:w-[35em]", tab == "apply" && "hidden")}>
              <DeployList
                deploys={domain?.deployConfig ?? []}
                onChange={(list: DeployConfig[]) => {
                  handelOnDeployListChange(list);
                }}
              />
            </div>
          </div>
        </div>
      </div>
    </>
  );
};

export default Edit;
