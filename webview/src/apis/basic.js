import { getAPI, putAPI } from '.'

export function getGlobalConfig() {
    return getAPI({ uri: '/gate/config', params: null })
}

export function putGlobalConfig({logpath, port}) {
    return putAPI({ uri: '/gate/config', params: {logpath, port} })
}