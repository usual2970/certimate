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

const Edit = () => {
  const {
    config: { accesses, emails, accessGroups },
  } = useConfig();

  const [domain, setDomain] = useState<Domain>();

  const location = useLocation();

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
    domain: z.string().regex(/^(?:\*\.)?([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$/, {
      message: "请输入正确的域名",
    }),
    email: z.string().email().optional(),
    access: z.string().regex(/^[a-zA-Z0-9]+$/, {
      message: "请选择DNS服务商授权配置",
    }),
    targetAccess: z.string().optional(),
    targetType: z.string().regex(/^[a-zA-Z0-9-]+$/, {
      message: "请选择部署服务类型",
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
        message: "部署授权和部署授权组至少选一个",
      });
      form.setError("targetAccess", {
        type: "manual",
        message: "部署授权和部署授权组至少选一个",
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
      let description = "域名编辑成功";
      if (req.id == "") {
        description = "域名添加成功";
      }

      toast({
        title: "成功",
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
          {domain?.id ? "编辑" : "新增"}域名
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
              基础设置
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
              高级设置
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
                      <FormLabel>域名</FormLabel>
                      <FormControl>
                        <Input placeholder="请输入域名" {...field} />
                      </FormControl>

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
                        <div>Email（申请证书需要提供邮箱）</div>
                        <EmailsEdit
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
                          onValueChange={(value) => {
                            form.setValue("email", value);
                          }}
                        >
                          <SelectTrigger>
                            <SelectValue placeholder="请选择邮箱" />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectGroup>
                              <SelectLabel>邮箱列表</SelectLabel>
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
                        <div>DNS 服务商授权配置</div>
                        <AccessEdit
                          trigger={
                            <div className="font-normal text-primary hover:underline cursor-pointer flex items-center">
                              <Plus size={14} />
                              新增
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
                            <SelectValue placeholder="请选择授权配置" />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectGroup>
                              <SelectLabel>服务商授权配置</SelectLabel>
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
                      <FormLabel>部署服务类型</FormLabel>
                      <FormControl>
                        <Select
                          {...field}
                          onValueChange={(value) => {
                            setTargetType(value);
                            form.setValue("targetType", value);
                          }}
                        >
                          <SelectTrigger>
                            <SelectValue placeholder="请选择部署服务类型" />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectGroup>
                              <SelectLabel>部署服务类型</SelectLabel>
                              {targetTypeKeys.map((key) => (
                                <SelectItem key={key} value={key}>
                                  <div className="flex items-center space-x-2">
                                    <img
                                      className="w-6"
                                      src={targetTypeMap.get(key)?.[1]}
                                    />
                                    <div>{targetTypeMap.get(key)?.[0]}</div>
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
                        <div>部署服务商授权配置</div>
                        <AccessEdit
                          trigger={
                            <div className="font-normal text-primary hover:underline cursor-pointer flex items-center">
                              <Plus size={14} />
                              新增
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
                            <SelectValue placeholder="请选择授权配置" />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectGroup>
                              <SelectLabel>
                                服务商授权配置{form.getValues().targetAccess}
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
                        <div>
                          部署配置组(用于将一个域名证书部署到多个 ssh 主机)
                        </div>
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
                      <FormLabel>变量</FormLabel>
                      <FormControl>
                        <Textarea
                          placeholder={`可在SSH部署中使用,形如:\nkey=val;\nkey2=val2;`}
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
                      <FormLabel>域名服务器</FormLabel>
                      <FormControl>
                        <Textarea
                          placeholder={`自定义域名服务器,多个用分号隔开,如:\n8.8.8.8;\n8.8.4.4;`}
                          {...field}
                          className="placeholder:whitespace-pre-wrap"
                        />
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
        </div>
      </div>
    </>
  );
};

export default Edit;
