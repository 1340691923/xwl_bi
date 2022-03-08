<template>
  <div class="app-container">
    <el-card class="box-card">
      <div class="filter-container">
        <el-tag size="medium" class="filter-item">搜索应用名</el-tag>
        <el-input v-model="input.app_name" style="width:300px" class="filter-item" @input="search(1)" />
        <el-tag size="medium" class="filter-item">状态</el-tag>
        <el-select v-model="input.is_close" clearable class="filter-item" @input="search(1)">
          <el-option label="开启" :value="Number(0)" />
          <el-option label="关闭" :value="Number(1)" />
        </el-select>
        <el-button type="primary" icon="el-icon-plus" class="filter-item" @click.native="dialogVisible = true">新建应用
        </el-button>
      </div>
      <el-table
        v-loading="tableLoading"
        border
        :data="tableData"
        style="width: 100%"
      >

        <el-table-column
          align="center"
          fixed
          prop="app_name"
          label="应用名"
          width="180"
        >
          <template slot-scope="scope">
            {{ scope.row.app_name }}({{ scope.row.id }})
          </template>
        </el-table-column>

        <el-table-column
          align="center"
          prop="app_id"
          label="APPID"
          width="180"
        />
        <el-table-column
          align="center"
          prop="app_key"
          label="APPKEY"
          width="180"
        />
        <el-table-column
          align="center"
          prop="descibe"
          label="描述"
          width="180"
        />
        <el-table-column label="状态" width="180" align="center">
          <template slot-scope="scope">
            <template v-if="scope.row.is_close == 1">
              <el-tag type="danger" style="margin-top: 2px">
                关闭
              </el-tag>&nbsp;&nbsp;
            </template>
            <template v-else>
              <el-tag type="primary" style="margin-top: 2px">
                开启
              </el-tag>&nbsp;&nbsp;
            </template>
          </template>
        </el-table-column>
        <el-table-column label="成员" width="180" align="center">
          <template slot-scope="scope">
            <template v-for="(app_manager,index) in getManagerName(scope.row.app_manager)">
              <el-tag :type="index%2 == 0 ?'success':'primary'" style="margin-top: 2px">
                {{ app_manager }}
              </el-tag>&nbsp;&nbsp;
            </template>
          </template>
        </el-table-column>
        <el-table-column label="数据保留月数" align="center">
          <template slot-scope="scope">
            {{ scope.row.save_mouth }}
          </template>
        </el-table-column>
        <el-table-column label="创建人" align="center">
          <template slot-scope="scope">
            {{ userConfig[scope.row.create_by] }}
          </template>
        </el-table-column>
        <el-table-column label="更新人" align="center">
          <template slot-scope="scope">
            {{ userConfig[scope.row.update_by] }}
          </template>
        </el-table-column>

        <el-table-column
          width="180"
          align="center"
          sortable
          prop="create_time"
          label="创建时间"
        />
        <el-table-column
          width="180"
          align="center"
          sortable
          prop="create_time"
          label="更新时间"
        />
        <el-table-column
          fixed="right"
          label="操作"
          width="450"
          align="center"
        >
          <template slot-scope="scope">

            <el-button
              size="mini"
              type="primary"
              class="copyParmas"
              :data-clipboard-text="JSON.stringify({appid:scope.row.app_id,appkey:scope.row.app_key})"
              icon="el-icon-copy-document"
              @click="copyParmas"
            >复制参数
            </el-button>
            <el-button
              size="mini"
              type="info"
              icon="el-icon-refresh"
              @click="resetKey(scope.row.id,scope.row.app_id,scope.row.app_name)"
            >更新秘钥
            </el-button>
            <el-button size="mini" type="success" icon="el-icon-edit" @click="openManagerForm(scope.row)">操作成员
            </el-button>
            <el-button
              v-if="scope.row.is_close == 1"
              size="mini"
              type="success"
              icon="el-icon-open"
              @click="statusOperation(0,scope.row.app_name,scope.row.app_id,scope.row.app_key,scope.row.id)"
            >开启
            </el-button>
            <el-button
              v-if="scope.row.is_close == 0"
              size="mini"
              type="danger"
              icon="el-icon-close"
              @click="statusOperation(1,scope.row.app_name,scope.row.app_id,scope.row.app_key,scope.row.id)"
            >关闭
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <div
        style="width: 98%; margin-left: 20px; position: sticky;bottom: 10px;height: 80px;z-index: 30;background: white; display: flex;align-items: center;justify-content: right;padding: 0 24px"
      >
        <div>
          <el-pagination
            :current-page="page"
            :page-sizes=" [10, 20, 30, 50,100,150,200]"
            :page-size="limit"
            layout="total, sizes, prev, pager, next, jumper"
            :total="total"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </div>
      <el-dialog :close-on-click-modal="false" :visible.sync="dialogVisible" title="新增应用" @close="formClose">
        <el-form :model="form" label-width="120px" label-position="left">
          <el-form-item label="应用名:">
            <el-input v-model="form.app_name" placeholder="应用名" />
          </el-form-item>
          <el-form-item label="应用描述:">
            <el-input v-model="form.descibe" placeholder="应用描述" />
          </el-form-item>
          <el-form-item label="数据保留月数:">
            <el-input v-model="form.save_mouth" type="number" placeholder="数据保留月数" />
          </el-form-item>
        </el-form>
        <div style="text-align:right;">
          <el-button type="danger" icon="el-icon-close" @click="formClose">返回</el-button>
          <el-button type="primary" icon="el-icon-plus" @click="addForm">添加</el-button>
        </div>
      </el-dialog>

      <el-dialog
        width="70%"
        :close-on-click-modal="false"
        :visible.sync="managerFormdialogVisible"
        title="操作成员"
        @close="managerFormdialogVisible = false"
      >
        <el-form :model="managerForm" label-width="120px" label-position="left">
          <el-form-item label="应用名">
            <el-input v-model="managerForm.app_name" placeholder="应用名" />
          </el-form-item>
          <el-form-item label="成员">
            <el-transfer
              v-if="managerFormdialogVisible"
              v-model="managerForm.app_manager"
              :titles="['全部用户', '当前应用成员']"
              :button-texts="['移除成员', '添加成员']"
              filterable
              :filter-method="filterMethod"
              filter-placeholder="请操作成员"
              :data="allUserConfig"
            />

          </el-form-item>
        </el-form>
        <div style="text-align:right;">
          <el-button type="danger" icon="el-icon-close" @click="managerFormdialogVisible = false">返回</el-button>
          <el-button type="primary" icon="el-icon-check" @click="addManager">添加</el-button>
        </div>
      </el-dialog>
    </el-card>
    <back-to-top />
  </div>
</template>

<script>
import Clipboard from 'clipboard'
import { Create, List, ResetAppkey, StatusAction, UpdateManager } from '@/api/app'
import { userList } from '@/api/user'

export default {
  name: 'Index',
  components: {
    BackToTop: () => import('@/components/BackToTop/index')
  },
  data() {
    return {
      dialogVisible: false,
      managerFormdialogVisible: false,
      form: {
        app_name: '',
        app_key: '',
        save_mouth: 1
      },
      input: {
        limit: 10,
        page: 1,
        app_name: '',
        is_close: null
      },
      total: 0,
      tableData: [],
      tableLoading: false,
      userConfig: {},
      formDefault: {},
      managerForm: {
        app_id: '',
        app_name: '',
        create_by: 0,
        app_manager: []
      },
      allUserConfig: []
    }
  },
  mounted() {
    this.formDefault = JSON.parse(JSON.stringify(this.form))
    this.init()
  },
  methods: {
    async addManager() {
      const form = {
        app_id: this.managerForm.app_id,
        app_manager: this.managerForm.app_manager.join(',')
      }
      const res = await UpdateManager(form)
      if (res.code != 0) {
        this.$message({
          showClose: true,
          offset: 60,
          type: 'error',
          message: res.msg
        })
        return
      } else {
        this.$message({
          showClose: true,
          offset: 60,
          type: 'success',
          message: res.msg
        })
        this.search(1)
        this.managerFormdialogVisible = false
      }
    },
    filterMethod(query, item) {
      return item.label.indexOf(query) > -1
    },
    async openManagerForm(row) {
      this.managerForm.app_id = row.app_id
      this.managerForm.app_name = row.app_name
      this.managerForm.app_manager = row.app_manager.split(',')

      this.allUserConfig = []
      const userListRes = await userList()
      if (userListRes) {
        for (var v of userListRes.data) {
          this.userConfig[v.id] = v.realname
          const obj = {
            label: v.realname,
            key: v.id.toString(),
            disabled: false
          }
          if (row.create_by == v.id) {
            obj.disabled = true
          }
          this.allUserConfig.push(
            obj
          )
        }
      }

      this.managerFormdialogVisible = true
    },

    async addForm() {
      const res = await Create(this.form)
      if (res.code != 0) {
        this.$message({
          showClose: true,
          offset: 60,
          type: 'error',
          message: res.msg
        })
        return
      } else {
        this.$message({
          showClose: true,
          offset: 60,
          type: 'success',
          message: res.msg
        })
        this.search(1)
      }
    },
    formClose() {
      this.form = JSON.parse(JSON.stringify(this.formDefault))
      this.dialogVisible = false
    },
    copyParmas() {
      var clipboard = new Clipboard('.copyParmas')
      clipboard.on('success', e => {
        this.$message({
          showClose: true,
          offset: 60,
          type: 'success',
          message: '复制成功！'
        })
      })
      clipboard.on('error', e => {
        this.$message({
          showClose: true,
          offset: 60,
          type: 'error',
          message: e
        })
      })
    },

    async resetKey(id, appid, appname) {
      this.$confirm('确定重置应用名为 【' + appname + '】 的这个秘钥吗?', '警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(async() => {
          const res = await ResetAppkey({ app_id: appid, id: id })
          if (res.code != 0) {
            this.$message({
              showClose: true,
              offset: 60,
              type: 'error',
              message: res.msg
            })
            return
          } else {
            this.$message({
              showClose: true,
              offset: 60,
              type: 'success',
              message: res.msg
            })
            this.search(1)
          }
        })
        .catch(err => {
          console.error(err)
        })
    },
    getManagerName(managerStr) {
      var arr = managerStr.split(',')
      for (const k in arr) {
        arr[k] = this.userConfig[arr[k]]
      }
      return arr
    },
    async statusOperation(is_close, app_name, app_id, app_key, id) {
      this.$confirm('确定修改应用名为 【' + app_name + '】 的状态吗?', '警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(async() => {
          const res = await StatusAction({ app_id: app_id, app_key: app_key, is_close: is_close, id: id })
          if (res.code != 0) {
            this.$message({
              showClose: true,
              offset: 60,
              type: 'error',
              message: res.msg
            })
            return
          } else {
            this.$message({
              showClose: true,
              offset: 60,
              type: 'success',
              message: res.msg
            })
            this.search(1)
          }
        })
        .catch(err => {
          console.error(err)
        })
    },
    async search(page) {
      !page ? this.input.page = 1 : this.input.page = page
      this.tableLoading = true
      if (this.input.is_close != null) {
        if (this.input.is_close.toString() == '') {
          this.input.is_close = null
        }
      }
      const res = await List(this.input)

      if (res.code != 0) {
        this.$message({
          showClose: true,
          offset: 60,
          type: 'error',
          message: res.msg
        })
        this.tableLoading = false
        return
      }
      this.tableData = res.data.list
      this.total = res.data.count

      this.tableLoading = false
    },
    async init() {
      const userListRes = await userList()
      if (userListRes) {
        for (let v of userListRes.data) {
          this.userConfig[v.id] = v.realname
          this.allUserConfig.push(
            {
              label: v.realname,
              key: v.id.toString(),
              disabled: false
            }
          )
        }
      }
      await this.search()
    }
  }
}
</script>

