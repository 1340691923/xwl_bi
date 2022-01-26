import request from '@/utils/request'

var api = '/api/app/'

export function Create(data) {
  return request({
    url: api + 'Create',
    method: 'post',
    data
  })
}

export function ResetAppkey(data) {
  return request({
    url: api + 'ResetAppkey',
    method: 'post',
    data
  })
}

export function List(data) {
  return request({
    url: api + 'List',
    method: 'post',
    data
  })
}

export function UpdateManager(data) {
  return request({
    url: api + 'UpdateManager',
    method: 'post',
    data
  })
}

export function Config(data) {
  return request({
    url: api + 'Config',
    method: 'post',
    data
  })
}
export function StatusAction(data) {
  return request({
    url: api + 'StatusAction',
    method: 'post',
    data
  })
}

