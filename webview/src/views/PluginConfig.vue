<template>
  <div>
    <transition name="fade">
      <router-view/>
    </transition>
  </div>
</template>

<script>
import { basicapi } from "@/apis";
export default {
  name: "Config"
};
</script>

<style scoped>
</style>
