import { useState } from "react";
import { ExpandOutlined as ExpandOutlinedIcon, MinusOutlined as MinusOutlinedIcon, PlusOutlined as PlusOutlinedIcon } from "@ant-design/icons";
import { Button, Card, Typography } from "antd";

import WorkflowElements from "@/components/workflow/WorkflowElements";
import { mergeCls } from "@/utils/css";

export type WorkflowElementsProps = {
  className?: string;
  style?: React.CSSProperties;
  disabled?: boolean;
};

const WorkflowElementsContainer = ({ className, style, disabled }: WorkflowElementsProps) => {
  const [scale, setScale] = useState(1);

  return (
    <div className={mergeCls("relative size-full overflow-hidden", className)} style={style}>
      <div className="size-full overflow-auto">
        <div className="relative z-[1]">
          <div className="origin-center transition-transform duration-300" style={{ zoom: `${scale}` }}>
            <div className="p-4">
              <WorkflowElements disabled={disabled} />
            </div>
          </div>
        </div>
      </div>

      <Card className="absolute bottom-4 right-6 z-[2] rounded-lg p-2 shadow-lg" styles={{ body: { padding: 0 } }}>
        <div className="flex items-center gap-2">
          <Button icon={<MinusOutlinedIcon />} disabled={scale <= 0.5} onClick={() => setScale((s) => Math.max(0.5, s - 0.1))} />
          <Typography.Text className="min-w-[3em] text-center">{Math.round(scale * 100)}%</Typography.Text>
          <Button icon={<PlusOutlinedIcon />} disabled={scale >= 2} onClick={() => setScale((s) => Math.min(2, s + 0.1))} />
          <Button icon={<ExpandOutlinedIcon />} onClick={() => setScale(1)} />
        </div>
      </Card>
    </div>
  );
};

export default WorkflowElementsContainer;
