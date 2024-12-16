import { memo } from "react";
import { useTranslation } from "react-i18next";
import { Avatar, Select, Space, Tag, Typography, type SelectProps } from "antd";

import { accessProvidersMap } from "@/domain/access";

export type AccessTypeSelectProps = Omit<SelectProps, "labelRender" | "options" | "optionFilterProp" | "optionLabelProp" | "optionRender">;

const AccessTypeSelect = memo((props: AccessTypeSelectProps) => {
  const { t } = useTranslation();

  const options = Array.from(accessProvidersMap.values()).map((item) => ({
    key: item.type,
    value: item.type,
    label: t(item.name),
  }));

  return (
    <Select
      {...props}
      labelRender={({ label, value }) => {
        if (label) {
          return (
            <Space className="max-w-full truncate" size={4}>
              <Avatar src={accessProvidersMap.get(String(value))?.icon} size="small" />
              {label}
            </Space>
          );
        }

        return <Typography.Text type="secondary">{props.placeholder}</Typography.Text>;
      }}
      options={options}
      optionFilterProp={undefined}
      optionLabelProp={undefined}
      optionRender={(option) => (
        <div className="flex items-center justify-between gap-4 max-w-full overflow-hidden">
          <Space className="flex-grow max-w-full truncate" size={4}>
            <Avatar src={accessProvidersMap.get(option.data.value)?.icon} size="small" />
            <Typography.Text ellipsis>{t(accessProvidersMap.get(option.data.value)?.name ?? "")}</Typography.Text>
          </Space>
          <div>
            {accessProvidersMap.get(option.data.value)?.usage === "apply" && (
              <>
                <Tag color="orange">{t("access.props.provider.usage.dns")}</Tag>
              </>
            )}
            {accessProvidersMap.get(option.data.value)?.usage === "deploy" && (
              <>
                <Tag color="blue">{t("access.props.provider.usage.host")}</Tag>
              </>
            )}
            {accessProvidersMap.get(option.data.value)?.usage === "all" && (
              <>
                <Tag color="orange">{t("access.props.provider.usage.dns")}</Tag>
                <Tag color="blue">{t("access.props.provider.usage.host")}</Tag>
              </>
            )}
          </div>
        </div>
      )}
    />
  );
});

export default AccessTypeSelect;
