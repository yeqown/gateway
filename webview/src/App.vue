<template>
  <div id="app">
    <el-container>
      <!-- header -->
      <el-header height="100">
        <!-- top bar -->
        <el-menu
          :default-active="activeIndex"
          class="el-menu-demo"
          mode="horizontal"
          @select="handleSelect"
        >
          <el-menu-item index="0">
            <img src="@/assets/logo.png" alt="logo" width="50px">
          </el-menu-item>
          <el-menu-item index="gate.global">概览</el-menu-item>
          <el-menu-item index="gate.basic">基本配置</el-menu-item>
          <el-submenu index="gate.plugin">
            <template slot="title">插件配置</template>
            <el-menu-item index="gate.plugin.cache">缓存插件</el-menu-item>
            <el-menu-item index="gate.plugin.ratelimit" disabled>限流插件</el-menu-item>
            <el-submenu index="gate.plugin.proxy">
              <template slot="title">代理插件</template>
              <el-menu-item index="gate.plugin.proxy.path">路径代理</el-menu-item>
              <el-menu-item index="gate.plugin.proxy.server">反向服务代理</el-menu-item>
              <el-menu-item index="gate.plugin.proxy.reverseServer">服务实例组</el-menu-item>
            </el-submenu>
          </el-submenu>
          <el-menu-item index="gate.plugins">网关插件管理</el-menu-item>
          <el-menu-item index="gate.repo">
            <a href="https://github.com/yeqown/gateway" target="_blank">Github</a>
          </el-menu-item>
        </el-menu>
        <!-- bread cumb -->
        <el-breadcrumb separator="/" class="breadcumb">
          <el-breadcrumb-item
            v-for="(breadcrumb,idx) in breadcrumbs"
            :key="idx"
            :to="breadcrumb.to"
          >{{breadcrumb.name}}</el-breadcrumb-item>
          <!-- <el-breadcrumb-item><a href="/">活动管理</a></el-breadcrumb-item> -->
        </el-breadcrumb>
      </el-header>

      <!-- body -->
      <el-main>
        <router-view/>
      </el-main>
    </el-container>

    <!-- <el-footer>
      <p>
        Copyright@yeqown 2018
        yeqwon@gmail.com
      </p>
    </el-footer> -->
  </div>
</template>

<script>
import {saveBreadcrumbs, loadBreadcrumbs} from './config'
export default {
  name: "app",
  data() {
    return {
      activeIndex: "1",
      breadcrumbs: [
        {
          name: '导航页',
          to: '/dashbord'
        }
      ]
    };
  },
  methods: {
    handleSelect(key, keyPath) {
      console.log("menu selectd: ", key);
      this.activeIndex = key;
      switch (key) {
        case "gate.global":
          this.$router.push("/dashbord");
          break;
        case "gate.basic":
          this.$router.push("/configs/basic");
          break;
        case "gate.plugin.cache":
          this.$router.push("/configs/plugin/cache")
          break;
        case "gate.plugins":
          this.$router.push("/plugins")
      }
    }
  },
  watch: {
    $route(newVal, oldVal) {
      if (newVal.params.breadcrumbs) {
        this.breadcrumbs = newVal.params.breadcrumbs
        saveBreadcrumbs(this.breadcrumbs)
      }
    }
  },
  created () {
    let breadcrumbs = loadBreadcrumbs()
    // console.log(breadcrumbs)
    if (breadcrumbs) {
      this.breadcrumbs = breadcrumbs
    }

  }
};
</script>

<style>
#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  width: 80%;
  margin: 0 auto;
  min-height: 600px;
  /* margin-top: 60px; */
}

.breadcumb {
  margin: 1em;
}
</style>
