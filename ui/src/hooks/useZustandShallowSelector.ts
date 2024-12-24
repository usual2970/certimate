import { pick, isArray } from "radash";

import { useRef } from "react";
import { shallow } from "zustand/shallow";

type MaybeMany<T> = T | readonly T[];

export default function <T extends object, TKeys extends keyof T>(paths: MaybeMany<TKeys>): (state: T) => Pick<T, TKeys> {
  const prev = useRef<Pick<T, TKeys>>({} as Pick<T, TKeys>);

  return (state: T) => {
    if (state) {
      const next = pick(state, isArray(paths) ? paths : [paths]);
      return shallow(prev.current, next) ? prev.current : (prev.current = next);
    }
    return prev.current;
  };
}
