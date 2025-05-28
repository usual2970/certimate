import { memo, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { Avatar, Select, type SelectProps, Space, Typography, theme } from "antd";

import { type CAProvider, caProvidersMap } from "@/domain/provider";

export type CAProviderSelectProps = Omit<
  SelectProps,
  "filterOption" | "filterSort" | "labelRender" | "options" | "optionFilterProp" | "optionLabelProp" | "optionRender"
> & {
  filter?: (record: CAProvider) => boolean;
};

const CAProviderSelect = ({ filter, ...props }: CAProviderSelectProps) => {
  const { t } = useTranslation();

  const { token: themeToken } = theme.useToken();

  const [options, setOptions] = useState<Array<{ key: string; value: string; label: string; data: CAProvider }>>([]);
  useEffect(() => {
    const allItems = Array.from(caProvidersMap.values());
    const filteredItems = filter != null ? allItems.filter(filter) : allItems;
    setOptions([
      {
        key: "",
        value: "",
        label: t("provider.default_ca_provider.label"),
        data: {} as CAProvider,
      },
      ...filteredItems.map((item) => ({
        key: item.type,
        value: item.type,
        label: t(item.name),
        data: item,
      })),
    ]);
  }, [filter]);

  const renderOption = (key: string) => {
    if (key === "") {
      return (
        <Space className="max-w-full grow overflow-hidden truncate" size={4}>
          <Typography.Text className="italic leading-loose" type="secondary" ellipsis italic>
            {t("provider.default_ca_provider.label")}
          </Typography.Text>
        </Space>
      );
    }

    const provider = caProvidersMap.get(key);
    return (
      <Space className="max-w-full grow overflow-hidden truncate" size={4}>
        <Avatar shape="square" src={provider?.icon} size="small" />
        <Typography.Text className="leading-loose" ellipsis>
          {t(provider?.name ?? "")}
        </Typography.Text>
      </Space>
    );
  };

  return (
    <Select
      {...props}
      filterOption={(inputValue, option) => {
        if (!option) return false;

        const value = inputValue.toLowerCase();
        return option.value.toLowerCase().includes(value) || option.label.toLowerCase().includes(value);
      }}
      labelRender={({ value }) => {
        if (value != null) {
          return renderOption(value as string);
        }

        return <span style={{ color: themeToken.colorTextPlaceholder }}>{props.placeholder}</span>;
      }}
      options={options}
      optionFilterProp={undefined}
      optionLabelProp={undefined}
      optionRender={(option) => renderOption(option.data.value)}
    />
  );
};

export default memo(CAProviderSelect);
