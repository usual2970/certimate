import { useMemo, useRef } from "react";
import { json } from "@codemirror/lang-json";
import { yaml } from "@codemirror/lang-yaml";
import { StreamLanguage } from "@codemirror/language";
import { powerShell } from "@codemirror/legacy-modes/mode/powershell";
import { shell } from "@codemirror/legacy-modes/mode/shell";
import { basicSetup } from "@uiw/codemirror-extensions-basic-setup";
import { vscodeDark, vscodeLight } from "@uiw/codemirror-theme-vscode";
import CodeMirror, { type ReactCodeMirrorProps, type ReactCodeMirrorRef } from "@uiw/react-codemirror";
import { useFocusWithin } from "ahooks";
import { theme } from "antd";

import { useBrowserTheme } from "@/hooks";
import { mergeCls } from "@/utils/css";

export interface CodeInputProps extends Omit<ReactCodeMirrorProps, "extensions" | "lang" | "theme"> {
  disabled?: boolean;
  language?: string | string[];
}

const CodeInput = ({ className, style, disabled, language, ...props }: CodeInputProps) => {
  const { token: themeToken } = theme.useToken();

  const { theme: browserTheme } = useBrowserTheme();

  const cmRef = useRef<ReactCodeMirrorRef>(null);
  const isFocusWithin = useFocusWithin(cmRef.current?.editor);

  const cmTheme = useMemo(() => {
    if (browserTheme === "dark") {
      return vscodeDark;
    }
    return vscodeLight;
  }, [browserTheme]);

  const cmExtensions = useMemo(() => {
    const temp: NonNullable<ReactCodeMirrorProps["extensions"]> = [
      basicSetup({
        foldGutter: false,
        dropCursor: false,
        allowMultipleSelections: false,
        indentOnInput: false,
      }),
    ];

    const langs = Array.isArray(language) ? language : [language];
    langs.forEach((lang) => {
      switch (lang) {
        case "shell":
          temp.push(StreamLanguage.define(shell));
          break;
        case "json":
          temp.push(json());
          break;
        case "powershell":
          temp.push(StreamLanguage.define(powerShell));
          break;
        case "yaml":
          temp.push(yaml());
          break;
      }
    });

    return temp;
  }, [language]);

  return (
    <div
      className={mergeCls(className, `hover:border-[${themeToken.colorPrimaryBorderHover}]`)}
      style={{
        ...(style ?? {}),
        border: `1px solid ${isFocusWithin ? (themeToken.Input?.activeBorderColor ?? themeToken.colorPrimaryBorder) : themeToken.colorBorder}`,
        borderRadius: `${themeToken.borderRadius}px`,
        backgroundColor: disabled ? themeToken.colorBgContainerDisabled : themeToken.colorBgContainer,
        boxShadow: isFocusWithin ? themeToken.Input?.activeShadow : undefined,
        overflow: "hidden",
      }}
    >
      <CodeMirror
        ref={cmRef}
        height="100%"
        style={{ height: "100%" }}
        {...props}
        basicSetup={{
          foldGutter: false,
          dropCursor: false,
          allowMultipleSelections: false,
          indentOnInput: false,
        }}
        extensions={cmExtensions}
        theme={cmTheme}
      />
    </div>
  );
};

export default CodeInput;
