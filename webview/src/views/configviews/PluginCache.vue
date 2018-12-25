<template>
  <div style="text-align:left">
    <el-row>
      <el-col :span="1" :offset="23">
        <el-button
          size="small"
          type="success"
          @click="dialogFormVisible=true"
          icon="el-icon-plus"
          circle
        />
      </el-col>
      <!-- <el-col :span="22" :offset="0">规则数总计：{{total}}</el-col> -->
    </el-row>
    <el-row style="margin-top:1em;">
      <el-col>
        <i-cache
          v-for="(nocache, idx) in nocaches"
          :key="idx"
          :nocache="nocache"
          :refresh="refresh"
        ></i-cache>
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
    <el-dialog title="发布规则" :visible.sync="dialogFormVisible" width="400px">
      <el-form ref="refNewForm" :model="newForm" label-position="left">
        <el-form-item label="正则表达式" label-width="100px">
          <el-input v-model="newForm.regular" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="是否启用" label-width="100px">
          <el-switch v-model="newForm.enabled"></el-switch>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false" size="small">取 消</el-button>
        <el-button type="primary" @click="hdlNewNocacheRule" size="small">确 定</el-button>
      </div>
    </el-dialog>
  </div>
</template>
<script>
import Cache from "@/components/Cache";
import { cacheapi } from "@/apis";
export default {
  name: "PluginCache",
  data() {
    return {
      curPage: 1,
      nocaches: [],
      total: 0,
      dialogFormVisible: false,
      newForm: {
        regular: "",
        enabled: true
      }
    };
  },
  components: {
    iCache: Cache
  },
  methods: {
    handleCurrentChange(page) {
      this.curPage = page;
      this.refresh();
    },
    async hdlNewNocacheRule() {
      let { enabled, regular } = this.newForm;
      this.dialogFormVisible = false;
      await cacheapi
        .newCacheRule({ enabled, regular })
        .then(data => this.$message.success(data.message))
        .catch(err => {
          console.error(err);
          return;
        });
      this.resetNewForm();
      this.refresh();
    },
    getCacheRules(page) {
      let limit = 8;
      let offset = (this.curPage - 1) * limit;
      cacheapi
        .getCacheRules({ limit, offset })
        .then(data => {
          // console.log(data);
          this.nocaches = data.rules;
          this.total = data.total;
          this.$message.success(data.message);
        })
        .catch(err => {
          console.error(err);
        });
    },
    refresh() {
      this.getCacheRules(this.curPage);
    },
    resetNewForm() {
      this.newForm = { regular: "", enabled: true };
    }
  },
  async created() {
    await this.getCacheRules(this.curPage);
  }
};
</script>


<style scoped>
</style>
