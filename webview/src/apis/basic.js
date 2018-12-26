import { getAPI, putAPI } from '.'

// export function getGlobalConfig() {
//     return getAPI({ uri: '/gate/config', params: null })
// }

// export function putGlobalConfig({logpath, port}) {
//     return putAPI({ uri: '/gate/config', params: {logpath, port} })
// }

// plugins api
export function getPluginsList() {
    return getAPI({uri: '/gate/plugins', params: null})
}

export function enablePlugin({idx, enabled}) {
    return getAPI({uri: "/gate/plugins/enable", params: {idx, enabled}})
}