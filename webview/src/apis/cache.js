import { getAPI, putAPI, postAPI, deleteAPI } from '.'

export function newCacheRule({ regular, enabled }) {
    return postAPI({ uri: "/gateapi/plugin/cache/rule", params: { regular, enabled } })
}

export function delCacheRule({ id }) {
    return deleteAPI({ uri: "/gateapi/plugin/cache/rule/" + id , params: null})
}

export function getCacheRules({limit = 10, offset = 0}) {
    return getAPI({ uri: "/gateapi/plugin/cache/rules", params: { offset, limit } })
}

export function editCacheRule({ id, regular, enabled }) {
    return putAPI({ uri: "/gateapi/plugin/cache/rule/" + id, params: { regular, enabled } })
}