import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import { setupStore, usePermissionStore } from './stores';
import { permission } from './directives/permission';
import 'ant-design-vue/dist/reset.css';
import './assets/styles/index.css';

const app = createApp(App);

// 配置状态管理
setupStore(app);

// 注册全局指令
app.directive('permission', permission);

// 重置路由生成状态（确保每次刷新页面都重新加载路由）
const permissionStore = usePermissionStore();
permissionStore.isRoutesGenerated = false;

// 配置路由
app.use(router);

app.mount('#app');
