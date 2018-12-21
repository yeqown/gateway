<template>
  <div class="wrapper">
    <el-row>
      <el-col :span="2">
        <span class="gray">ID</span>
        {{cache.id}}
      </el-col>
      <el-col :span="4" :offset="2" style="text-align:left">
        <span class="gray">REGEXP </span>
        {{cache.regular}}
      </el-col>
      <el-col :span="8" :offset="0" style="text-align:left">
        <span class="gray">ENABLED </span>

        <el-switch :value="true" disabled></el-switch>
      </el-col>
      <el-col :span="4" :offset="4">
        <el-button type="success" size="small">修改</el-button>
        <el-button type="danger" size="small">删除</el-button>
      </el-col>
    </el-row>
  </div>
</template>

<script>
export default {
  name: "Cache",
  props: {
    cache: Object
  },
  created() {
    console.log(cache);
  }
};
</script>

<style scoped>
.wrapper {
  margin-bottom: 1em;
  padding: 0.5em;
  background-color: #f4f7f7;
}
.gray {
  color: gray;
}
</style>