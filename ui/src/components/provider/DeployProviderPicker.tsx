import { memo, useState } from "react";
import { useTranslation } from "react-i18next";
import { useDebounceEffect } from "ahooks";
import { Avatar, Card, Col, Empty, Flex, Input, Row, Typography } from "antd";

import Show from "@/components/Show";
import { deployProvidersMap } from "@/domain/provider";

export type DeployProviderPickerProps = {
  className?: string;
  style?: React.CSSProperties;
  onSelect?: (value: string) => void;
};

const DeployProviderPicker = ({ className, style, onSelect }: DeployProviderPickerProps) => {
  const { t } = useTranslation();

  const allProviders = Array.from(deployProvidersMap.values());
  const [providers, setProviders] = useState(allProviders);
  const [keyword, setKeyword] = useState<string>();
  useDebounceEffect(
    () => {
      if (keyword) {
        setProviders(
          allProviders.filter((provider) => {
            const value = keyword.toLowerCase();
            return provider.type.toLowerCase().includes(value) || provider.name.toLowerCase().includes(value);
          })
        );
      } else {
        setProviders(allProviders);
      }
    },
    [keyword],
    { wait: 300 }
  );

  const handleProviderTypeSelect = (value: string) => {
    onSelect?.(value);
  };

  return (
    <div className={className} style={style}>
      <Input.Search placeholder={t("workflow_node.deploy.search.provider_type.placeholder")} onChange={(e) => setKeyword(e.target.value.trim())} />

      <div className="mt-4">
        <Show when={providers.length > 0} fallback={<Empty image={Empty.PRESENTED_IMAGE_SIMPLE} />}>
          <Row gutter={[16, 16]}>
            {providers.map((provider, index) => {
              return (
                <Col key={index} span={12}>
                  <Card
                    className="w-full h-16 shadow-sm overflow-hidden"
                    styles={{ body: { height: "100%", padding: "0.5rem 1rem" } }}
                    hoverable
                    onClick={() => {
                      handleProviderTypeSelect(provider.type);
                    }}
                  >
                    <Flex className="size-full overflow-hidden" align="center" gap={8}>
                      <Avatar src={provider.icon} size="small" />
                      <Typography.Text className="line-clamp-2">{t(provider.name)}</Typography.Text>
                    </Flex>
                  </Card>
                </Col>
              );
            })}
          </Row>
        </Show>
      </div>
    </div>
  );
};

export default memo(DeployProviderPicker);
