import { WorkflowOutput, WorkflowRunLog, WorkflowRunLogItem } from "@/domain/workflow";
import { Sheet, SheetContent, SheetDescription, SheetHeader, SheetTitle } from "../ui/sheet";
import { Check, X } from "lucide-react";

type WorkflowLogDetailProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  log?: WorkflowRunLog;
};
const WorkflowLogDetail = ({ open, onOpenChange, log }: WorkflowLogDetailProps) => {
  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent className="sm:max-w-5xl">
        <SheetHeader>
          <SheetTitle>日志</SheetTitle>
        </SheetHeader>

        <div className="flex flex-col">
          {log?.succeed ? (
            <div className="mt-5 flex justify-between bg-green-100 p-5 rounded-md items-center">
              <div className="flex space-x-2 items-center">
                <div className="w-8 h-8 bg-green-500 flex items-center justify-center rounded-full text-white">
                  <Check size={18} />
                </div>
                <div className="text-stone-700">成功</div>
              </div>

              <div className="text-muted-foreground">{new Date(log.created).toLocaleString()}</div>
            </div>
          ) : (
            <div className="mt-5 flex justify-between bg-green-100 p-5 rounded-md items-center">
              <div className="flex space-x-2 items-center">
                <div className="w-8 h-8 bg-red-500 flex items-center justify-center rounded-full text-white">
                  <X size={18} />
                </div>
                <div className="text-stone-700">失败</div>
              </div>

              <div className="text-red-500">{log?.error}</div>

              <div className="text-muted-foreground">{log?.created && new Date(log.created).toLocaleString()}</div>
            </div>
          )}

          <div className="bg-black p-5 mt-5 rounded-md text-stone-200 flex flex-col space-y-3">
            {log?.log.map((item: WorkflowRunLogItem, i) => {
              return (
                <div key={i} className="flex flex-col space-y-2">
                  <div className="">{item.nodeName}</div>
                  <div className="flex flex-col space-y-1">
                    {item.outputs.map((output: WorkflowOutput) => {
                      return (
                        <>
                          <div className="flex text-sm space-x-2">
                            <div>[{output.time}]</div>
                            {output.error ? (
                              <>
                                <div className="text-red-500">{output.error}</div>
                              </>
                            ) : (
                              <>
                                <div>{output.content}</div>
                              </>
                            )}
                          </div>
                        </>
                      );
                    })}
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </SheetContent>
    </Sheet>
  );
};

export default WorkflowLogDetail;

