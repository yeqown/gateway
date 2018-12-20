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
            <el-input-number v-model="form.port"></el-input-number>
          </el-form-item>
          <el-form-item label="日志文件夹" prop="logpath">
            <el-input v-model="form.logpath"></el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" :disabled="confirmDisabled">确认</el-button>
            <el-button type="danger">重置</el-button>
          </el-form-item>
        </el-form>
      </el-col>
    </el-row>
  </div>
</template>
<script>
import { basicapi } from "@/apis/index";
export default {
  name: "BaseConfig",
  data() {
    return {
      confirmDisabled: true,
      form: {
        logpath: "",
        port: ""
      },
      rules: {
        logpath: [
          { required: true, message: "请输入日志文件夹地址", trigger: "blur" },
          { min: 3, max: 20, message: "长度在1到20个字符", trigger: "blur" }
        ],
        port: [
          { required: true, message: "请输入网关端口号", trigger: "change" },
          {
            min: 9000,
            max: 65555,
            message: "在9000-65555之间",
            trigger: "blur"
          }
        ]
      }
    };
  },
  created() {
    basicapi
      .getGlobalConfig()
      .then(data => {
        console.log(data);
      })
      .catch(err => {
        this.$notify.error({
          title: "提示",
          message: err.message
        });
      });
  }
};
</script>

<style scoped>
</style>
