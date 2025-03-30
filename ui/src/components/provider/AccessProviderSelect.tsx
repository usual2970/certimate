import { memo, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { Avatar, Select, type SelectProps, Space, Tag, Typography } from "antd";

import { ACCESS_USAGES, type AccessProvider, accessProvidersMap } from "@/domain/provider";

export type AccessProviderSelectProps = Omit<
  SelectProps,
  "filterOption" | "filterSort" | "labelRender" | "options" | "optionFilterProp" | "optionLabelProp" | "optionRender"
> & {
  filter?: (record: AccessProvider) => boolean;
  showOptionTags?: boolean;
};

const AccessProviderSelect = ({ filter, showOptionTags, ...props }: AccessProviderSelectProps = { showOptionTags: true }) => {
  const { t } = useTranslation();

  const [options, setOptions] = useState<Array<{ key: string; value: string; label: string; data: AccessProvider }>>([]);
  useEffect(() => {
    const allItems = Array.from(accessProvidersMap.values());
    const filteredItems = filter != null ? allItems.filter(filter) : allItems;
    setOptions(
      filteredItems.map((item) => ({
        key: item.type,
        value: item.type,
        label: t(item.name),
        disabled: item.builtin,
        data: item,
      }))
    );
  }, [filter]);

  const renderOption = (key: string) => {
    const provider = accessProvidersMap.get(key);
    return (
      <div className="flex max-w-full items-center justify-between gap-4 overflow-hidden">
        <Space className="max-w-full grow truncate" size={4}>
          <Avatar src={provider?.icon} size="small" />
          <Typography.Text className="leading-loose" type={provider?.builtin ? "secondary" : undefined} delete={provider?.builtin ? true : undefined} ellipsis>
            {t(provider?.name ?? "")}
          </Typography.Text>
        </Space>
        {showOptionTags && (
          <div>
            {provider?.usages?.includes(ACCESS_USAGES.DNS) && <Tag color="peru">{t("access.props.provider.usage.dns")}</Tag>}
            {provider?.usages?.includes(ACCESS_USAGES.HOSTING) && <Tag color="royalblue">{t("access.props.provider.usage.hosting")}</Tag>}
            {/* {provider?.usages?.includes(ACCESS_USAGES.CA) && <Tag color="crimson">{t("access.props.provider.usage.ca")}</Tag>} */}
            {/* {provider?.usages?.includes(ACCESS_USAGES.NOTIFICATION) && <Tag color="mediumaquamarine">{t("access.props.provider.usage.notification")}</Tag>} */}
          </div>
        )}
      </div>
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
          return <Typography.Text type="secondary">{props.placeholder}</Typography.Text>;
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

export default memo(AccessProviderSelect);
