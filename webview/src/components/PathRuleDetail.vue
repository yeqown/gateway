<template>
  <div v-if="pathRule" class="wrapper">
    <el-row style="width:60%;margin:0 auto;">
      <el-form :model="pathRule" ref="pathRule" label-width="100px">
        <el-form-item v-if="pathRule.id" prop="id" label="配置id">
          <el-tag type="info">{{pathRule.id}}</el-tag>
        </el-form-item>
        <el-form-item prop="path" label="请求路径">
          <el-input v-model="pathRule.path" size="small" placeholder="/gateway/uri2"></el-input>
        </el-form-item>
        <el-form-item prop="rewrite_path" label="重写路径">
          <el-input v-model="pathRule.rewrite_path" size="small" placeholder="/srv/uri1"></el-input>
        </el-form-item>
        <el-form-item prop="server_name" label="服务组别名">
          <el-input v-model="pathRule.server_name" size="small" placeholder="groupName1"></el-input>
        </el-form-item>
        <el-form-item prop="need_combine" label="组合请求">
          <el-switch v-model="pathRule.need_combine" size="small"></el-switch>
          <el-popover
            placement="top-start"
            title="提示"
            width="200"
            trigger="hover"
            style="margin-left: 1em; font-size: 1.2em;line-height:2em;"
          >
            <i class="el-icon-info" slot="reference"></i>
            <p>如果不开启则默认会URI直接代理，不会同时请求两个服务并组合结果</p>
          </el-popover>
        </el-form-item>
        <el-form-item prop="method" label="请求方法">
          <el-select v-model="pathRule.method" placeholder="请选择" size="small">
            <el-option value="GET" label="GET"></el-option>
            <el-option value="POST" label="POST"></el-option>
            <el-option value="DELETE" label="DELETE"></el-option>
            <el-option value="PUT" label="PUT"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item
          v-for="(cfg, idx) in pathRule.combine_req_cfgs"
          :label="'组合请求_' + idx"
          :key="cfg.id + idx"
          v-if="pathRule.need_combine && pathRule.combine_req_cfgs && pathRule.combine_req_cfgs.length"
        >
          <i-combine-req-cfg
            :cfg="pathRule.combine_req_cfgs[idx]"
            :idx="idx"
            :delFunc="removeCombineReqCfg"
            :change="updateCombineReqCfg"
          ></i-combine-req-cfg>
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            size="small"
            @click="hdlEditPathRule"
          >{{(typ === 'new')?'新增':'修改'}}</el-button>
          <el-button @click="addCombineReqCfg" size="small">新增组合配置</el-button>
        </el-form-item>
      </el-form>
    </el-row>
  </div>
</template>

<script>
import { proxyapi } from "@/apis";
import CombineReqCfg from "@/components/PathRuleCombineReqCfg";
export default {
  name: "PathRuleDetail",
  components: {
    iCombineReqCfg: CombineReqCfg
  },
  data() {
    return {
      typ: "edit",
      pathRule: {
        path: "",
        rewrite_path: "",
        method: "GET",
        server_name: "",
        combine_req_cfgs: [],
        need_combine: false,
        id: ""
      }
    };
  },
  methods: {
    getPathRuleDetail(id) {
      proxyapi
        .getPathRuleByID({ id })
        .then(data => {
          this.pathRule = data.rule;
          this.$message.success(data.message);
        })
        .catch(err => this.$message.error(err.message));
    },
    removeCombineReqCfg(idx) {
      // var index = this.pathRule.combine_req_cfgs.indexOf(item);
      // console.log("del combine-req-cfg idx", idx);
      if (idx !== -1) {
        this.pathRule.combine_req_cfgs.splice(idx, 1);
      }
    },
    updateCombineReqCfg(idx, cfg) {
      // var index = this.pathRule.combine_req_cfgs.indexOf(item);
      console.log("update combin-req-cfg idx", idx, cfg);
      if (idx !== -1) {
        this.pathRule.combine_req_cfgs[idx] = cfg;
      }
    },
    addCombineReqCfg() {
      if (!this.pathRule.need_combine) {
        return;
      }
      this.pathRule.combine_req_cfgs.push({
        server_name: "",
        path: "",
        field: "",
        method: "",
        id: ""
      });
    },
    hdlEditPathRule() {
      let {
        id,
        path,
        rewrite_path,
        method,
        server_name,
        combine_req_cfgs,
        need_combine
      } = this.pathRule;
      // new a path rule
      if (this.typ !== "edit") {
        proxyapi
          .newPathRule({
            path,
            rewrite_path,
            method,
            server_name,
            combine_req_cfgs,
            need_combine
          })
          .then(data => {
            this.$message.success(data.message);
          })
          .catch(err => this.$message.error(err.message));
        return;
      }
      // update a path rule
      proxyapi
        .editPathRule({
          id,
          path,
          rewrite_path,
          method,
          server_name,
          combine_req_cfgs,
          need_combine
        })
        .then(data => {
          this.$message.success(data.message);
        })
        .catch(err => this.$message.error(err.message));
      return;
    }
  },
  created() {
    // this.$on("del-combine-req-cfg", idx => this.removeCombineReqCfg(idx));
    let { id } = this.$route.params;
    if (id === "new") {
      this.typ = "new";
      return;
    }
    this.getPathRuleDetail(id);
  }
};
</script>

<style scoped>
.wrapper {
  text-align: left;
}
</style>
