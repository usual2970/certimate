import { memo, useEffect, useMemo, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { Avatar, Card, Col, Empty, Flex, Input, type InputRef, Row, Tabs, Tooltip, Typography } from "antd";

import Show from "@/components/Show";
import { DEPLOY_CATEGORIES, deployProvidersMap } from "@/domain/provider";

export type DeployProviderPickerProps = {
  className?: string;
  style?: React.CSSProperties;
  autoFocus?: boolean;
  placeholder?: string;
  onSelect?: (value: string) => void;
};

const DeployProviderPicker = ({ className, style, autoFocus, placeholder, onSelect }: DeployProviderPickerProps) => {
  const { t } = useTranslation();

  const [keyword, setKeyword] = useState<string>();
  const keywordInputRef = useRef<InputRef>(null);
  useEffect(() => {
    if (autoFocus) {
      setTimeout(() => keywordInputRef.current?.focus(), 1);
    }
  }, []);

  const [category, setCategory] = useState<string>(DEPLOY_CATEGORIES.ALL);

  const providers = useMemo(() => {
    return Array.from(deployProvidersMap.values())
      .filter((provider) => {
        if (keyword) {
          const value = keyword.toLowerCase();
          return provider.type.toLowerCase().includes(value) || t(provider.name).toLowerCase().includes(value);
        }

        return true;
      })
      .filter((provider) => {
        if (category && category !== DEPLOY_CATEGORIES.ALL) {
          return provider.category === category;
        }

        return true;
      });
  }, [keyword, category]);

  const handleProviderTypeSelect = (value: string) => {
    onSelect?.(value);
  };

  return (
    <div className={className} style={style}>
      <Input.Search ref={keywordInputRef} placeholder={placeholder} onChange={(e) => setKeyword(e.target.value.trim())} />

      <div className="mt-4">
        <Flex>
          <Tabs
            defaultActiveKey={DEPLOY_CATEGORIES.ALL}
            items={[
              DEPLOY_CATEGORIES.ALL,
              DEPLOY_CATEGORIES.CDN,
              DEPLOY_CATEGORIES.STORAGE,
              DEPLOY_CATEGORIES.LOADBALANCE,
              DEPLOY_CATEGORIES.FIREWALL,
              DEPLOY_CATEGORIES.AV,
              DEPLOY_CATEGORIES.SERVERLESS,
              DEPLOY_CATEGORIES.WEBSITE,
              DEPLOY_CATEGORIES.OTHER,
            ].map((key) => ({
              key: key,
              label: t(`provider.category.${key}`),
            }))}
            size="small"
            tabBarStyle={{ marginLeft: "-1rem" }}
            tabPosition="left"
            onChange={(key) => setCategory(key)}
          />

          <div className="flex-1">
            <Show when={providers.length > 0} fallback={<Empty image={Empty.PRESENTED_IMAGE_SIMPLE} />}>
              <Row gutter={[16, 16]}>
                {providers.map((provider, index) => {
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
                        <Tooltip title={t(provider.name)} mouseEnterDelay={1}>
                          <Flex className="size-full overflow-hidden" align="center" gap={8}>
                            <Avatar src={provider.icon} size="small" />
                            <Typography.Text className="line-clamp-2 flex-1">{t(provider.name)}</Typography.Text>
                          </Flex>
                        </Tooltip>
                      </Card>
                    </Col>
                  );
                })}
              </Row>
            </Show>
          </div>
        </Flex>
      </div>
    </div>
  );
};

export default memo(DeployProviderPicker);
