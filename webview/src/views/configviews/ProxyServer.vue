<template>
  <div style="text-align:left">
    <el-row>
      <el-col :span="1" :offset="23">
        <el-button
          size="small"
          type="primary"
          @click="dialogFormVisible=true"
          icon="el-icon-plus"
          circle
        />
      </el-col>
      <!-- <el-col :span="22" :offset="0">规则数总计：{{total}}</el-col> -->
    </el-row>
    <el-row style="margin-top:1em;">
      <el-col>
        <i-rule
          v-for="(rule, idx) in serverRules"
          :key="idx"
          :rule="rule"
          :reverseGroups="reverseGroups"
          :refresh="refresh"
        ></i-rule>
      </el-col>
    </el-row>
    <el-row>
      <el-pagination
        @current-change="handleCurrentChange"
        :current-page="curPage"
        :page-size="8"
        layout="total, prev, pager, next, jumper"
        :total="total"
      ></el-pagination>
    </el-row>
    <!-- new dialog -->
    <el-dialog title="新增规则" :visible.sync="dialogFormVisible" width="400px">
      <el-form ref="refNewForm" :model="newForm" label-position="left">
        <el-form-item label="前缀" label-width="100px">
          <el-input v-model="newForm.prefix" autocomplete="off" placeholder="/prefix"></el-input>
        </el-form-item>
        <el-form-item label="服务器组别" label-width="100px">
          <el-select v-model="newForm.server_name">
            <el-option
              v-for="(opt, idx) in reverseGroups"
              :value="opt.name"
              :label="`${opt.name}:${opt.count}`"
              :key="idx"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="取消前缀" label-width="100px">
          <el-switch v-model="newForm.need_strip_prefix"></el-switch>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false" size="small">取 消</el-button>
        <el-button type="primary" @click="hdlNewServerRule" size="small">确 定</el-button>
      </div>
    </el-dialog>
  </div>
</template>
<script>
import ServerRuleDetail from "@/components/ServerRuleDetail";
import { proxyapi } from "@/apis";
export default {
  name: "ServerRule",
  data() {
    return {
      curPage: 1,
      serverRules: [],
      reverseGroups: [],
      total: 0,
      dialogFormVisible: false,
      newForm: {
        server_name: "",
        prefix: "",
        need_strip_prefix: true
      }
    };
  },
  components: {
    iRule: ServerRuleDetail
  },
  methods: {
    handleCurrentChange(page) {
      this.curPage = page;
      this.refresh();
    },
    async hdlNewServerRule() {
      let { server_name, need_strip_prefix, prefix } = this.newForm;
      this.dialogFormVisible = false;
      await proxyapi
        .newServerRule({ server_name, need_strip_prefix, prefix })
        .then(data => this.$message.success(data.message))
        .catch(err => {
          console.error(err);
          return;
        });
      this.resetNewForm();
      this.refresh();
    },
    getReverseGroups() {
      this.reverseGroups = [];
      proxyapi
        .getReverseSrvGroup()
        .then(data => {
          let g = data.groups;
          Object.keys(g).map(key => {
            this.reverseGroups.push({ name: key, count: g[key] });
          });
          this.$message.success(data.message);
        })
        .catch(err => this.$message.error(err.message));
    },
    getServerRules(page) {
      let limit = 8;
      let offset = (this.curPage - 1) * limit;
      proxyapi
        .getServerRules({ limit, offset })
        .then(data => {
          // console.log(data);
          this.serverRules = data.rules;
          this.total = data.total;
          this.$message.success(data.message);
        })
        .catch(err => {
          console.error(err);
        });
    },
    refresh() {
      this.getServerRules(this.curPage);
      this.getReverseGroups();
    },
    resetNewForm() {
      this.newForm = { server_name: "", need_strip_prefix: true, prefix: "" };
    }
  },
  async created() {
    // await this.getServerRules(this.curPage);
    this.refresh();
  }
};
</script>


<style scoped>
</style>
