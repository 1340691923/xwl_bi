<template>
  <div class="app-container">
    <el-card class="box-card">

      <div class="filter-container">
        <el-button type="primary" icon="el-icon-plus" class="filter-item" @click="handleAddRole">新建用户</el-button>
      </div>
      <el-table
        v-loading="loading"
        border
        :data="rolesList"
      >
        <el-table-column
          label="序号"
          align="center"
          fixed
          width="50"
        >
          <template slot-scope="scope">
            {{ scope.$index+1 }}
          </template>
        </el-table-column>
        <el-table-column align="center" label="id" width="220">
          <template slot-scope="scope">
            {{ scope.row.id }}
          </template>
        </el-table-column>
        <el-table-column align="center" label="用户名" width="200">
          <template slot-scope="scope">
            {{ scope.row.username }}
          </template>
        </el-table-column>

        <el-table-column align="center" label="所属角色" width="200">
          <template slot-scope="scope">
            {{ scope.row.role_name }}
          </template>
        </el-table-column>
        <el-table-column align="center" label="真实姓名" width="200">
          <template slot-scope="scope">
            {{ scope.row.realname }}
          </template>
        </el-table-column>

        <el-table-column align="center" label="创建时间">
          <template slot-scope="scope">
            {{ scope.row.create_time }}
          </template>
        </el-table-column>

        <el-table-column align="center" label="最后登录时间">
          <template slot-scope="scope">
            {{ scope.row.last_login_time }}
          </template>
        </el-table-column>
        <el-table-column align="center" label="是否被禁用">
          <template slot-scope="scope">
            <span v-if="scope.row.is_del==1" style="color:red">已禁用</span>
            <span v-else>否</span>
          </template>
        </el-table-column>
        <el-table-column align="center" label="操作" fixed="right" width="300">
          <template slot-scope="scope">
            <el-button-group>
              <el-button
                v-if="scope.row.is_del == 0"
                type="info"
                size="small"
                icon="el-icon-delete-solid"
                @click="banUser(scope.row.id,scope.row.is_del==0?1:0)"
              >禁用
              </el-button>
              <el-button
                v-if="scope.row.is_del == 1"
                type="success"
                size="small"
                icon="el-icon-edit"
                @click="banUser(scope.row.id,scope.row.is_del==0?1:0)"
              >解封
              </el-button>
              <el-button type="primary" size="small" icon="el-icon-edit" @click="handleEdit(scope)">编辑</el-button>
              <el-button type="danger" size="small" icon="el-icon-delete" @click="handleDelete(scope)">删除</el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>

      <el-dialog
        :close-on-click-modal="false"
        :visible.sync="dialogVisible"
        :title="dialogType==='edit'?'编辑用户信息':'新建用户信息'"
      >
        <el-form :model="role" label-width="100px" label-position="left">

          <el-form-item label="用户名">
            <el-input v-model="role.username" placeholder="用户名" />
          </el-form-item>
          <el-form-item label="密码">
            <el-input v-model="role.password" placeholder="密码" />
          </el-form-item>
          <el-form-item label="真实姓名">
            <el-input v-model="role.realname" placeholder="真实姓名" />
          </el-form-item>
          <el-form-item label="请选择角色">
            <el-select v-model="role.role_id" placeholder="请选择角色" clearable filterable>
              <el-option
                v-for="item in chanCfgList"
                :key="item.id"
                :label="item.name"
                :value="item.id"
              />
            </el-select>
          </el-form-item>

        </el-form>
        <div style="text-align:right;">
          <el-button type="danger" icon="el-icon-close" @click="dialogVisible=false">取消</el-button>
          <el-button type="primary" icon="el-icon-check" @click="confirmRole">确认</el-button>
        </div>
      </el-dialog>
    </el-card>
    <back-to-top />
  </div>
</template>

<script>
import { deepClone } from '@/utils'
import { Ban, DelUser, InsertUser, roleOption, UpdateUser, userList } from '@/api/user'

const defaultUser = {
  id: '',
  password: '',
  realname: '',
  role_id: '',
  username: '',
  role_name: ''
}

export default {
  components: {
    BackToTop: () => import('@/components/BackToTop/index')
  },

  data() {
    return {
      loading: false,
      role: Object.assign({}, defaultUser),
      routes: [],
      rolesList: [],
      dialogVisible: false,
      dialogType: 'new',
      checkStrictly: false,
      defaultProps: {
        children: 'children',
        label: 'title'
      },
      chanCfgList: [],
      chanCfgMap: []
    }
  },
  computed: {},
  created() {
    this.init()
  },
  methods: {

    async init() {
      await this.getRoleOpt()
      await this.getUserList()
    },

    async banUser(id, typ) {
      const res = await Ban({ id: id, typ: typ })
      if (res.code != 0) {
        this.$message({
          offset: 60,

          message: res.msg,
          type: 'error'
        })
        return
      }

      this.init()
    },
    async getRoleOpt() {
      const res = await roleOption()
      for (var v of res.data) {
        this.chanCfgList.push(v)
        this.chanCfgMap[v['id']] = v['name']
      }
    },

    async getUserList() {
      this.loading = true
      const res = await userList()
      this.loading = false
      for (var k in res.data) {
        res.data[k]['role_name'] = this.chanCfgMap[res.data[k]['role_id']]
      }

      this.rolesList = res.data
    },

    handleAddRole() {
      this.role = Object.assign({}, defaultUser)
      this.dialogType = 'new'
      this.dialogVisible = true
    },
    handleEdit(scope) {
      this.dialogType = 'edit'
      this.dialogVisible = true
      this.checkStrictly = true
      this.role = deepClone(scope.row)
      this.role.password = ''
    },
    handleDelete({ $index, row }) {
      this.$confirm('确定删除该用户吗?', '警告', {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(async() => {
          await DelUser({ id: row.id })
          this.rolesList.splice($index, 1)
          this.$message({
            offset: 60,

            type: 'success',
            message: 'Delete succed!'
          })
        })
        .catch(err => {
          console.error(err)
        })
    },
    async confirmRole() {
      const isEdit = this.dialogType === 'edit'
      if (this.role.password.length < 5 && this.role.password.trim()!="") {
        this.$message({
          offset: 60,

          message: '密码长度必须大于5',
          type: 'error'
        })
        return false
      }
      if (this.role.realname == '') {
        this.$message({
          offset: 60,

          message: '用户真实姓名不能为空',
          type: 'error'
        })
        return false
      }
      if (this.role.username == '') {
        this.$message({
          offset: 60,

          message: '用户名不能为空',
          type: 'error'
        })
        return false
      }
      if (this.role.password == '') {
        this.$message({
          offset: 60,

          message: '密码不能为空',
          type: 'error'
        })
        return false
      }
      if (isEdit) {
        await UpdateUser(this.role)
        this.getUserList()
      } else {
        const { data } = await InsertUser(this.role)
        this.getUserList()
        /* this.role.id = data
            this.rolesList.push(this.role)*/
      }

      const { username, password, id, role_name, realname } = this.role
      this.dialogVisible = false
      this.$notify({
        title: 'Success',
        dangerouslyUseHTMLString: true,
        message: `
            <div>id: ${id}</div>
            <div>用户名: ${username}</div>
            <div>密码: ${password}</div>
            <div>角色: ${role_name}</div>
            <div>真实姓名: ${realname}</div>
          `,
        type: 'success'
      })
    }
  }
}
</script>

