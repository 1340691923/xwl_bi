import request from '@/utils/request'

var api = '/api/metadata/'

export function MetaEventList(data) {
  return request({
    url: api + 'MetaEventList',
    method: 'post',
    data
  })
}

export function UpdateShowName(data) {
  return request({
    url: api + 'UpdateShowName',
    method: 'post',
    data
  })
}

export function UpdateAttrShowName(data) {
  return request({
    url: api + 'UpdateAttrShowName',
    method: 'post',
    data
  })
}

export function AttrManager(data) {
  return request({
    url: api + 'AttrManager',
    method: 'post',
    data
  })
}

export function MetaEventListByAttr(data) {
  return request({
    url: api + 'MetaEventListByAttr',
    method: 'post',
    data
  })
}

export function AttrManagerByMeta(data) {
  return request({
    url: api + 'AttrManagerByMeta',
    method: 'post',
    data
  })
}

export function GetAnalyseSelectOptions(data) {
  return request({
    url: api + 'GetAnalyseSelectOptions',
    method: 'post',
    data
  })
}
export function UpdateAttrInvisible(data) {
  return request({
    url: api + 'UpdateAttrInvisible',
    method: 'post',
    data
  })
}
