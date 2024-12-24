import { memo, useEffect } from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { Collapse, Switch, Tooltip } from "antd";
import z from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { ChevronsUpDown as ChevronsUpDownIcon, Plus as PlusIcon, CircleHelp as CircleHelpIcon } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectGroup, SelectItem, SelectLabel, SelectTrigger, SelectValue } from "@/components/ui/select";
import AccessEditModal from "@/components/access/AccessEditModal";
import EmailsEdit from "@/components/certimate/EmailsEdit";
import StringList from "@/components/certimate/StringList";
import { accessProvidersMap } from "@/domain/access";
import { useZustandShallowSelector } from "@/hooks";
import { useAccessStore } from "@/stores/access";
import { useContactStore } from "@/stores/contact";
import { WorkflowNode, WorkflowNodeConfig } from "@/domain/workflow";
import { useWorkflowStore } from "@/stores/workflow";
import { usePanel } from "./PanelProvider";

type ApplyFormProps = {
  data: WorkflowNode;
};

const ApplyForm = ({ data }: ApplyFormProps) => {
  const { updateNode } = useWorkflowStore(useZustandShallowSelector(["updateNode"]));

  const { accesses } = useAccessStore();
  const { emails, fetchEmails } = useContactStore();

  useEffect(() => {
    fetchEmails();
  }, []);

  const { t } = useTranslation();

  const { hidePanel } = usePanel();

  const formSchema = z.object({
    domain: z.string().min(1, {
      message: "common.errmsg.domain_invalid",
    }),
    email: z.string().email("common.errmsg.email_invalid").optional(),
    access: z.string().regex(/^[a-zA-Z0-9]+$/, {
      message: "domain.application.form.access.placeholder",
    }),
    keyAlgorithm: z.string().optional(),
    nameservers: z.string().optional(),
    timeout: z.number().optional(),
    disableFollowCNAME: z.boolean().optional(),
  });

  let config: WorkflowNodeConfig = {
    domain: "",
    email: "",
    access: "",
    keyAlgorithm: "RSA2048",
    nameservers: "",
    timeout: 60,
    disableFollowCNAME: true,
  };
  if (data) config = data.config ?? config;

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      domain: config.domain as string,
      email: config.email as string,
      access: config.access as string,
      keyAlgorithm: config.keyAlgorithm as string,
      nameservers: config.nameservers as string,
      timeout: config.timeout as number,
      disableFollowCNAME: config.disableFollowCNAME as boolean,
    },
  });

  const onSubmit = async (config: z.infer<typeof formSchema>) => {
    updateNode({ ...data, config, validated: true });
    hidePanel();
  };

  return (
    <>
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
                <FormLabel className="flex justify-between w-full">
                  <div>{t("domain.application.form.email.label") + " " + t("domain.application.form.email.tips")}</div>
                  <EmailsEdit
                    trigger={
                      <div className="flex items-center font-normal cursor-pointer text-primary hover:underline">
                        <PlusIcon size={14} />
                        {t("common.button.add")}
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
                      <SelectValue placeholder={t("domain.application.form.email.placeholder")} />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectGroup>
                        <SelectLabel>{t("domain.application.form.email.list")}</SelectLabel>
                        {emails.map((item) => (
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

          {/* DNS 服务商授权 */}
          <FormField
            control={form.control}
            name="access"
            render={({ field }) => (
              <FormItem>
                <FormLabel className="flex justify-between w-full">
                  <div>{t("domain.application.form.access.label")}</div>
                  <AccessEditModal
                    mode="add"
                    trigger={
                      <div className="flex items-center font-normal cursor-pointer text-primary hover:underline">
                        <PlusIcon size={14} />
                        {t("common.button.add")}
                      </div>
                    }
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
                                <img className="w-6" src={accessProvidersMap.get(item.configType)?.icon} />
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

          <div>
            <hr />
            <Collapse
              bordered={false}
              ghost={true}
              items={[
                {
                  key: "advanced",
                  styles: {
                    header: { paddingLeft: 0, paddingRight: 0 },
                    body: { paddingLeft: 0, paddingRight: 0 },
                  },
                  label: <>{t("domain.application.form.advanced_settings.label")}</>,
                  children: (
                    <div className="flex flex-col space-y-8">
                      {/* 证书算法 */}
                      <FormField
                        control={form.control}
                        name="keyAlgorithm"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>{t("domain.application.form.key_algorithm.label")}</FormLabel>
                            <Select
                              {...field}
                              value={field.value}
                              onValueChange={(value) => {
                                form.setValue("keyAlgorithm", value);
                              }}
                            >
                              <SelectTrigger>
                                <SelectValue placeholder={t("domain.application.form.key_algorithm.placeholder")} />
                              </SelectTrigger>
                              <SelectContent>
                                <SelectGroup>
                                  <SelectItem value="RSA2048">RSA2048</SelectItem>
                                  <SelectItem value="RSA3072">RSA3072</SelectItem>
                                  <SelectItem value="RSA4096">RSA4096</SelectItem>
                                  <SelectItem value="RSA8192">RSA8192</SelectItem>
                                  <SelectItem value="EC256">EC256</SelectItem>
                                  <SelectItem value="EC384">EC384</SelectItem>
                                </SelectGroup>
                              </SelectContent>
                            </Select>
                          </FormItem>
                        )}
                      />

                      {/* DNS */}
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

                      {/* DNS 超时时间 */}
                      <FormField
                        control={form.control}
                        name="timeout"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>{t("domain.application.form.timeout.label")}</FormLabel>
                            <FormControl>
                              <Input
                                type="number"
                                placeholder={t("domain.application.form.timeout.placeholder")}
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

                      {/* 禁用 CNAME 跟随 */}
                      <FormField
                        control={form.control}
                        name="disableFollowCNAME"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>
                              <div className="flex">
                                <span className="mr-1">{t("domain.application.form.disable_follow_cname.label")} </span>
                                <Tooltip
                                  title={
                                    <p>
                                      {t("domain.application.form.disable_follow_cname.tips")}
                                      <a
                                        className="text-primary"
                                        target="_blank"
                                        href="https://letsencrypt.org/2019/10/09/onboarding-your-customers-with-lets-encrypt-and-acme/#the-advantages-of-a-cname"
                                      >
                                        {t("domain.application.form.disable_follow_cname.tips_link")}
                                      </a>
                                    </p>
                                  }
                                >
                                  <CircleHelpIcon size={14} />
                                </Tooltip>
                              </div>
                            </FormLabel>
                            <FormControl>
                              <div>
                                <Switch
                                  defaultChecked={field.value}
                                  onChange={(value) => {
                                    form.setValue(field.name, value);
                                  }}
                                />
                              </div>
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />
                    </div>
                  ),
                  extra: <ChevronsUpDownIcon size={14} />,
                  forceRender: true,
                  showArrow: false,
                },
              ]}
            />
          </div>

          <div className="flex justify-end">
            <Button type="submit">{t("common.button.save")}</Button>
          </div>
        </form>
      </Form>
    </>
  );
};

export default memo(ApplyForm);
