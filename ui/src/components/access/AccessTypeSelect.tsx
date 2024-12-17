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

  const renderOption = (key: string) => {
    const provider = accessProvidersMap.get(key);
    return (
      <div className="flex items-center justify-between gap-4 max-w-full overflow-hidden">
        <Space className="flex-grow max-w-full truncate" size={4}>
          <Avatar src={provider?.icon} size="small" />
          <Typography.Text className="leading-loose" ellipsis>
            {t(provider?.name ?? "")}
          </Typography.Text>
        </Space>
        <div>
          {provider?.usage === "apply" && (
            <>
              <Tag color="orange">{t("access.props.provider.usage.dns")}</Tag>
            </>
          )}
          {provider?.usage === "deploy" && (
            <>
              <Tag color="blue">{t("access.props.provider.usage.host")}</Tag>
            </>
          )}
          {provider?.usage === "all" && (
            <>
              <Tag color="orange">{t("access.props.provider.usage.dns")}</Tag>
              <Tag color="blue">{t("access.props.provider.usage.host")}</Tag>
            </>
          )}
        </div>
      </div>
    );
  };

  return (
    <Select
      {...props}
      labelRender={({ label, value }) => {
        if (label) {
          return renderOption(value as string);
        }

        return <Typography.Text type="secondary">{props.placeholder}</Typography.Text>;
      }}
      options={options}
      optionFilterProp={undefined}
      optionLabelProp={undefined}
      optionRender={(option) => renderOption(option.data.value)}
    />
  );
});

export default AccessTypeSelect;
