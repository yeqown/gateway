import axios from 'axios'
import queryString from 'query-string'
import * as basicapi from './basic'
export { basicapi }

const baseURL = 'http://localhost:8989/gateapi'

var instance = axios.create({
    baseURL: baseURL,
    timeout: 5000,
    headers: {},
})

instance.interceptors.request.use((config) => {
    // before request
    console.log(config)
    return config
}, (error) => {
    // handle with error
    console.log(error)
    return Promise.reject(error)
})

instance.interceptors.response.use((response) => {
    console.log(response)
    if (response.status > 300) {
        throw Error("wrong status code got")
    }
    return response.data
}, (error) => {
    console.log(error)
    return Promise.reject(error)
})

function requestAPI(config) {
    return instance.request(config)
}

export function getAPI({ uri, params }) {
    let qs = queryString.stringify(params)
    return requestAPI({
        method: 'get',
        url: uri,
        params: qs,
        responseType: 'json'
    })
}

export function postAPI({ uri, params, headers }) {
    if (!headers) {
        // throw Error('headers cannot be empty')
        let headers = {}
    }
    headers['Content-Type'] = 'application/x-www-form-urlencoded'

    return requestAPI({
        method: 'get',
        url: uri,
        data: params,
        responseType: 'json',
        headers: headers
    })
}

export function deleteAPI({ uri, params }) {
    return requestAPI({
        method: 'delete',
        url: uri,
        data: params,
        responseType: 'json',
    })
}

export function putAPI({ uri, params }) {
    return requestAPI({
        method: 'put',
        url: uri,
        data: params,
        responseType: 'json'
    })
}