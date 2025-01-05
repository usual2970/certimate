import { useState } from "react";
import { useTranslation } from "react-i18next";
import { Button, Space } from "antd";

import Panel from "./Panel";
import PanelContext, { type ShowPanelOptions, type ShowPanelWithConfirmOptions } from "./PanelContext";

export const PanelProvider = ({ children }: { children: React.ReactNode }) => {
  const { t } = useTranslation();

  const [open, setOpen] = useState(false);
  const [options, setOptions] = useState<ShowPanelOptions>();

  const showPanel = (options: ShowPanelOptions) => {
    setOpen(true);
    setOptions(options);
  };

  const showPanelWithConfirm = (options: ShowPanelWithConfirmOptions) => {
    const updateOptionsFooter = (confirmLoading: boolean) => {
      setOptions({
        ...options,
        footer: (
          <Space className="w-full justify-end">
            <Button
              {...options.cancelButtonProps}
              onClick={() => {
                if (confirmLoading) return;

                options.onCancel?.();

                hidePanel();
              }}
            >
              {options.cancelText ?? t("common.button.cancel")}
            </Button>
            <Button
              loading={confirmLoading}
              type={options.okButtonProps?.type ?? "primary"}
              {...options.okButtonProps}
              onClick={async () => {
                updateOptionsFooter(true);

                try {
                  const ret = await options.onOk?.();
                  if (ret != null && !ret) return;
                } catch {
                  return;
                } finally {
                  updateOptionsFooter(false);
                }

                hidePanel();
              }}
            >
              {options.okText ?? t("common.button.ok")}
            </Button>
          </Space>
        ),
        onClose: () => Promise.resolve(!confirmLoading),
      });
    };

    showPanel(options);
    updateOptionsFooter(false);
  };

  const hidePanel = () => {
    setOpen(false);
    setOptions(undefined);
  };

  const handleOpenChange = (open: boolean) => {
    setOpen(open);

    if (!open) {
      setOptions(undefined);
    }
  };

  return (
    <PanelContext.Provider value={{ open, show: showPanel, confirm: showPanelWithConfirm, hide: hidePanel }}>
      {children}

      <Panel open={open} {...options} onOpenChange={handleOpenChange}>
        {options?.children}
      </Panel>
    </PanelContext.Provider>
  );
};
