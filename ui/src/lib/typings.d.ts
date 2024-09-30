declare type Recordable<T = any> = Record<string, T>;

declare interface ViteEnv {
  VITE_API_URL: string;
  VITE_PORT: number;
  VITE_DROP_LOG: boolean;
  VITE_BUILD_GZIP: boolean;
}
