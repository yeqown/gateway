<template>
  <div>
    <el-row>
      <el-col :span="12" :offset="6">
        <el-form
          :model="form"
          :rules="rules"
          ref="form"
          label-width="100px"
          label-position="right"
          class="form"
        >
          <el-form-item label="网关端口号" prop="port">
            <el-input-number v-model="form.port"/>
          </el-form-item>
          <el-form-item label="域名配置" prop="baseURL">
            <el-input v-model="form.baseURL"/>
          </el-form-item>
          <el-form-item label="日志文件夹" prop="logpath">
            <el-input v-model="form.logpath"/>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="updateBasicConfig">确认</el-button>
            <el-button type="danger" @click="reset">重置</el-button>
          </el-form-item>
        </el-form>
      </el-col>
    </el-row>
  </div>
</template>
<script>
import { basicapi, resetBaseURL, baseURL } from "@/apis/index";
export default {
  name: "BaseConfig",
  data() {
    return {
      // confirmDisabled: true,
      form: {
        logpath: "",
        port: "",
        baseURL: "http://localhost:8989"
      },
      rules: {
        logpath: [
          {
            required: true,
            message: "请输入日志文件夹地址",
            trigger: "change"
          },
          { min: 3, max: 20, message: "长度在1到20个字符", trigger: "blur" }
        ],
        port: [
          { required: true, message: "请输入网关端口号", trigger: "change" },
          {
            type: "number",
            min: 100,
            max: 65555,
            message: "在9000-65555之间",
            trigger: "blur"
          }
        ],
        baseURL: [{ required: true, message: "请输入服务端域名" }]
      }
    };
  },
  methods: {
    reset() {
      // resetBaseURL({ baseURL: "http://localhost:7777" });
    },
    updateBasicConfig() {
      this.$refs.form.validate((isValid, obj) => {
        if (isValid) {
          basicapi
            .putGlobalConfig(this.form)
            .then(data => {
              this.$message.success(data.message);
            }).catch(err => {
              this.$message.error(err.message);
            })
          return;
        }
        // invalid ...
        let errmsg = "";
        Object.keys(obj).map(key => {
          let errmsgInner = `${key}:`;
          obj[key].forEach(item => {
            // console.log(item)
            errmsgInner += item.message;
          });
          errmsg += errmsgInner + ";";
        });
        this.$message.error(errmsg);
      });
    }
  },
  created() {
    // const loading = this.$loading({
    //   lock: true,
    //   text: "请求服务中...",
    //   spinner: "el-icon-loading",
    //   background: "rgba(0, 0, 0, 0.7)"
    // });
    this.form.baseURL = baseURL
    // console.log(this.form)
    basicapi
      .getGlobalConfig()
      .then(data => {
        this.form.port = data.port
        this.form.logpath = data.logpath 
        // console.log(data, this.form)
      })
      .catch(err => {
        this.$notify.error({
          title: "错误提示",
          message: err.message
        });
      });
    // loading.close();
  }
};
</script>

<style scoped>
</style>
