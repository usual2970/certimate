import React, { useEffect } from "react";
import { Select, SelectContent, SelectGroup, SelectItem, SelectLabel, SelectTrigger, SelectValue } from "../ui/select";
import { accessProvidersMap } from "@/domain/access";
import { useTranslation } from "react-i18next";
import { useConfigContext } from "@/providers/config";
import { deployTargetsMap } from "@/domain/domain";

type AccessSelectProps = {
  providerType: string;
  value: string;
  onValueChange: (val: string) => void;
};
const AccessSelect = ({ value, onValueChange, providerType }: AccessSelectProps) => {
  const [localValue, setLocalValue] = React.useState<string>("");
  const { t } = useTranslation();
  const {
    config: { accesses },
  } = useConfigContext();

  useEffect(() => {
    setLocalValue(value);
  }, [value]);

  const targetAccesses = accesses.filter((item) => {
    console.log(item, providerType);
    return item.configType === deployTargetsMap.get(providerType)?.provider;
  });

  return (
    <>
      <Select
        value={localValue}
        onValueChange={(val: string) => {
          setLocalValue(val);
          onValueChange(val);
        }}
      >
        <SelectTrigger className="mt-2">
          <SelectValue placeholder={t("domain.deployment.form.access.placeholder")} />
        </SelectTrigger>
        <SelectContent>
          <SelectGroup>
            <SelectLabel>{t("domain.deployment.form.access.list")}</SelectLabel>
            {targetAccesses.map((item) => (
              <SelectItem key={item.id} value={item.id}>
                <div className="flex items-center space-x-2">
                  <img className="w-6" src={accessProvidersMap.get(item.configType)?.icon} />
                  <div>{item.name}</div>
                </div>
              </SelectItem>
            ))}
          </SelectGroup>
        </SelectContent>
      </Select>
    </>
  );
};

export default AccessSelect;
