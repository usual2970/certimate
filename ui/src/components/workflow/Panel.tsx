// components/AddNodePanel.tsx
import { ScrollArea } from "../ui/scroll-area";
import { Sheet, SheetContent, SheetTitle } from "../ui/sheet";

type AddNodePanelProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  children: React.ReactNode;
  name: string;
};

const Panel = ({ open, onOpenChange, children, name }: AddNodePanelProps) => {
  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent className="sm:max-w-[640px] p-0">
        <SheetTitle className="bg-primary px-6 py-4 text-white">{name}</SheetTitle>

        <ScrollArea className="px-6 py-4 flex-col space-y-5 h-[90vh]">{children}</ScrollArea>
      </SheetContent>
    </Sheet>
  );
};

export default Panel;
