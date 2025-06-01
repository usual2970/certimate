import { memo, useEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { Avatar, Select, type SelectProps, Space, Tag, Typography, theme } from "antd";

import Show from "@/components/Show";
import { ACCESS_USAGES, type AccessProvider, type AccessUsageType, accessProvidersMap } from "@/domain/provider";

export type AccessProviderSelectProps = Omit<
  SelectProps,
  "filterOption" | "filterSort" | "labelRender" | "options" | "optionFilterProp" | "optionLabelProp" | "optionRender"
> & {
  filter?: (record: AccessProvider) => boolean;
  showOptionTags?: boolean | { [key in AccessUsageType]?: boolean };
};

const AccessProviderSelect = ({ filter, showOptionTags, ...props }: AccessProviderSelectProps = { showOptionTags: true }) => {
  const { t } = useTranslation();

  const { token: themeToken } = theme.useToken();

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

  const showOptionTagForDNS = useMemo(() => {
    return typeof showOptionTags === "object" ? !!showOptionTags[ACCESS_USAGES.DNS] : !!showOptionTags;
  }, [showOptionTags]);
  const showOptionTagForHosting = useMemo(() => {
    return typeof showOptionTags === "object" ? !!showOptionTags[ACCESS_USAGES.HOSTING] : !!showOptionTags;
  }, [showOptionTags]);
  const showOptionTagForCA = useMemo(() => {
    return typeof showOptionTags === "object" ? !!showOptionTags[ACCESS_USAGES.CA] : !!showOptionTags;
  }, [showOptionTags]);
  const showOptionTagForNotification = useMemo(() => {
    return typeof showOptionTags === "object" ? !!showOptionTags[ACCESS_USAGES.NOTIFICATION] : !!showOptionTags;
  }, [showOptionTags]);

  const renderOption = (key: string) => {
    const provider = accessProvidersMap.get(key) ?? ({ type: "", name: "", icon: "", usages: [] } as unknown as AccessProvider);
    return (
      <div className="flex max-w-full items-center justify-between gap-4 overflow-hidden">
        <Space className="max-w-full grow truncate" size={4}>
          <Avatar shape="square" src={provider.icon} size="small" />
          <Typography.Text className="leading-loose" type={provider.builtin ? "secondary" : undefined} ellipsis>
            {t(provider.name)}
          </Typography.Text>
        </Space>
        <div className="origin-right scale-[75%]">
          <Show when={provider.builtin}>
            <Tag>{t("access.props.provider.builtin")}</Tag>
          </Show>
          <Show when={showOptionTagForDNS && provider.usages.includes(ACCESS_USAGES.DNS)}>
            <Tag color="orange">{t("access.props.provider.usage.dns")}</Tag>
          </Show>
          <Show when={showOptionTagForHosting && provider.usages.includes(ACCESS_USAGES.HOSTING)}>
            <Tag color="geekblue">{t("access.props.provider.usage.hosting")}</Tag>
          </Show>
          <Show when={showOptionTagForCA && provider.usages.includes(ACCESS_USAGES.CA)}>
            <Tag color="magenta">{t("access.props.provider.usage.ca")}</Tag>
          </Show>
          <Show when={showOptionTagForNotification && provider.usages.includes(ACCESS_USAGES.NOTIFICATION)}>
            <Tag color="cyan">{t("access.props.provider.usage.notification")}</Tag>
          </Show>
        </div>
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

export default memo(AccessProviderSelect);
