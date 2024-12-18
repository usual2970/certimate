import { useEffect, useState } from "react";
import { Avatar, Select, Space, Typography, type SelectProps } from "antd";

import { accessProvidersMap, type AccessModel } from "@/domain/access";
import { useAccessStore } from "@/stores/access";

export type AccessTypeSelectProps = Omit<
  SelectProps,
  "filterOption" | "filterSort" | "labelRender" | "options" | "optionFilterProp" | "optionLabelProp" | "optionRender"
> & {
  filter?: (record: AccessModel) => boolean;
};

const AccessSelect = ({ filter, ...props }: AccessTypeSelectProps) => {
  const { accesses, fetchAccesses } = useAccessStore();
  useEffect(() => {
    fetchAccesses();
  }, [fetchAccesses]);

  const [options, setOptions] = useState<Array<{ key: string; value: string; label: string; data: AccessModel }>>([]);
  useEffect(() => {
    const items = filter != null ? accesses.filter(filter) : accesses;
    setOptions(
      items.map((item) => ({
        key: item.id,
        value: item.id,
        label: item.name,
        data: item,
      }))
    );
  }, [accesses, filter]);

  const renderOption = (key: string) => {
    const access = accesses.find((e) => e.id === key);
    if (!access) {
      return (
        <Space className="flex-grow max-w-full truncate" size={4}>
          <Avatar size="small" />
          <Typography.Text className="leading-loose" ellipsis>
            {key}
          </Typography.Text>
        </Space>
      );
    }

    const provider = accessProvidersMap.get(access.configType);
    return (
      <Space className="flex-grow max-w-full truncate" size={4}>
        <Avatar src={provider?.icon} size="small" />
        <Typography.Text className="leading-loose" ellipsis>
          {access.name}
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
      optionFilterProp="label"
      optionLabelProp={undefined}
      optionRender={(option) => renderOption(option.data.value)}
    />
  );
};

export default AccessSelect;
