<template>
  <div class="wrapper" style="line-height:40px" v-if="nocache">
    <el-row>
      <el-col :span="6">
        <span class="gray">ID</span>
        {{nocache.id}}
      </el-col>
      <el-col :span="10" :offset="0" style="text-align:left">
        <span class="gray">正则表达式：</span>
        {{nocache.regular}}
      </el-col>
      <el-col :span="4" :offset="0" style="text-align:left">
        <span class="gray">是否启用：</span>
        <el-switch v-model="nocache.enabled" @change="hdlEditNocacheRule"></el-switch>
      </el-col>
      <el-col :span="2" :offset="2">
        <el-button
          type="primary"
          size="small"
          @click="dialogFormVisible = true"
          circle
          icon="el-icon-edit"
        ></el-button>
        <el-button type="danger" size="small" @click="delNocacheRule" circle icon="el-icon-delete"></el-button>
      </el-col>
    </el-row>
    <!-- edit dialog -->
    <el-dialog :title="`修改规则: ${nocache.id}`" :visible.sync="dialogFormVisible" width="400px">
      <el-form ref="refNewForm" :model="nocache" label-position="left">
        <el-form-item label="正则表达式" label-width="100px">
          <el-input v-model="nocache.regular" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="是否启用" label-width="100px">
          <el-switch v-model="nocache.enabled"></el-switch>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false" size="small">取 消</el-button>
        <el-button type="primary" @click="hdlEditNocacheRule" size="small">确 定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { cacheapi } from "@/apis";
export default {
  name: "Cache",
  data() {
    return {
      dialogFormVisible: false
    };
  },
  props: {
    nocache: Object,
    refresh: Function
  },
  methods: {
    async hdlEditNocacheRule() {
      let { id, regular, enabled } = this.nocache;
      // console.log(id, regular, enabled);
      await cacheapi
        .editCacheRule({ id, regular, enabled })
        .then(data => this.$message.success(data.message))
        .catch(err => console.error(err));
      this.dialogFormVisible = false;
      this.refresh();
    },
    async delNocacheRule() {
      let { id } = this.nocache;
      // console.log(this.nocache, id)
      await cacheapi
        .delCacheRule({ id })
        .then(data => this.$message.success(data.message))
        .catch(err => console.error(err));
      this.refresh();
    },
    async editNocacheRule() {
      this.refresh();
    }
  },
  created() {
    // console.log(this.nocache);
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