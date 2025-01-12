import { memo } from "react";
import { useTranslation } from "react-i18next";
import { Avatar, Select, type SelectProps, Space, Typography } from "antd";

import { applyDNSProvidersMap } from "@/domain/provider";

export type ApplyDNSProviderSelectProps = Omit<
  SelectProps,
  "filterOption" | "filterSort" | "labelRender" | "options" | "optionFilterProp" | "optionLabelProp" | "optionRender"
>;

const ApplyDNSProviderSelect = (props: ApplyDNSProviderSelectProps) => {
  const { t } = useTranslation();

  const options = Array.from(applyDNSProvidersMap.values()).map((item) => ({
    key: item.type,
    value: item.type,
    label: t(item.name),
  }));

  const renderOption = (key: string) => {
    const provider = applyDNSProvidersMap.get(key);
    return (
      <Space className="max-w-full grow overflow-hidden truncate" size={4}>
        <Avatar src={provider?.icon} size="small" />
        <Typography.Text className="leading-loose" ellipsis>
          {t(provider?.name ?? "")}
        </Typography.Text>
      </Space>
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
};

export default memo(ApplyDNSProviderSelect);
