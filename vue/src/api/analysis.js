import request from '@/utils/request'

var api = '/api/analysis/'

export function GetConfigs(data) {
  return request({
    url: api + 'GetConfigs',
    method: 'post',
    data
  })
}

export function LoadPropQuotas(data) {
  return request({
    url: api + 'LoadPropQuotas',
    method: 'post',
    data
  })
}

export function GetValues(data) {
  return request({
    url: api + 'GetValues',
    method: 'post',
    data
  })
}

export function FunnelList(data) {
  return request({
    url: api + 'FunnelList',
    method: 'post',
    data
  })
}

export function RetentionList(data) {
  return request({
    url: api + 'RetentionList',
    method: 'post',
    data
  })
}

export function TraceList(data) {
  return request({
    url: api + 'TraceList',
    method: 'post',
    data
  })
}

export function UserAttrList(data) {
  return request({
    url: api + 'UserAttrList',
    method: 'post',
    data
  })
}

export function EventList(data) {
  return request({
    url: api + 'EventList',
    method: 'post',
    data
  })
}

export function UserList(data) {
  return request({
    url: api + 'UserList',
    method: 'post',
    data
  })
}

export function UserEventDetailList(data) {
  return request({
    url: api + 'UserEventDetailList',
    method: 'post',
    data
  })
}

export function UserEventCountList(data) {
  return request({
    url: api + 'UserEventCountList',
    method: 'post',
    data
  })
}
