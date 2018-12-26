const BREADCRUMBS_KEY = 'breadcumbs'

export function saveBreadcrumbs(breadcrumbs) {
  let s = JSON.stringify(breadcrumbs)
  sessionStorage.setItem(BREADCRUMBS_KEY, s)
}

export function loadBreadcrumbs() {
  let s = sessionStorage.getItem(BREADCRUMBS_KEY)
  let breadcrumbs = JSON.parse(s)
  return breadcrumbs
}

export const plgCfg = {};
(function(plugins){
  plugins["plugin.httplog"] = {
    name: "日志记录",
    description:
      "日志插件是一个用于控制网关对于每一个请求进行控制，重新调度的插件",
    key: "plugin.httplog",
    to: "/configs/plugin/httplog"
  }
  plugins["plugin.proxy"] = {
    name: "代理",
    description:
      "代理插件是一个用于控制网关对于每一个请求进行控制，重新调度的插件",
    key: "plugin.proxy",
    to: "/configs/plugin/proxy"
  }
  plugins["plugin.cache"] = {
    name: "缓存",
    description:
      "缓存插件是一个用于控制网关对于每一个请求进行控制，重新调度的插件",
    key: "plugin.cache",
    to: "/configs/plugin/cache"
  }
  plugins["plugin.ratelimit"] = {
    name: "流量控制",
    description:
      "流量控制插件是一个用于控制网关对于每一个请求进行控制，重新调度的插件",
    key: "plugin.ratelimit",
    to: "/configs/plugin/ratelimit"
  }
})(plgCfg)