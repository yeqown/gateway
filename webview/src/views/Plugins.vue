<template>
  <div style="text-align:left">
    <el-row>
      <h4>插件列表</h4>
    </el-row>
    <el-row>
      <el-col :span="24" :offset="0">
        <div class="plugin-status-wrapper" v-for="(plg, idx) in plugins" :key="plg.name">
          <el-row :class="plg.status">
            <el-col :span="2" :offset="1" style="line-height:60px">
              <el-switch
                v-model="plg.enabled"
                active-color="#13ce66"
                inactive-color="#ff4949"
                :disabled="plg.name=== 'plugin.proxy'"
                @change="hdlChange(idx)"
              ></el-switch>
            </el-col>
            <el-col :span="15" :offset="1">
              <el-alert
                :title="plgCfg[plg.name].name"
                :description="plgCfg[plg.name].description"
                :type="(plg.status==='working')?'success':'info'"
                show-icon
                :closable="false"
              ></el-alert>
            </el-col>
            <el-col :span="2" :offset="3" style="line-height:60px">
              <router-link :to="plg.to || '/hell'" class="link">前往设置</router-link>
            </el-col>
          </el-row>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { basicapi } from "@/apis";
import { plgCfg } from "@/config";
export default {
  name: "Plugins",
  data() {
    return {
      plugins: [],
      plgCfg: plgCfg
    };
  },
  methods: {
    async hdlChange(idx) {
      // this.plugins[idx].enabled = !this.plugins[idx].enabled
      let plg = this.plugins[idx];
      let msg = `插件 ${plg.name} 已关闭！`;
      if (plg.enabled) {
        msg = `插件 ${plg.name} 已启用！`;
      }

      // console.log(idx, this.plugins[idx].enabled)
      await this.enablePlugin(idx, this.plugins[idx].enabled);
      await this.refresh();
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
    },
    refresh() {
      basicapi
        .getPluginsList()
        .then(data => {
          this.plugins = data.plugins;
        })
        .catch(err => console.error(err));
    },
    enablePlugin(idx, enabled) {
      basicapi
        .enablePlugin({ idx, enabled })
        .then(data => {
          console.log(data);
        })
        .catch(err => console.log(err));
    }
  },
  async created() {
    await this.refresh();
  }
};
</script>

<style scoped>
.plugin-status-wrapper {
  margin-bottom: 1em;
}

.working {
  background-color: #f0f9eb;
}

.stopped {
  background-color: #f4f4f5;
}

.link {
  color: #795548;
}
</style>
