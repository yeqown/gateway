<template>
  <div class="wrapper" style="line-height:40px" v-if="rule">
    <el-row>
      <el-col :span="6">
        <span class="gray" style="margin-right:0.5em">ID</span>
        {{rule.id}}
      </el-col>
      <el-col :span="5" :offset="0" style="text-align:left">
        <span class="gray" style="margin-right:0.5em">前缀匹配</span>
        {{rule.prefix}}
      </el-col>
      <el-col :span="5" :offset="0" style="text-align:left">
        <span class="gray" style="margin-right:0.5em">服务组别名</span>
        <router-link
          :to="{path: `reverse_server/${rule.server_name}`, params: {group: rule.server_name}}"
        >{{rule.server_name}}</router-link>
      </el-col>
      <el-col :span="4" :offset="0" style="text-align:left">
        <span class="gray" style="margin-right:0.5em">清除前缀
          <el-popover placement="top-start" title="关于清除前缀" width="200" trigger="hover">
            <i slot="reference" class="el-icon-info"></i>
            <p>
              原请求地址
              <span style="color:#409EFF">{{rule.prefix}}/path1</span>，不开启清除前缀为
              <span style="color:#409EFF">{{rule.prefix}}/path1</span>，开启则收到的请求为
              <span style="color:#409EFF">/path1</span>
            </p>
          </el-popover>
        </span>
        <el-switch v-model="rule.need_strip_prefix" @change="hdlEditServerRule" size="small"></el-switch>
      </el-col>
      <el-col :span="2" :offset="2">
        <el-button
          type="primary"
          size="small"
          @click="dialogFormVisible = true"
          circle
          icon="el-icon-edit"
        ></el-button>
        <el-button type="danger" size="small" @click="delServerRule" circle icon="el-icon-delete"></el-button>
      </el-col>
    </el-row>
    <!-- edit dialog -->
    <el-dialog :title="`修改规则: ${rule.id}`" :visible.sync="dialogFormVisible" width="400px">
      <el-form ref="refNewForm" :model="rule" label-position="left">
        <el-form-item label="前缀" label-width="100px">
          <el-input v-model="rule.prefix" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="服务器组别" label-width="100px">
          <el-select v-model="rule.server_name">
            <el-option
              v-for="(opt, idx) in reverseGroups"
              :value="opt.name"
              :label="`${opt.name}:${opt.count}`"
              :key="idx"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="取消前缀" label-width="100px">
          <el-switch v-model="rule.need_strip_prefix"></el-switch>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false" size="small">取 消</el-button>
        <el-button type="primary" @click="hdlEditServerRule" size="small">确 定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { proxyapi } from "@/apis";
export default {
  name: "ServerRuleDetail",
  data() {
    return {
      dialogFormVisible: false
    };
  },
  props: {
    rule: Object,
    refresh: Function,
    reverseGroups: Array
  },
  methods: {
    async hdlEditServerRule() {
      let { id, server_name, need_strip_prefix, prefix } = this.rule;
      await proxyapi
        .editServerRule({ id, server_name, need_strip_prefix, prefix })
        .then(data => this.$message.success(data.message))
        .catch(err => console.error(err));
      this.dialogFormVisible = false;
      this.refresh();
    },
    async delServerRule() {
      let { id } = this.rule;
      await proxyapi
        .delServerRule({ id })
        .then(data => this.$message.success(data.message))
        .catch(err => console.error(err));
      this.refresh();
    },
    async editServerRule() {
      this.refresh();
    }
  },
  created() {
    // console.log(this.rule);
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