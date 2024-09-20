export type Setting = {
  id?: string;
  name?: string;
  content: EmailsSetting;
};

type EmailsSetting = {
  emails: string[];
};
