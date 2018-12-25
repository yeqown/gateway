import axios from 'axios'
// import queryString from 'query-string'
import { Notification } from 'element-ui'
import * as basicapi from './basic'
import * as cacheapi from './cache'
import * as proxyapi from './proxy'

export { basicapi, cacheapi, proxyapi}

export const baseURL = 'http://localhost:8989'

var defaultHeaders = {};

(function (headers) {
    headers["Content-Type"] = "application/x-www-form-urlencoded"
})(defaultHeaders);


var instance = axios.create({
    baseURL: baseURL,
    timeout: 5000,
    headers: defaultHeaders,
})

// resetBaseURL
export function resetBaseURL({ baseURL }) {
    instance.defaults.baseURL = baseURL
    // console.log(baseURL, instance.defaults)
}

instance.interceptors.request.use((config) => {
    return config
}, (error) => {
    console.log(error)
    return Promise.reject(error)
})

instance.interceptors.response.use((response) => {
    // console.log(response)
    if (response.status > 300) {
        throw Error("wrong status got, " + response.status)
    }
    if (response.data.code !== 0) {
        Notification.error(response.data.message)
        throw Error(response.data.message)
    }
    return response.data

}, (error) => {
    Message.error(error.message)
})

function requestAPI(config) {
    return instance.request(config)
}

export function getAPI({ uri, params }) {
    // console.log(qs)
    return requestAPI({
        method: 'get',
        url: uri,
        params: params,
        responseType: 'json'
    })
}

export function postAPI({ uri, params = null, headers = defaultHeaders }) {
    headers['Content-Type'] = 'application/x-www-form-urlencoded'
    return requestAPI({
        method: 'post',
        url: uri,
        data: serializeForm(params),
        responseType: 'json',
        headers: headers
    })
}

export function deleteAPI({ uri, params = null }) {
    return requestAPI({
        method: 'delete',
        url: uri,
        data: serializeForm(params),
        responseType: 'json',
    })
}

export function putAPI({ uri, params = null, headers = defaultHeaders }) {
    return requestAPI({
        method: 'put',
        url: uri,
        data: serializeForm(params),
        responseType: 'json',
        headers: headers
    })
}

function serializeForm(params) {
    if (!params) {
        return new FormData()
    }
    let body = new FormData()
    Object.keys(params).map(key => {
        body.append(key, params[key])
    })

    return body
}