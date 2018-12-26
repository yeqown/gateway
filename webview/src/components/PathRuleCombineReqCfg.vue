<template>
  <div class="wrapper" v-if="cfg">
    <el-row>
      <el-col :span="1" :offset="23">
        <i class="el-icon-delete del-combine-req-cfg" @click="delCombineReqCfg"></i>
      </el-col>
    </el-row>
    <el-form :model="cfg" label-width="100px">
      <el-form-item v-if="cfg.id" prop="id" label="id">
        <el-tag type="info" size="small">{{cfg.id}}</el-tag>
      </el-form-item>
      <el-form-item prop="path" label="请求路径">
        <el-input v-model="cfg.path" size="small" placeholder="/gateway/uri2"></el-input>
      </el-form-item>
      <el-form-item prop="field" label="响应结果字段">
        <el-input v-model="cfg.field" size="small" placeholder="repsonse field"></el-input>
      </el-form-item>
      <el-form-item prop="server_name" label="服务器组别名">
        <el-input v-model="cfg.server_name" size="small" placeholder="groupName"></el-input>
      </el-form-item>
      <el-form-item prop="method" label="请求方法" size="small">
        <el-select v-model="cfg.method" placeholder="请选择">
          <el-option value="GET" label="GET"></el-option>
          <el-option value="POST" label="POST"></el-option>
          <el-option value="DELETE" label="DELETE"></el-option>
          <el-option value="PUT" label="PUT"></el-option>
        </el-select>
        <!-- <el-input v-model="cfg.method" size="small"></el-input> -->
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
export default {
  name: "CombineReqCfg",
  props: {
    cfg: Object,
    idx: Number,
    delFunc: Function,
    change: Function
  },
  methods: {
    delCombineReqCfg() {
      this.delFunc(this.idx);
    }
  },
  watch: {
    cfg: {
      handler: function(newVal, oldVal) {
        console.log("changed:", newVal, oldVal);
        this.change(this.idx, this.cfg);
      },
      deep: true
    }
  }
  // created() {
  //   console.log(this.cfg, this.idx);
  // }
};
</script>

<style scoped>
.wrapper {
  background-color: #f3f3f5;
  padding: 0.5em;
  border-radius: 5px;
}
.del-combine-req-cfg {
  cursor: pointer;
}
.del-combine-req-cfg:hover {
  color: red;
}
</style>

