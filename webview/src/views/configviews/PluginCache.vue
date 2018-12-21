<template>
  <div>
    <el-row>
      <el-col :span="2" :offset="0">
        <el-button size="small" type="primary">新增</el-button>
      </el-col>
      <el-col :span="22" :offset="0">Here here nothing</el-col>
    </el-row>
    <el-row style="margin-top:1em;">
      <el-col>
        <i-cache v-for="(nocache, idx) in nocaches" :key="idx" :cache="nocache"></i-cache>
      </el-col>
    </el-row>
  </div>
</template>
<script>
import Cache from "@/components/Cache.vue";
import { cacheapi } from "@/apis";
export default {
  name: "BaseConfig",
  data() {
    return {
      limit: 10,
      offset: 0,
      nocaches: []
    };
  },
  components: {
    iCache: Cache
  },
  created() {
    cacheapi
      .getCacheRules()
      .then(data => {
        // console.log(data);
        this.nocaches = data.rules
      })
      .catch(err => {
        console.error(err);
      });
  }
};
</script>


<style scoped>
</style>
