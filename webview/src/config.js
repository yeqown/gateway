// export const breadcrumb = {
//     'configs': '配置',
//     'gate': '网关',
//     'cache': '缓存',
//     'dashbord': '导航',
//     'basic': '基本',
// }

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