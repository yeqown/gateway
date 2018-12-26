<template>
  <div>
    <template v-if="pathRules && pathRules.length">
      <el-row style="margin-bottom: 1em">
        <!-- <el-col :span="2">
          <span>组数: {{total}}</span>
        </el-col>-->
        <el-col :span="1" :offset="23">
          <el-button
            type="success"
            size="small"
            icon="el-icon-plus"
            circle
            @click="$router.push({path: 'pathrule/new', params: {id: 'new'}})"
          />
        </el-col>
      </el-row>
      <!-- table row -->
      <el-row style="margin-bottom:2em;">
        <el-table :data="pathRules" style="width: 100%">
          <el-table-column label="请求方法" width="100">
            <template slot-scope="scope">
              <el-tag :type="judgeMethoTagType(scope.row.method)" size="mini">{{ scope.row.method }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="URI" width="200">
            <template slot-scope="scope">
              <span style="margin-left: 10px">{{ scope.row.path }}</span>
            </template>
          </el-table-column>
          <el-table-column label="重定向URI" width="200">
            <template slot-scope="scope">
              <span>{{ scope.row.rewrite_path || '组合请求为空' }}</span>
            </template>
          </el-table-column>
          <el-table-column label="组别" width="200">
            <template slot-scope="scope">
              <span>
                <router-link
                  class="link"
                  :to="{path: `reverse_server/${scope.row.server_name}`, params: {group: scope.row.server_name}}"
                >{{ scope.row.server_name }}</router-link>
              </span>
            </template>
          </el-table-column>
          <el-table-column label="组合" width="200">
            <template slot-scope="scope">
              <el-switch :value="scope.row.need_combine" disabled></el-switch>
            </template>
          </el-table-column>
          <el-table-column label="操作">
            <template slot-scope="scope">
              <span id="more-rule-detail">
                <i class="el-icon-more" @click="hdlClickDetail(scope.row.id)"></i>
              </span>
              <span id="del-rule-btn">
                <i class="el-icon-delete" @click="hdlDelPathRule(scope.row.id)"></i>
              </span>
            </template>
          </el-table-column>
        </el-table>
      </el-row>
      <!-- pagination -->
      <el-row>
        <el-pagination
          @current-change="handleCurrentChange"
          :current-page="curPage"
          :page-size="8"
          layout="total, prev, pager, next, jumper"
          :total="total"
        ></el-pagination>
      </el-row>
    </template>
    <template v-else style="min-height:400px;">
      <h2 style="color:gray">
        <i class="el-icon-info"></i> 暂无配置
      </h2>
    </template>
  </div>
</template>
<script>
import { proxyapi } from "@/apis";
export default {
  name: "PathRule",
  data() {
    return {
      curPage: 1,
      pathRules: [],
      total: 0
    };
  },
  methods: {
    handleCurrentChange(page) {
      this.curPage = page;
    },
    hdlClickDetail(id) {
      this.$router.push({ path: `pathrule/${id}`, params: { id: id } });
    },
    async hdlDelPathRule(id) {
      await proxyapi
        .delPathRule({ id })
        .then(data => this.$message.success(data.message))
        .catch(err => this.$message.error(err.message));
      this.refresh();
    },
    getPathRules() {
      let limit = 8;
      let offset = (this.curPage - 1) * limit;
      proxyapi
        .getPathRules({ limit, offset })
        .then(data => {
          this.pathRules = data.rules;
          this.total = data.total;
          this.$message.success(data.message);
        })
        .catch(err => this.$message.error(err.message));
    },
    judgeMethoTagType(method) {
      switch (method.toLowerCase()) {
        case "get":
          return "success";
        case "post":
          return "error";
        default:
          return "success";
      }
    },
    refresh() {
      this.getPathRules();
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
#more-rule-detail {
  margin-right: 1em;
  cursor: pointer;
}
#more-rule-detail:hover {
  color: green;
}
#del-rule-btn {
  margin-right: 1em;
  cursor: pointer;
}
#del-rule-btn:hover {
  color: red;
}
</style>
