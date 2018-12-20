<template>
  <div>
    <!-- <el-row :gutter="1" class="row" v-show="false">
      <el-col :span="8" key="gate.basic" :offset="3">
        <el-card :body-style="{ padding: '0px',cursor: 'pointer' }" shadow="hover">
          <img src="@/assets/gateway.jpg" class="image">
          <div style="padding: 10px;">
            <span class="card-title">
              <i class="el-icon-setting" style="color:green"/>
              网关基本配置
            </span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8" key="plugin.ratelimit" :offset="3">
        <el-card :body-style="{ padding: '0px',cursor: 'pointer' }" shadow="hover">
          <img src="@/assets/gateway.jpg" class="image">
          <div style="padding: 10px;">
            <span class="card-title">
              <i class="el-icon-d-caret" style="color:green"/>
              流量限制
            </span>
          </div>
        </el-card>
      </el-col>
    </el-row>
    <el-row :gutter="1" class="row" v-show="false">
      <el-col :span="8" key="plugin.proxy" :offset="3">
        <el-card :body-style="{ padding: '0px',cursor: 'pointer' }" shadow="hover">
          <img src="@/assets/gateway.jpg" class="image">
          <div style="padding: 10px;">
            <span class="card-title">
              <i class="el-icon-share" style="color:green"/>
              代理
            </span>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8" key="plugin.cache" :offset="3">
        <el-card :body-style="{ padding: '0px',cursor: 'pointer' }" shadow="hover">
          <img src="@/assets/gateway.jpg" class="image">
          <div style="padding: 10px;">
            <span class="card-title">
              <i class="el-icon-refresh" style="color:green"/>
              缓存
            </span>
          </div>
        </el-card>
      </el-col>
    </el-row> -->

    <el-row class="row text">
      <el-col :offset="0" :span="24">
        <h3>关于网关</h3>
        <el-collapse v-model="activeNames" @change="handleChange">
          <el-collapse-item title="一致性 Consistency" name="1">
            <div>与现实生活一致：与现实生活的流程、逻辑保持一致，遵循用户习惯的语言和概念；</div>
            <div>在界面中一致：所有的元素和结构需保持一致，比如：设计样式、图标和文本、元素的位置等。</div>
          </el-collapse-item>
          <el-collapse-item title="反馈 Feedback" name="2">
            <div>控制反馈：通过界面样式和交互动效让用户可以清晰的感知自己的操作；</div>
            <div>页面反馈：操作后，通过页面元素的变化清晰地展现当前状态。</div>
          </el-collapse-item>
          <el-collapse-item title="效率 Efficiency" name="3">
            <div>简化流程：设计简洁直观的操作流程；</div>
            <div>清晰明确：语言表达清晰且表意明确，让用户快速理解进而作出决策；</div>
            <div>帮助用户识别：界面简单直白，让用户快速识别而非回忆，减少用户记忆负担。</div>
          </el-collapse-item>
          <el-collapse-item title="可控 Controllability" name="4">
            <div>用户决策：根据场景可给予用户操作建议或安全提示，但不能代替用户进行决策；</div>
            <div>结果可控：用户可以自由的进行操作，包括撤销、回退和终止当前操作等。</div>
          </el-collapse-item>
        </el-collapse>
      </el-col>
    </el-row>
  </div>
</template>

<script>
export default {
  name: "Dash",
  data() {
    return {
      activeNames: ["1", "2", "3", "4"]
    };
  },
  methods: {
    handleChange(val) {
      // console.log(val);
    }
  }
};
</script>

<style scoped>
.row {
  margin-bottom: 20px;
}
.image {
  height: 200px;
}
.card-title {
  font-size: 1.2em;
  color: #443931;
}

.text {
  text-align: left;
}
</style>
