import request from '@/utils/request'

var api = '/api/pannel/'

export function ReportTableList(data) {
  return request({
    url: api + 'ReportTableList',
    method: 'post',
    data
  })
}

export function DeleteReportTableByID(data) {
  return request({
    url: api + 'DeleteReportTableByID',
    method: 'post',
    data
  })
}

export function AddReportTable(data) {
  return request({
    url: api + 'AddReportTable',
    method: 'post',
    data
  })
}

export function FindNameCount(data) {
  return request({
    url: api + 'FindNameCount',
    method: 'post',
    data
  })
}

export function FindRtById(data) {
  return request({
    url: api + 'FindRtById',
    method: 'post',
    data
  })
}

export function GetPannelList(data) {
  return request({
    url: api + 'GetPannelList',
    method: 'post',
    data
  })
}

export function NewDir(data) {
  return request({
    url: api + 'NewDir',
    method: 'post',
    data
  })
}

export function NewPannel(data) {
  return request({
    url: api + 'NewPannel',
    method: 'post',
    data
  })
}

export function Rename(data) {
  return request({
    url: api + 'Rename',
    method: 'post',
    data
  })
}

export function MovePannel2Dir(data) {
  return request({
    url: api + 'MovePannel2Dir',
    method: 'post',
    data
  })
}

export function DeletePannel(data) {
  return request({
    url: api + 'DeletePannel',
    method: 'post',
    data
  })
}

export function DeleteDir(data) {
  return request({
    url: api + 'DeleteDir',
    method: 'post',
    data
  })
}

export function CopyPannel(data) {
  return request({
    url: api + 'CopyPannel',
    method: 'post',
    data
  })
}

export function UpdatePannelRt(data) {
  return request({
    url: api + 'UpdatePannelRt',
    method: 'post',
    data
  })
}

export function UpdatePannelManager(data) {
  return request({
    url: api + 'UpdatePannelManager',
    method: 'post',
    data
  })
}

export function RtListByAppid(data) {
  return request({
    url: api + 'RtListByAppid',
    method: 'post',
    data
  })
}

