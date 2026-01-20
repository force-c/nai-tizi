declare global {
  interface Window {
    // 扩展 window 对象
  }
}

export {};

// 环境变量类型
interface ImportMetaEnv {
  readonly VITE_APP_TITLE: string;
  readonly VITE_API_BASE_URL: string;
  readonly VITE_CLIENT_KEY: string;
  readonly VITE_CLIENT_SECRET: string;
  readonly VITE_USE_MOCK: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
