import React from "react";
import { useTranslation } from "react-i18next";
import { Avatar, Select, Space, Typography, type SelectProps } from "antd";

import { accessProvidersMap } from "@/domain/access";

export type AccessProviderSelectProps = Omit<SelectProps, "options" | "optionFilterProp" | "optionLabelProp" | "optionRender"> & {
  className?: string;
};

const AccessProviderSelect = React.memo((props: AccessProviderSelectProps) => {
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
            <Space className="max-w-full truncate" align="center" size={4}>
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
        <Space className="max-w-full truncate" align="center" size={4}>
          <Avatar src={accessProvidersMap.get(option.data.value)?.icon} size="small" />
          <Typography.Text ellipsis>{t(accessProvidersMap.get(option.data.value)?.name ?? "")}</Typography.Text>
        </Space>
      )}
    />
  );
});

export default AccessProviderSelect;
