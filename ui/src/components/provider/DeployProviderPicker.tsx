import { memo, useState } from "react";
import { useTranslation } from "react-i18next";
import { Avatar, Card, Col, Empty, Flex, Input, Row, Typography } from "antd";

import Show from "@/components/Show";
import { deployProvidersMap } from "@/domain/provider";

export type DeployProviderPickerProps = {
  className?: string;
  style?: React.CSSProperties;
  placeholder?: string;
  onSelect?: (value: string) => void;
};

const DeployProviderPicker = ({ className, style, placeholder, onSelect }: DeployProviderPickerProps) => {
  const { t } = useTranslation();

  const [keyword, setKeyword] = useState<string>();

  const providers = Array.from(deployProvidersMap.values());
  const filteredProviders = providers.filter((provider) => {
    if (keyword) {
      const value = keyword.toLowerCase();
      return provider.type.toLowerCase().includes(value) || provider.name.toLowerCase().includes(value);
    }

    return true;
  });

  const handleProviderTypeSelect = (value: string) => {
    onSelect?.(value);
  };

  return (
    <div className={className} style={style}>
      <Input.Search placeholder={placeholder} onChange={(e) => setKeyword(e.target.value.trim())} />

      <div className="mt-4">
        <Show when={filteredProviders.length > 0} fallback={<Empty image={Empty.PRESENTED_IMAGE_SIMPLE} />}>
          <Row gutter={[16, 16]}>
            {filteredProviders.map((provider, index) => {
              return (
                <Col key={index} xs={24} md={12} span={12}>
                  <Card
                    className="h-16 w-full overflow-hidden shadow-sm"
                    styles={{ body: { height: "100%", padding: "0.5rem 1rem" } }}
                    hoverable
                    onClick={() => {
                      handleProviderTypeSelect(provider.type);
                    }}
                  >
                    <Flex className="size-full overflow-hidden" align="center" gap={8}>
                      <Avatar src={provider.icon} size="small" />
                      <Typography.Text className="line-clamp-2 flex-1">{t(provider.name)}</Typography.Text>
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
