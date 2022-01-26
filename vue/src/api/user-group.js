import request from '@/utils/request'

var api = '/api/user_group/'

export function AddUserGroup(data) {
  return request({
    url: api + 'AddUserGroup',
    method: 'post',
    data
  })
}

export function ModifyUserGroup(data) {
  return request({
    url: api + 'ModifyUserGroup',
    method: 'post',
    data
  })
}

export function DeleteUserGroup(data) {
  return request({
    url: api + 'DeleteUserGroup',
    method: 'post',
    data
  })
}

export function UserGroupList(data) {
  return request({
    url: api + 'UserGroupList',
    method: 'post',
    data
  })
}

export function UserGroupSelect(data) {
  return request({
    url: api + 'UserGroupSelect',
    method: 'post',
    data
  })
}
