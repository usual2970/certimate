declare global {
  type Nullish<T> = {
    [P in keyof T]?: T[P] | null | undefined;
  };
}

export {};
