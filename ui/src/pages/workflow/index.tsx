import { Button } from "@/components/ui/button";
import { Plus } from "lucide-react";
import { useNavigate } from "react-router-dom";

const Workflow = () => {
  const navigate = useNavigate();
  const handleCreateClick = () => {
    navigate("/workflow/detail");
  };
  return (
    <>
      <div className="flex justify-between items-center">
        <div className="text-muted-foreground">工作流</div>
        <Button onClick={handleCreateClick}>
          <Plus size={16} />
          新建工作流
        </Button>
      </div>
    </>
  );
};

export default Workflow;
