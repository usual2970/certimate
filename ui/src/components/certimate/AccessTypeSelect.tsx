import { Check, ChevronsUpDown } from "lucide-react";

import { cn } from "@/components/ui/utils";
import { Button } from "@/components/ui/button";
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from "@/components/ui/command";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { accessProvidersMap } from "@/domain/access";
import { useTranslation } from "react-i18next";
import { useEffect, useState } from "react";

type AccessTypeSelectProps = {
  value: string;
  onChange: (value: string) => void;
  placeholder: string;
  searchPlaceholder: string;
  className?: string;
};

export function AccessTypeSelect({ value, onChange, placeholder, searchPlaceholder, className }: AccessTypeSelectProps) {
  const [open, setOpen] = useState(false);
  const [locValue, setLocValue] = useState("");
  const { t } = useTranslation();
  const [search, setSearch] = useState("");
  const filteredProviders = Array.from(accessProvidersMap.entries());

  useEffect(() => {
    setLocValue(value);
  }, [value]);

  const handleOnSelect = (currentValue: string) => {
    const newValue = currentValue === locValue ? "" : currentValue;
    setLocValue(newValue);
    setSearch("");
    setOpen(false);
    onChange(newValue);
  };

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button variant="outline" role="combobox" aria-expanded={open} className={cn("justify-between z-50", className)}>
          {locValue ? (
            <div className="flex space-x-2 items-center">
              <img src={accessProvidersMap.get(locValue)?.icon} className="h-6 w-6" />
              <div>{t(accessProvidersMap.get(locValue)?.name ?? "")}</div>
            </div>
          ) : (
            <>{placeholder}</>
          )}
          <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className={cn("p-0  w-full")}>
        <Command className="">
          <CommandInput
            placeholder={searchPlaceholder}
            value={search}
            onValueChange={(val: string) => {
              setSearch(val);
            }}
          />
          <CommandList>
            <CommandEmpty>{t("access.authorization.form.type.search.notfound")}</CommandEmpty>
            <CommandGroup>
              {filteredProviders.map(([key, provider]) => (
                <CommandItem key={key} value={key} onSelect={handleOnSelect} keywords={provider.searchContent.split(":")}>
                  <Check className={cn("mr-2 h-4 w-4", locValue === key ? "opacity-100" : "opacity-0")} />
                  <div className="flex space-x-2">
                    <img src={provider.icon} className="h-6 w-6" />
                    <div className="font-medium">{t(provider.name)}</div>
                  </div>
                </CommandItem>
              ))}
            </CommandGroup>
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  );
}
