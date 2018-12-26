<template>
  <div>
    <template v-if="current">
      <router-link to="path">URI</router-link>
      <br>
      <router-link to="server">Server</router-link>
      <br>
      <router-link to="reverse_server">Reverse Server</router-link>
      <br>
    </template>
    <transition name="fade">
      <router-view></router-view>
    </transition>
  </div>
</template>
<script>
export default {
  name: "PluginProxy",
  data() {
    return {
      current: false
    };
  },
  watch: {
    $route: function(newVal, oldVal) {
      // console.log("config router", newVal)
      if (newVal.path === "/configs/plugin/proxy") {
        this.current = true;
      } else {
        this.current = false;
      }
    }
  }
};
</script>

<style scoped>
</style>
