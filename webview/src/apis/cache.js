import { getAPI, putAPI, postAPI, deleteAPI } from '.'

export function getCacheRules(limit = 10, offset = 0) {
    return getAPI({ uri: "/plugin/cacherules", params: { offset, limit } })
}

export function newCacheRule({ regexp }) {
    return postAPI({ uri: "/plugin/cacherule", params: { regexp } })
}