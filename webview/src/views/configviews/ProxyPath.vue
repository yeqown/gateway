<template>
  <div>
    <el-row style="margin-bottom: 1em">
      <el-col :span="2">
        <span>组数: {{total}}</span>
      </el-col>
      <el-col :span="1" :offset="21">
        <el-button
          type="success"
          size="small"
          icon="el-icon-plus"
          circle
          @click="dialogFormVisible = true"
        />
      </el-col>
    </el-row>
  </div>
</template>
<script>
import { proxyapi } from "@/apis";
export default {
  name: "PathRule",
  data() {
    return {
      // groups: [],
      total: 0,
      dialogFormVisible: false,
      // form: {
      //   name: "",
      //   addr: "",
      //   weight: "",
      //   group: ""
      // }
    };
  },
  methods: {
    refresh() {
      console.log('PorxyPath refresh')
    }
  },
  created() {
    this.refresh();
  }
};
</script>

<style scoped>
.group-wrapper {
  /* font-size: 1.2em; */
  background-color: #fbfbfb;
  text-align: left;
  padding: 1em;
  margin-bottom: 1em;
}
.link {
  color: green;
}
</style>
