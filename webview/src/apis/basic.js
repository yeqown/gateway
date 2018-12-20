import {getAPI} from '.'

export function getGlobalConfig() {
    return getAPI({uri:'config', params:null})
}