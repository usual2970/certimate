import { forwardRef, memo, useEffect, useImperativeHandle } from "react";
import { useTranslation } from "react-i18next";
import { Button, Form, type FormInstance, Input, Upload, type UploadProps } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { validateCertificate, validatePrivateKey } from "@/api/certificates";
import { type WorkflowNodeConfigForUpload } from "@/domain/workflow";
import { useAntdForm } from "@/hooks";
import { getErrMsg } from "@/utils/error";
import { readFileContent } from "@/utils/file";

type UploadNodeConfigFormFieldValues = Partial<WorkflowNodeConfigForUpload>;

export type UploadNodeConfigFormProps = {
  className?: string;
  style?: React.CSSProperties;
  disabled?: boolean;
  initialValues?: UploadNodeConfigFormFieldValues;
  onValuesChange?: (values: UploadNodeConfigFormFieldValues) => void;
};

export type UploadNodeConfigFormInstance = {
  getFieldsValue: () => ReturnType<FormInstance<UploadNodeConfigFormFieldValues>["getFieldsValue"]>;
  resetFields: FormInstance<UploadNodeConfigFormFieldValues>["resetFields"];
  validateFields: FormInstance<UploadNodeConfigFormFieldValues>["validateFields"];
};

const initFormModel = (): UploadNodeConfigFormFieldValues => {
  return {};
};

const UploadNodeConfigForm = forwardRef<UploadNodeConfigFormInstance, UploadNodeConfigFormProps>(
  ({ className, style, disabled, initialValues, onValuesChange }, ref) => {
    const { t } = useTranslation();

    const formSchema = z.object({
      certificateId: z.string().optional(),
      domains: z.string().optional(),
      certificate: z
        .string({ message: t("workflow_node.upload.form.certificate.placeholder") })
        .min(1, t("workflow_node.upload.form.certificate.placeholder"))
        .max(5120, t("common.errmsg.string_max", { max: 5120 })),
      privateKey: z
        .string({ message: t("workflow_node.upload.form.private_key.placeholder") })
        .min(1, t("workflow_node.upload.form.private_key.placeholder"))
        .max(5120, t("common.errmsg.string_max", { max: 5120 })),
    });
    const formRule = createSchemaFieldRule(formSchema);
    const { form: formInst, formProps } = useAntdForm({
      name: "workflowNodeUploadConfigForm",
      initialValues: initialValues ?? initFormModel(),
    });

    const certificate = Form.useWatch("certificate", formInst);
    const privateKey = Form.useWatch("privateKey", formInst);

    useEffect(() => {
      if (certificate && privateKey) {
        formInst.validateFields(["certificate", "privateKey"]);
      }
    }, [certificate, privateKey]);

    const handleFormChange = (_: unknown, values: z.infer<typeof formSchema>) => {
      onValuesChange?.(values as UploadNodeConfigFormFieldValues);
    };

    useImperativeHandle(ref, () => {
      return {
        getFieldsValue: () => {
          return formInst.getFieldsValue(true);
        },
        resetFields: (fields) => {
          return formInst.resetFields(fields as (keyof UploadNodeConfigFormFieldValues)[]);
        },
        validateFields: (nameList, config) => {
          return formInst.validateFields(nameList, config);
        },
      } as UploadNodeConfigFormInstance;
    });

    const handleCertificateFileChange: UploadProps["onChange"] = async ({ file }) => {
      if (file && file.status !== "removed") {
        const certificate = await readFileContent(file.originFileObj ?? (file as unknown as File));

        try {
          const resp = await validateCertificate(certificate);
          formInst.setFields([
            {
              name: "certificate",
              value: certificate,
              errors: [],
            },
            {
              name: "domains",
              value: resp.data.domains,
            },
          ]);
        } catch (e) {
          formInst.setFields([
            {
              name: "certificate",
              value: "",
              errors: [getErrMsg(e)],
            },
          ]);
        }
      } else {
        formInst.setFieldValue("certificate", "");
      }
      onValuesChange?.(formInst.getFieldsValue(true));
    };

    const handlePrivateKeyFileChange: UploadProps["onChange"] = async ({ file }) => {
      if (file && file.status !== "removed") {
        const privateKey = await readFileContent(file.originFileObj ?? (file as unknown as File));
        try {
          await validatePrivateKey(privateKey);
          formInst.setFields([
            {
              name: "privateKey",
              value: privateKey,
              errors: [],
            },
          ]);
        } catch (e) {
          formInst.setFields([
            {
              name: "privateKey",
              errors: [getErrMsg(e)],
            },
          ]);
        }
      } else {
        formInst.setFieldValue("privateKey", "");
      }

      onValuesChange?.(formInst.getFieldsValue(true));
    };

    return (
      <Form className={className} style={style} {...formProps} disabled={disabled} layout="vertical" scrollToFirstError onValuesChange={handleFormChange}>
        <Form.Item name="domains" label={t("workflow_node.upload.form.domains.label")} rules={[formRule]}>
          <Input readOnly />
        </Form.Item>

        <Form.Item name="certificate" label={t("workflow_node.upload.form.certificate.label")} rules={[formRule]}>
          <Input.TextArea
            readOnly
            autoSize={{ minRows: 5, maxRows: 10 }}
            placeholder={t("workflow_node.upload.form.certificate.placeholder")}
            value={certificate}
          />
          <div className="mt-2 text-right">
            <Upload beforeUpload={() => false} maxCount={1} onChange={handleCertificateFileChange}>
              <Button>{t("workflow_node.upload.form.certificate.button")}</Button>
            </Upload>
          </div>
        </Form.Item>

        <Form.Item name="privateKey" label={t("workflow_node.upload.form.private_key.label")} rules={[formRule]}>
          <Input.TextArea
            readOnly
            autoSize={{ minRows: 5, maxRows: 10 }}
            placeholder={t("workflow_node.upload.form.private_key.placeholder")}
            value={privateKey}
          />
          <div className="mt-2 text-right">
            <Upload beforeUpload={() => false} maxCount={1} onChange={handlePrivateKeyFileChange}>
              <Button>{t("workflow_node.upload.form.private_key.button")}</Button>
            </Upload>
          </div>
        </Form.Item>
      </Form>
    );
  }
);

export default memo(UploadNodeConfigForm);
