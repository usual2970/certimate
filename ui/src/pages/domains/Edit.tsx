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

const Edit = () => {
  const {
    config: { accesses, emails },
  } = useConfig();

  const [domain, setDomain] = useState<Domain>();

  const location = useLocation();

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
    domain: z.string().regex(/^(?:\*\.)?([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$/, {
      message: "请输入正确的域名",
    }),
    email: z.string().email().optional(),
    access: z.string().regex(/^[a-zA-Z0-9]+$/, {
      message: "请选择DNS服务商授权配置",
    }),
    targetAccess: z.string().regex(/^[a-zA-Z0-9]+$/, {
      message: "请选择部署服务商配置",
    }),
    targetType: z.string().regex(/^[a-zA-Z0-9-]+$/, {
      message: "请选择部署服务类型",
    }),
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
      });
    }
  }, [domain, form]);

  const [targetType, setTargetType] = useState(domain ? domain.targetType : "");

  const targetAccesses = accesses.filter((item) => {
    if (item.usage == "apply") {
      return false;
    }

    if (targetType == "") {
      return true;
    }
    const types = form.getValues().targetType.split("-");
    return item.configType === types[0];
  });

  const { toast } = useToast();

  const navigate = useNavigate();

  const onSubmit = async (data: z.infer<typeof formSchema>) => {
    const req: Domain = {
      id: data.id as string,
      crontab: "0 0 * * *",
      domain: data.domain,
      email: data.email,
      access: data.access,
      targetAccess: data.targetAccess,
      targetType: data.targetType,
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
        <div className="border-b dark:border-stone-500 h-10 text-muted-foreground">
          {domain?.id ? "编辑" : "新增"}域名
        </div>
        <div className="max-w-[35em] mx-auto mt-10">
          <Form {...form}>
            <form
              onSubmit={form.handleSubmit(onSubmit)}
              className="space-y-8 dark:text-stone-200"
            >
              <FormField
                control={form.control}
                name="domain"
                render={({ field }) => (
                  <FormItem>
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
                  <FormItem>
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
                            {emails.content.emails.map((item) => (
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

              <FormField
                control={form.control}
                name="access"
                render={({ field }) => (
                  <FormItem>
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
                name="targetType"
                render={({ field }) => (
                  <FormItem>
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
                  <FormItem>
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
                            <SelectLabel>服务商授权配置</SelectLabel>
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

              <div className="flex justify-end">
                <Button type="submit">保存</Button>
              </div>
            </form>
          </Form>
        </div>
      </div>
    </>
  );
};

export default Edit;
