<template>
  <div>
    <el-row style="margin-bottom: 1em; text-align:left">
      <el-col :span="4">组名：{{groupName}}</el-col>
      <el-col :span="4">实例数：{{total}}</el-col>
      <el-col :span="2" :offset="14">
        <el-button
          type="info"
          size="small"
          icon="el-icon-edit"
          circle
          @click="dialogFormVisible2 = true"
        />
        <el-button
          size="small"
          type="success"
          @click="openNewReverseSrvDialog"
          icon="el-icon-plus"
          circle
        />
      </el-col>
    </el-row>
    <template v-if="servers && servers.length">
      <el-row v-for="(s, idx) in servers" :key="idx" class="reverse-srv-instance">
        <el-col :span="6">Id：{{s.id}}</el-col>
        <el-col :span="3">服务名：{{s.name}}</el-col>
        <el-col :span="2">权重：{{s.weight}}</el-col>
        <el-col :span="6" :offset="0">服务地址：{{s.addr}}</el-col>
        <el-col :span="3" :offset="4">
          <el-button
            type="primary"
            icon="el-icon-info"
            size="small"
            circle
            @click="console.log(s.id)"
          ></el-button>
          <el-button
            type="success"
            icon="el-icon-edit"
            size="small"
            circle
            @click="hdlEditReverseServer(s)"
          ></el-button>
          <el-button
            type="danger"
            icon="el-icon-delete"
            size="small"
            circle
            @click="hdlDelReverseServer(s.id)"
          ></el-button>
        </el-col>
      </el-row>
    </template>
    <template v-else>
      <h2 style="color:gray">
        <i class="el-icon-info"></i> 暂无配置
      </h2>
    </template>
    <el-row>
      <el-pagination
        @current-change="handleCurrentChange"
        :current-page="curPage"
        :page-size="8"
        layout="total, prev, pager, next, jumper"
        :total="total"
      ></el-pagination>
    </el-row>

    <!-- new reverse server config dialog -->
    <el-dialog :title="!form.id?'新增实例':'更新实例配置'" :visible.sync="dialogFormVisible" width="500px">
      <el-form ref="refForm" :model="form" label-position="left">
        <el-form-item label="组别" label-width="40px">
          <el-input v-model="form.group" disabled></el-input>
        </el-form-item>
        <el-form-item label="名字" label-width="40px">
          <el-input v-model="form.name" autocomplete="off" placeholder="api.dev or api1"></el-input>
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
        <el-button type="primary" @click="hdlNewReverseSrv" size="small">确 定</el-button>
      </div>
    </el-dialog>

    <!-- new reverse server config dialog -->
    <el-dialog title="更新组名" :visible.sync="dialogFormVisible2" width="400px">
      <el-alert
        title="更新组名并不会同步更新其他相关配置，请确保此次操作在您的预期之内!"
        type="error"
        style="text-align:left;margin-bottom:1em"
      ></el-alert>
      <el-form ref="refForm" label-position="left">
        <el-form-item label="旧组名" label-width="80px">
          <el-input v-model="groupName" disabled></el-input>
        </el-form-item>
        <el-form-item label="新组名" label-width="80px">
          <el-input v-model="newGroupName" placeholder="groupName"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible2 = false" size="small">取 消</el-button>
        <el-button type="primary" @click="hdlRenameReverseSrvGroupName" size="small">确 定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { proxyapi } from "@/apis";
export default {
  name: "ReverseServerGroup",
  data() {
    return {
      curPage: 1,
      servers: [],
      total: 0,
      groupName: "",
      newGroupName: "",
      dialogFormVisible: false,
      dialogFormVisible2: false,
      form: {
        id: "",
        name: "",
        group: "",
        addr: "",
        weight: 5
      }
    };
  },
  methods: {
    handleCurrentChange(page) {
      this.curPage = page;
    },
    hdlRenameReverseSrvGroupName() {
      this.dialogFormVisible2 = false;
      proxyapi
        .renameReverseSrvGroupName({
          group: this.groupName,
          newname: this.newGroupName
        })
        .then(data => this.$message.success(data.message))
        .catch(err => this.$message.error(err.message));
      console.log(this.newGroupName);
    },
    hdlEditReverseServer(s) {
      let { id, name, group, addr, weight } = s;
      this.form = { id, name, group, addr, weight };
      this.dialogFormVisible = true;
    },
    async hdlDelReverseServer(id) {
      await proxyapi
        .delReverseSrv({ id: id })
        .then(data => this.$message.success(data.message))
        .catch(err => this.$message.error(err.message));
      this.getReverseSrvGroupDetail();
    },
    openNewReverseSrvDialog() {
      this.dialogFormVisible = true;
      this.form = { name: "", group: this.groupName, addr: "", weight: 5 };
    },
    getReverseSrvGroupDetail() {
      let { group } = this.$route.params;
      this.groupName = group;
      proxyapi
        .getReverseSrvGroupDetail({ group })
        .then(data => {
          this.servers = data.group;
          this.total = data.total;
          this.$message.success(data.message);
        })
        .catch(err => this.$message.error(err.message));
    },
    async hdlNewReverseSrv() {
      this.dialogFormVisible = false;
      let { id, group, name, addr, weight } = this.form;

      if (id) {
        // update
        await proxyapi
          .editReverseSrv({ id, name, group, addr, weight })
          .then(data => this.$message.success(data.message))
          .catch(err => this.$message.error(err.message));
      } else {
        // new
        await proxyapi
          .newReverseSrv({ group, name, addr, weight })
          .then(data => this.$message.success(data.message))
          .catch(err => this.$message.error(err.message));
      }
      this.getReverseSrvGroupDetail();
    }
  },
  created() {
    // console.log(this.$route.params);
    this.getReverseSrvGroupDetail();
  }
};
</script>

<style scoped>
.reverse-srv-instance {
  background-color: #fbfbfb;
  padding: 0.5em;
  margin-bottom: 1em;
  text-align: left;
  line-height: 60px;
}
</style>
