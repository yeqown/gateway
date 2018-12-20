import VueRouter from "vue-router"
// import HelloWorld from './components/HelloWorld'
import Dash from './views/Dash'
import Config from './views/Config'
import Plugins from './views/Plugins'
import PluginConfig from './views/PluginConfig'
import BaseConfig from './views/configviews/GatewayBasicConf'
import PluginCache from './views/configviews/PluginCache'

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
      //   return '/configs/basic'
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
          meta: { breadcrumb: '插件' },
          children: [
            {
              name: 'plugin.cache',
              path: 'cache',
              component: PluginCache,
              meta: { breadcrumb: '缓存' }
            }
          ]
        }
      ]
    },
    {
      name: 'plugins-manage',
      path: '/plugins',
      component: Plugins,
      meta: {breadcrumb: '插件管理'}
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