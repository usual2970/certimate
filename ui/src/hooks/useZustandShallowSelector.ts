import { useRef } from "react";
import { isArray, pick } from "radash";
import { shallow } from "zustand/shallow";

type MaybeMany<T> = T | readonly T[];

export type UseZustandShallowSelectorReturns<T extends object, TKeys extends keyof T> = (state: T) => Pick<T, TKeys>;

/**
 * 选择并获取指定的状态。
 * 基于 `zustand.useShallow` 二次封装，以减少样板代码。
 * @param {Array} paths 要选择的状态键名。
 * @returns {UseZustandShallowSelectorReturns}
 *
 * @example
 * ```js
 * // 使用示例：
 * const { foo, bar, baz } = useStore(useZustandShallowSelector(["foo", "bar", "baz"]));
 *
 * // 以上代码等效于：
 * const { foo, bar, baz } = useStore((state) => ({
 *  foo: state.foo,
 *  bar: state.bar,
 *  baz: state.baz,
 * }));
 * ```
 */
const useZustandShallowSelector = <T extends object, TKeys extends keyof T>(paths: MaybeMany<TKeys>): UseZustandShallowSelectorReturns<T, TKeys> => {
  const prev = useRef<Pick<T, TKeys>>({} as Pick<T, TKeys>);

  return (state: T) => {
    if (state) {
      const next = pick(state, isArray(paths) ? paths : [paths]);
      return shallow(prev.current, next) ? prev.current : (prev.current = next);
    }
    return prev.current;
  };
};

export default useZustandShallowSelector;
