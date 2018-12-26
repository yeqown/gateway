import VueRouter from "vue-router"
import Dash from './views/Dash'
import Config from './views/Config'
import Plugins from './views/Plugins'
import PluginConfig from './views/PluginConfig'
import BaseConfig from './views/configviews/GatewayBasicConf'
import PluginCache from './views/configviews/PluginCache'
import PluginProxy from '@/views/configviews/PluginProxy'
import ReverseServer from '@/views/configviews/ProxyReverseServer'
import PathRule from '@/views/configviews/ProxyPath'
import ServerRule from '@/views/configviews/ProxyServer'

import ReverseServerGroup from '@/components/ReverseServerGroup'
import PathRuleDetail from '@/components/PathRuleDetail'


const router = new VueRouter({
  routes: [
    {
      path: '/',
      redirect: _ => {
        return '/dashbord'
      },
    },
    {
      name: 'dashbord',
      path: '/dashbord',
      component: Dash,
      meta: { breadcrumb: '导航页' }
    },
    {
      name: 'gatebase',
      path: '/configs/',
      // redirect: _ => {
      //   return '/plugins'
      // },
      component: Config,
      meta: { breadcrumb: '配置' },
      children: [
        {
          name: 'base',
          path: 'basic',
          component: BaseConfig,
          meta: { breadcrumb: '基本' }
        },
        {
          name: 'plugin',
          path: 'plugin/',
          component: PluginConfig,
          redirect: _ => {
            return '/plugins'
          },
          meta: { breadcrumb: '插件' },
          children: [
            {
              name: 'plugin.cache',
              path: 'cache',
              component: PluginCache,
              meta: { breadcrumb: '缓存插件' }
            },
            {
              name: 'plugin.proxy',
              path: 'proxy/',
              component: PluginProxy,
              meta: { breadcrumb: '代理插件' },
              children: [
                {
                  name: 'plugin.proxy.reverseServer',
                  path: 'reverse_server',
                  component: ReverseServer,
                  meta: { breadcrumb: '反向代理组' }
                },
                {
                  name: 'plugin.proxy.reverseServer.group',
                  path: 'reverse_server/:group',
                  component: ReverseServerGroup,
                  meta: { breadcrumb: '组别详情' }
                },
                {
                  name: 'plugin.proxy.serverrule',
                  path: 'serverrule',
                  component: ServerRule,
                  meta: { breadcrumb: '服务' }
                },
                {
                  name: 'plugin.proxy.pathrule',
                  path: 'pathrule',
                  component: PathRule,
                  meta: { breadcrumb: 'URI' }
                },
                {
                  name: 'plugin.proxy.pathrule.detail',
                  path: 'pathrule/:id',
                  component: PathRuleDetail,
                  meta: { breadcrumb: '详情' }
                },
              ]
            }
          ]
        }
      ]
    },
    {
      name: 'plugins-manage',
      path: '/plugins',
      component: Plugins,
      meta: { breadcrumb: '插件管理' }
    }
  ]
})

router.beforeEach((to, from, next) => {
  let _ = from
  // console.log('router.beforeEach called', to, from)
  to.params.breadcrumbs = []
  to.matched.forEach(matched => {
    // console.log(matched)
    to.params.breadcrumbs.push({
      name: matched.meta.breadcrumb,
      to: matched.path
    })
  })
  next()
})

// router.afterEach((to, from) => {
//   console.log("router.afterEach called", to, from)
// })

export default router