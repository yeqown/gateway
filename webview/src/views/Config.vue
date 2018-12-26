<template>
  <div>
    <transition name="fade">
      <router-view/>
    </transition>
    <template v-if="current">
      <el-row :gutter="1" class="row">
        <el-col :span="8" key="gate.basic" :offset="3">
          <el-card :body-style="{ padding: '0px',cursor: 'pointer' }" shadow="hover">
            <div @click="navigate('/configs/basic')">
              <img src="@/assets/gateway.jpg" class="image">
              <div style="padding: 10px;">
                <span class="card-title">
                  <i class="el-icon-setting" style="color:green"/>
                  网关基本配置
                </span>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="8" key="plugin.ratelimit" :offset="3">
          <el-card :body-style="{ padding: '0px',cursor: 'pointer' }" shadow="hover">
            <div @click="navigate('/configs/plugin/ratelimit')">
              <img src="@/assets/gateway.jpg" class="image">
              <div style="padding: 10px;">
                <span class="card-title">
                  <i class="el-icon-d-caret" style="color:green"/>
                  流量限制
                </span>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
      <el-row :gutter="1" class="row">
        <el-col :span="8" key="plugin.proxy" :offset="3">
          <el-card :body-style="{ padding: '0px',cursor: 'pointer' }" shadow="hover">
            <div @click="navigate('/configs/plugin/proxy')">
              <img src="@/assets/gateway.jpg" class="image">
              <div style="padding: 10px;">
                <span class="card-title">
                  <i class="el-icon-share" style="color:green"/>
                  代理
                </span>
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="8" key="plugin.cache" :offset="3">
          <el-card :body-style="{ padding: '0px',cursor: 'pointer' }" shadow="hover">
            <div @click="navigate('/configs/plugin/cache')">
              <img src="@/assets/gateway.jpg" class="image">
              <div style="padding: 10px;">
                <span class="card-title">
                  <i class="el-icon-refresh" style="color:green"/>
                  缓存
                </span>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </template>
  </div>
</template>
<script>
export default {
  name: "Config",
  data() {
    return {
      current: false
    };
  },
  watch: {
    $route: function(newVal, oldVal) {
      // console.log("config router", newVal)
      if (newVal.path === "/configs") {
        this.current = true;
      } else {
        this.current = false;
      }
    }
  },
  methods: {
    navigate: function(to) {
      // console.log("clicked", to);
      this.$router.push(to);
    }
  },
  created() {
    // console.log(this.$route)
    if (this.$route.path === "/configs") {
      this.current = true;
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
