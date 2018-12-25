import { getAPI, putAPI, postAPI, deleteAPI } from '.'

// reverse server config api functions

export function getReverseSrvGroup() {
    return getAPI({ uri: "/gateapi/plugin/proxy/reversesrvgroups", params: null })
}

export function getReverseSrvGroupDetail({ group, limit = 10, offset = 0 }) {
    return getAPI({
        uri: "/gateapi/plugin/proxy/reversesrv/" + group,
        params: { limit, offset }
    })
}

export function newReverseSrv({ group, name, addr, weight }) {
    return postAPI({
        uri: "/gateapi/plugin/proxy/reversesrv",
        params: { group, name, addr, weight }
    })
}

export function delReverseSrv({ group, id }) {
    return deleteAPI({
        uri: "/gateapi/plugin/proxy/reversesrv/" + group + "/" + id,
        params: null
    })
}

export function editReverseSrv({ group, id, name, weight, addr }) {
    return putAPI({
        uri: "/gateapi/plugin/proxy/reversesrv/" + group + "/" + id,
        params: { name, weight, addr, group }
    })
}

export function delReverseSrvGroup({ group }) {
    return deleteAPI({
        uri: "/gateapi/plugin/proxy/reversesrv/" + group,
        params: null
    })
}