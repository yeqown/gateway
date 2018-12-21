<template>
  <div style="text-align:left">
    <el-row>
      <h4>插件列表及状态展示</h4>
    </el-row>
    <el-row>
      <el-col :span="24" :offset="0">
        <div class="plugin-status-wrapper" v-for="(plg, idx) in plugins" :key="plg.key">
          <el-row :class="plg.status">
            <el-col :span="2" :offset="1" style="line-height:60px">
              <el-switch
                v-model="plg.enabled"
                active-color="#13ce66"
                inactive-color="#ff4949"
                :disabled="plg.key=== 'plugin.proxy'"
                @change="hdlChange(idx)"
              ></el-switch>
            </el-col>
            <el-col :span="15" :offset="1">
              <el-alert
                :title="plg.name"
                :description="plg.description"
                :type="plg.status"
                show-icon
                :closable="false"
              ></el-alert>
            </el-col>
            <el-col :span="2" :offset="3" style="line-height:60px">
              <router-link :to="plg.to" class="link">前往设置</router-link>
            </el-col>
          </el-row>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script>
// import { Loading } from "element-ui";
export default {
  name: "Plugins",
  data() {
    return {
      plugins: [
        {
          name: "代理",
          description:
            "代理插件是一个用于控制网关对于每一个请求进行控制，重新调度的插件",
          key: "plugin.proxy",
          enabled: true,
          status: "success",
          to: "/configs/plugin/proxy"
        },
        {
          name: "缓存",
          description:
            "缓存插件是一个用于控制网关对于每一个请求进行控制，重新调度的插件",
          key: "plugin.cache",
          enabled: true,
          status: "warning",
          to: "/configs/plugin/cache"
        },
        {
          name: "流量控制",
          description:
            "流量控制插件是一个用于控制网关对于每一个请求进行控制，重新调度的插件",
          key: "plugin.ratelimit",
          enabled: true,
          status: "error",
          to: "/configs/plugin/ratelimit"
        }
      ]
    };
  },
  methods: {
    hdlChange(idx) {
      // this.plugins[idx].enabled = !this.plugins[idx].enabled
      let plg = this.plugins[idx];
      let msg = `插件 ${plg.name} 已关闭！`;
      if (plg.enabled) {
        msg = `插件 ${plg.name} 已启用！`;
      }

      const loading = this.$loading({
        lock: true,
        text: "请求服务中...",
        spinner: "el-icon-loading",
        background: "rgba(0, 0, 0, 0.7)"
      });
      setTimeout(() => {
        loading.close();
        this.$notify.success({
          title: "消息通知",
          message: msg.toUpperCase(),
          duration: 2000
        });
      }, 1000);
    }
  }
};
</script>

<style scoped>
.plugin-status-wrapper {
  margin-bottom: 1em;
}
.success {
  background-color: #f0f9eb;
}
.warning {
  background-color: #fdf6ec;
}
.error {
  background-color: #fef0f0;
}

.link {
  color: #795548;
}
</style>
