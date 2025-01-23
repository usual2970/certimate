import { useTranslation } from "react-i18next";
import { useControllableValue } from "ahooks";
import { Alert, Drawer, Typography } from "antd";
import dayjs from "dayjs";

import Show from "@/components/Show";
import { WORKFLOW_RUN_STATUSES, type WorkflowRunModel } from "@/domain/workflowRun";
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

  const triggerEl = useTriggerElement(trigger, { onClick: () => setOpen(true) });

  return (
    <>
      {triggerEl}

      <Drawer destroyOnClose open={open} loading={loading} placement="right" title={`WorkflowRun #${data?.id}`} width={640} onClose={() => setOpen(false)}>
        <Show when={!!data}>
          <Show when={data!.status === WORKFLOW_RUN_STATUSES.SUCCEEDED}>
            <Alert showIcon type="success" message={<Typography.Text type="success">{t("workflow_run.props.status.succeeded")}</Typography.Text>} />
          </Show>

          <Show when={data!.status === WORKFLOW_RUN_STATUSES.FAILED}>
            <Alert showIcon type="error" message={<Typography.Text type="danger">{t("workflow_run.props.status.failed")}</Typography.Text>} />
          </Show>

          <div className="mt-4 rounded-md bg-black p-4 text-stone-200">
            <div className="flex flex-col space-y-4">
              {data!.logs?.map((item, i) => {
                return (
                  <div key={i} className="flex flex-col space-y-2">
                    <div className="font-semibold">{item.nodeName}</div>
                    <div className="flex flex-col space-y-1">
                      {item.outputs?.map((output, j) => {
                        return (
                          <div key={j} className="flex space-x-2 text-sm" style={{ wordBreak: "break-word" }}>
                            <div className="whitespace-nowrap">[{dayjs(output.time).format("YYYY-MM-DD HH:mm:ss")}]</div>
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
