import request from '@/utils/request'

var api = '/api/realdata/'

export function List(data) {
  return request({
    url: api + 'List',
    method: 'post',
    data
  })
}

export function FailDataList(data) {
  return request({
    url: api + 'FailDataList',
    method: 'post',
    data
  })
}

export function FailDataDesc(data) {
  return request({
    url: api + 'FailDataDesc',
    method: 'post',
    data
  })
}

export function ReportCount(data) {
  return request({
    url: api + 'ReportCount',
    method: 'post',
    data
  })
}

export function EventFailDesc(data) {
  return request({
    url: api + 'EventFailDesc',
    method: 'post',
    data
  })
}

export function DelDebugDeviceID(data) {
  return request({
    url: api + 'DelDebugDeviceID',
    method: 'post',
    data
  })
}
export function DebugDeviceIDList(data) {
  return request({
    url: api + 'DebugDeviceIDList',
    method: 'post',
    data
  })
}
export function AddDebugDeviceID(data) {
  return request({
    url: api + 'AddDebugDeviceID',
    method: 'post',
    data
  })
}

