import { memo, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { Avatar, Select, type SelectProps, Space, Typography } from "antd";

import { type ApplyCAProvider, applyCAProvidersMap } from "@/domain/provider";

export type ApplyCAProviderSelectProps = Omit<
  SelectProps,
  "filterOption" | "filterSort" | "labelRender" | "options" | "optionFilterProp" | "optionLabelProp" | "optionRender"
> & {
  filter?: (record: ApplyCAProvider) => boolean;
};

const ApplyCAProviderSelect = ({ filter, ...props }: ApplyCAProviderSelectProps) => {
  const { t } = useTranslation();

  const [options, setOptions] = useState<Array<{ key: string; value: string; label: string; data: ApplyCAProvider }>>([]);
  useEffect(() => {
    const allItems = Array.from(applyCAProvidersMap.values());
    const filteredItems = filter != null ? allItems.filter(filter) : allItems;
    setOptions([
      {
        key: "",
        value: "",
        label: "provider.default_ca_provider.label",
        data: {} as ApplyCAProvider,
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

    const provider = applyCAProvidersMap.get(key);
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
      filterOption={(inputValue, option) => {
        if (!option) return false;

        const value = inputValue.toLowerCase();
        return option.value.toLowerCase().includes(value) || option.label.toLowerCase().includes(value);
      }}
      labelRender={({ label, value }) => {
        if (!label) {
          return <Typography.Text type="secondary">{props.placeholder || t("provider.default_ca_provider.label")}</Typography.Text>;
        }

        return renderOption(value as string);
      }}
      options={options}
      optionFilterProp={undefined}
      optionLabelProp={undefined}
      optionRender={(option) => renderOption(option.data.value)}
    />
  );
};

export default memo(ApplyCAProviderSelect);
