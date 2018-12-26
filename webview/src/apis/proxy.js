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
export function renameReverseSrvGroupName({ group, newname }) {
    return putAPI({
        uri: "/gateapi/plugin/proxy/reversesrv/" + group,
        params: { newname}
    })
}
export function delReverseSrvGroup({ group }) {
    return deleteAPI({
        uri: "/gateapi/plugin/proxy/reversesrv/" + group,
        params: null
    })
}

// server rules config api functions

export function getServerRules({ limit, offset }) {
    return getAPI({ uri: "/gateapi/plugin/proxy/srvrules", params: { limit, offset } })
}

export function newServerRule({ prefix, server_name, need_strip_prefix }) {
    return postAPI({ uri: "/gateapi/plugin/proxy/srvrule", params: { prefix, server_name, need_strip_prefix } })
}

export function editServerRule({ id, prefix, server_name, need_strip_prefix }) {
    return putAPI({ uri: "/gateapi/plugin/proxy/srvrule/" + id, params: { prefix, server_name, need_strip_prefix } })
}

export function delServerRule({ id }) {
    return deleteAPI({ uri: "/gateapi/plugin/proxy/srvrule/" + id, params: null })
}