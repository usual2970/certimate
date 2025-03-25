import { useEffect, useState } from "react";
import { Avatar, Select, type SelectProps, Space, Typography } from "antd";

import { type AccessModel } from "@/domain/access";
import { accessProvidersMap } from "@/domain/provider";
import { useZustandShallowSelector } from "@/hooks";
import { useAccessesStore } from "@/stores/access";

export type AccessTypeSelectProps = Omit<
  SelectProps,
  "filterOption" | "filterSort" | "labelRender" | "loading" | "options" | "optionFilterProp" | "optionLabelProp" | "optionRender"
> & {
  filter?: (record: AccessModel) => boolean;
};

const AccessSelect = ({ filter, ...props }: AccessTypeSelectProps) => {
  const { accesses, loadedAtOnce, fetchAccesses } = useAccessesStore(useZustandShallowSelector(["accesses", "loadedAtOnce", "fetchAccesses"]));
  useEffect(() => {
    fetchAccesses();
  }, []);

  const [options, setOptions] = useState<Array<{ key: string; value: string; label: string; data: AccessModel }>>([]);
  useEffect(() => {
    const filteredItems = filter != null ? accesses.filter(filter) : accesses;
    setOptions(
      filteredItems.map((item) => ({
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
        <Space className="max-w-full grow truncate" size={4}>
          <Avatar size="small" />
          <Typography.Text className="leading-loose" ellipsis>
            {key}
          </Typography.Text>
        </Space>
      );
    }

    const provider = accessProvidersMap.get(access.provider);
    return (
      <Space className="max-w-full grow truncate" size={4}>
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
      loading={!loadedAtOnce}
      options={options}
      optionFilterProp="label"
      optionLabelProp={undefined}
      optionRender={(option) => renderOption(option.data.value)}
    />
  );
};

export default AccessSelect;
