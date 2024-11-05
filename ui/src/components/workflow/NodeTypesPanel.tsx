import { WorkflowNodeType } from "@/domain/workflow";
import { CloudUpload, GitFork, Megaphone, NotebookPen } from "lucide-react";

type NodeTypesPanelProps = {
  onTypeSelected: (type: WorkflowNodeType) => void;
};

const NodeTypesPanel = ({ onTypeSelected }: NodeTypesPanelProps) => {
  return (
    <>
      <div className="flex space-x-2">
        <div
          className="flex w-1/2 items-center space-x-2 hover:bg-stone-100 p-2 rounded-md"
          onClick={() => {
            onTypeSelected(WorkflowNodeType.Apply);
          }}
        >
          <div className="bg-primary h-12 w-12 flex items-center justify-center rounded-full">
            <NotebookPen className="text-white" size={18} />
          </div>

          <div className="text-slate-600">申请</div>
        </div>
        <div
          className="flex w-1/2 items-center space-x-2  hover:bg-stone-100 p-2 rounded-md"
          onClick={() => {
            onTypeSelected(WorkflowNodeType.Deploy);
          }}
        >
          <div className="bg-primary h-12 w-12 flex items-center justify-center rounded-full">
            <CloudUpload className="text-white" size={18} />
          </div>

          <div className="text-slate-600">部署</div>
        </div>
      </div>
      <div className="flex space-x-2">
        <div
          className="flex w-1/2 items-center space-x-2 hover:bg-stone-100 p-2 rounded-md"
          onClick={() => {
            onTypeSelected(WorkflowNodeType.Branch);
          }}
        >
          <div className="bg-primary h-12 w-12 flex items-center justify-center rounded-full">
            <GitFork className="text-white" size={18} />
          </div>

          <div className="text-slate-600">分支</div>
        </div>
        <div
          className="flex w-1/2 items-center space-x-2 hover:bg-stone-100 p-2 rounded-md"
          onClick={() => {
            onTypeSelected(WorkflowNodeType.Notify);
          }}
        >
          <div className="bg-primary h-12 w-12 flex items-center justify-center rounded-full">
            <Megaphone className="text-white" size={18} />
          </div>

          <div className="text-slate-600">推送</div>
        </div>
      </div>
    </>
  );
};

export default NodeTypesPanel;
