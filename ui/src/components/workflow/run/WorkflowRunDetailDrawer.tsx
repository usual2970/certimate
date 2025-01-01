import { useTranslation } from "react-i18next";
import { useControllableValue } from "ahooks";
import { Alert, Drawer, Typography } from "antd";
import dayjs from "dayjs";

import Show from "@/components/Show";
import { type WorkflowRunModel } from "@/domain/workflowRun";
import { useTriggerElement } from "@/hooks";

export type WorkflowRunDetailDrawerProps = {
  data?: WorkflowRunModel;
  loading?: boolean;
  open?: boolean;
  trigger?: React.ReactNode;
  onOpenChange?: (open: boolean) => void;
};

const WorkflowRunDetailDrawer = ({ data, loading, trigger, ...props }: WorkflowRunDetailDrawerProps) => {
  const { t } = useTranslation();

  const [open, setOpen] = useControllableValue<boolean>(props, {
    valuePropName: "open",
    defaultValuePropName: "defaultOpen",
    trigger: "onOpenChange",
  });

  const triggerDom = useTriggerElement(trigger, { onClick: () => setOpen(true) });

  return (
    <>
      {triggerDom}

      <Drawer closable destroyOnClose open={open} loading={loading} placement="right" title={`runlog-${data?.id}`} width={640} onClose={() => setOpen(false)}>
        <Show when={!!data}>
          <Show when={data!.succeed}>
            <Alert showIcon type="success" message={<Typography.Text type="success">{t("workflow_run.props.status.succeeded")}</Typography.Text>} />
          </Show>

          <Show when={!!data!.error}>
            <Alert showIcon type="error" message={<Typography.Text type="danger">{t("workflow_run.props.status.failed")}</Typography.Text>} />
          </Show>

          <div className="mt-4 p-4 bg-black text-stone-200 rounded-md">
            <div className="flex flex-col space-y-3">
              {data!.log.map((item, i) => {
                return (
                  <div key={i} className="flex flex-col space-y-2">
                    <div>{item.nodeName}</div>
                    <div className="flex flex-col space-y-1">
                      {item.outputs.map((output, j) => {
                        return (
                          <div key={j} className="flex text-sm space-x-2">
                            <div>[{dayjs(output.time).format("YYYY-MM-DD HH:mm:ss")}]</div>
                            {output.error ? <div className="text-red-500">{output.error}</div> : <div>{output.content}</div>}
                          </div>
                        );
                      })}
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        </Show>
      </Drawer>
    </>
  );
};

export default WorkflowRunDetailDrawer;
