<template>
  <div>
    <el-row style="margin-bottom: 1em">
      <el-col :span="2">
        <span>组数: {{total}}</span>
      </el-col>
      <el-col :span="1" :offset="21">
        <el-button
          type="primary"
          size="small"
          icon="el-icon-plus"
          circle
          @click="dialogFormVisible = true"
        />
      </el-col>
    </el-row>
    <template v-if="groups && groups.length">
      <el-row v-for="(g, idx) in groups" :key="idx" class="group-wrapper">
        <el-col :span="1">
          <i class="el-icon-info" style="color:#03a9f4; font-size:1.5em; line-height: 40px"/>
        </el-col>
        <el-col :span="4">组别：{{g.name}}</el-col>
        <el-col :span="4">实例数：{{g.count}}</el-col>
        <el-col :span="1" :offset="13" style="line-height: 40px">
          <router-link
            :to="{path:`reverse_server/${g.name}`, params: {group: g.name}}"
            class="link"
          >
            <i class="el-icon-setting" style="color:gray;font-size:1.5em;"/>
            <!-- 详细 -->
          </router-link>
        </el-col>
        <el-col :span="1">
          <el-button
            icon="el-icon-delete"
            type="danger"
            size="small"
            circle
            @click="hdlDelReverseSrvGroup(g.name)"
          />
        </el-col>
      </el-row>
    </template>
    <template v-else>
      <h2 style="color:gray">
        <i class="el-icon-info"></i> 暂无配置
      </h2>
    </template>
    <!-- new reverse server config dialog -->
    <el-dialog :title="!form.id?'新增实例':'更新实例配置'" :visible.sync="dialogFormVisible" width="500px">
      <el-form ref="refForm" :model="form" label-position="left">
        <el-form-item label="组别" label-width="40px">
          <el-input v-model="form.group" placeholder="groupName1"></el-input>
        </el-form-item>
        <el-form-item label="名字" label-width="40px">
          <el-input v-model="form.name" autocomplete="off" placeholder="api.dev OR api1"></el-input>
        </el-form-item>
        <el-form-item label="地址" label-width="40px">
          <el-input v-model="form.addr" autocomplete="off" placeholder="http://api.example.com"></el-input>
        </el-form-item>
        <el-form-item label="权重" label-width="40px">
          <el-input v-model="form.weight" placeholder="5"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false" size="small">取 消</el-button>
        <el-button type="primary" @click="hdlNewReverseServer" size="small">确 定</el-button>
      </div>
    </el-dialog>
  </div>
</template>
<script>
import { proxyapi } from "@/apis";
export default {
  name: "ReverseServer",
  data() {
    return {
      groups: [],
      total: 0,
      dialogFormVisible: false,
      form: {
        name: "",
        addr: "",
        weight: "",
        group: ""
      }
    };
  },
  methods: {
    async hdlNewReverseServer() {
      let { name, addr, group, weight } = this.form;
      this.dialogFormVisible = false;
      this.resetNewForm()
      await proxyapi
        .newReverseSrv({ name, addr, group, weight })
        .then(data => this.$message.success(data.message))
        .catch(err => this.$message.error(err.message));
      this.getReverseSrvGroup();
    },
    async hdlDelReverseSrvGroup(groupName) {
      await proxyapi
        .delReverseSrvGroup({ group: groupName })
        .then(data => this.$message.success(data.message))
        .catch(err => this.$message.error(err.message));
      this.getReverseSrvGroup();
    },
    getReverseSrvGroup() {
      // reinit
      this.groups = [];
      this.total = 0;
      proxyapi
        .getReverseSrvGroup()
        .then(data => {
          let g = data.groups;
          let count = 0;
          Object.keys(g).map(k => {
            count++;
            this.groups.push({ name: k, count: g[k] });
          });
          this.total = count;
          this.$message.success(data.message);
        })
        .catch(err => this.$notify.error(err.message));
    },
    resetNewForm() {
      this.form = { name: "", addr: "", group: "", weight: "" };
    },
    refresh() {
      this.getReverseSrvGroup();
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
